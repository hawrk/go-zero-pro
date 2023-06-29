package logic

import (
	"algo_assess/assess-rpc-server/assessservice"
	"algo_assess/global"
	"context"
	"fmt"
	"github.com/spf13/cast"
	"io"
	"net/http"
	"os"
	"runtime"
	"time"

	"algo_assess/assess-api-server/internal/svc"
	"algo_assess/assess-api-server/internal/types"

	"github.com/xuri/excelize/v2"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	SheetName = "sheet1"

	Charater = [26]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R",
		"S", "T", "U", "V", "W", "X", "Y", "Z"}
	// EconomyHeader 经济性Excel表头
	EconomyHeader = [19]string{"用户ID", "用户名称", "算法厂商", "算法名称", "母单号", "证券代码", "证券名称", "行业", "市值", "流动性",
		"交易量", "盈亏", "收益率", "手续费", "流量费", "撤单率", "最小拆单单位", "成交效率", "创建时间"}
	// ProgressHeader 完成度Excel表头
	ProgressHeader = [14]string{"用户ID", "用户名称", "算法厂商", "算法名称", "母单号", "证券代码", "证券名称", "行业", "市值", "流动性",
		"完成度", "母单贴合度", "成交量贴合度", "创建时间"}
	// RiskHeader 风险度Excel表头
	RiskHeader = [14]string{"用户ID", "用户名称", "算法厂商", "算法名称", "母单号", "证券代码", "证券名称", "行业", "市值", "流动性",
		"最小贴合度", "收益率", "回撤比例", "创建时间"}
	// AssessHeader 绩效Excel表头
	AssessHeader = [13]string{"用户ID", "用户名称", "算法厂商", "算法名称", "母单号", "证券代码", "证券名称", "行业", "市值", "流动性",
		"VWAP滑点值", "绩效收益率", "创建时间"}
	// StabilityHeader 稳定性Excel表头
	StabilityHeader = [16]string{"用户ID", "用户名称", "算法厂商", "算法名称", "母单号", "证券代码", "证券名称", "行业", "市值", "流动性",
		"VWAP滑点标准差", "收益率标准差", "贴合度", "成交量贴合度标准差", "时间贴合度标准差", "创建时间"}
)

type UserProfileExportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserProfileExportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserProfileExportLogic {
	return &UserProfileExportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserProfileExportLogic) UserProfileExport(w http.ResponseWriter, req *types.ProfileExportReq) (err error) {
	l.Logger.Infof("in UserProfileExport, get req:%+v", req)
	start := cast.ToInt64(time.Unix(req.StartTime, 0).Format(global.TimeFormatMinInt))
	end := cast.ToInt64(time.Unix(req.EndTime, 0).Format(global.TimeFormatMinInt))
	var algoId int32
	if req.AlgoName != "" {
		// 先反查一下算法ID
		alReq := &assessservice.ChooseAlgoReq{
			ChooseType: 4,
			AlgoName:   req.AlgoName,
		}

		alRsp, err := l.svcCtx.AssessClient.ChooseAlgoInfo(l.ctx, alReq)
		if err != nil {
			l.Logger.Error("call rpc ChooseAlgoInfo error:", err)
			return nil
		}
		algoId = alRsp.GetAlgoId()
	}
	// 拼返回数据
	var lists []*assessservice.ProfileInfo // 汇总数据
	var total int64                        //  实际返回的数量
	var retTotal int64                     // 数据总条数
	var page int32 = 1                     // 当前页数
	firstReq := true                       // 首次请求
	for total < retTotal || firstReq {
		pReq := &assessservice.AlgoProfileReq{
			Provider:     req.Provider,
			AlgoTypeName: req.AlgoTypeName,
			AlgoId:       algoId,
			AlgoName:     req.AlgoName,
			UserId:       req.UserId,
			StartTime:    start,
			EndTime:      end,
			Page:         page,
			Limit:        1000,
			ProfileType:  req.ProfileType,
		}
		rsp, err := l.svcCtx.AssessClient.GetAlgoProfile(l.ctx, pReq)
		if err != nil {
			l.Logger.Error("profile rpc call error:", err)
			return nil
		}
		firstReq = false
		total += int64(len(rsp.GetInfo()))
		retTotal = rsp.GetTotal()
		lists = append(lists, rsp.GetInfo()...)
		page++
		if page > 100 { // 应该没有超过100 *1000 这么多的数据吧
			break
		}
		l.Logger.Info("in UserProfileExport, total:", total, ", retTotal:", retTotal)
	}
	// 拼装所有返回数据
	var fn, fp string // fn fileName fp filePath
	if req.ProfileType == 1 {
		fn, fp = l.BuildEconomyFile(lists)
	} else if req.ProfileType == 2 {
		fn, fp = l.BuildProgressFile(lists)
	} else if req.ProfileType == 3 {
		fn, fp = l.BuildRiskFile(lists)
	} else if req.ProfileType == 4 {
		fn, fp = l.BuildAssessFile(lists)
	} else if req.ProfileType == 5 {
		fn, fp = l.BuildStabilityFile(lists)
	} else {
		l.Logger.Error("unsupported profile type:", req.ProfileType)
		return nil
	}
	//fn := l.GetFileName()
	file, err := os.Open(fp)
	if err != nil {
		l.Logger.Error("open file error:", err)
		return nil
	}
	header := w.Header()
	header.Add("Content-Type", "application/octet-stream")
	header.Add("Content-Disposition", "filename="+fn)
	io.Copy(w, file)

	return nil
}

func (l *UserProfileExportLogic) BuildEconomyFile(info []*assessservice.ProfileInfo) (fn, fp string) {
	f := excelize.NewFile()
	index := f.NewSheet(SheetName)
	// 设置表头
	for k, v := range Charater {
		// A1, B1, C1,D1
		if k < len(EconomyHeader) {
			f.SetCellValue(SheetName, fmt.Sprintf("%s1", v), EconomyHeader[k])
		} else {
			break
		}
	}
	// A2,B2,C2,D2
	for i, j := range info {
		f.SetCellValue(SheetName, fmt.Sprintf("A%d", i+2), j.GetAccoutId())
		f.SetCellValue(SheetName, fmt.Sprintf("B%d", i+2), j.GetAccountName())
		f.SetCellValue(SheetName, fmt.Sprintf("C%d", i+2), j.GetProvider())
		f.SetCellValue(SheetName, fmt.Sprintf("D%d", i+2), j.GetAlgoName())
		f.SetCellValue(SheetName, fmt.Sprintf("E%d", i+2), j.GetAlgoOrderId())
		f.SetCellValue(SheetName, fmt.Sprintf("F%d", i+2), j.GetSecId())
		f.SetCellValue(SheetName, fmt.Sprintf("G%d", i+2), j.GetSecName())
		f.SetCellValue(SheetName, fmt.Sprintf("H%d", i+2), j.GetIndustry())
		f.SetCellValue(SheetName, fmt.Sprintf("I%d", i+2), getFundType(j.GetFundType()))
		f.SetCellValue(SheetName, fmt.Sprintf("J%d", i+2), getLiquidity(j.GetLiquidity()))

		f.SetCellValue(SheetName, fmt.Sprintf("K%d", i+2), fmt.Sprintf("%.2f", float64(j.GetTradeVol())/10000))
		f.SetCellValue(SheetName, fmt.Sprintf("L%d", i+2), fmt.Sprintf("%.2f", float64(j.GetProfit())/10000))
		f.SetCellValue(SheetName, fmt.Sprintf("M%d", i+2), fmt.Sprintf("%.2f%%", j.GetProfitRate()))
		f.SetCellValue(SheetName, fmt.Sprintf("N%d", i+2), fmt.Sprintf("%.2f", float64(j.GetTotalFee())/10000))
		f.SetCellValue(SheetName, fmt.Sprintf("O%d", i+2), fmt.Sprintf("%.2f", float64(j.GetCrossFee())/10000))
		f.SetCellValue(SheetName, fmt.Sprintf("P%d", i+2), fmt.Sprintf("%.2f%%", j.GetCancelRate()))
		f.SetCellValue(SheetName, fmt.Sprintf("Q%d", i+2), j.GetMinSplitOrder())
		f.SetCellValue(SheetName, fmt.Sprintf("R%d", i+2), fmt.Sprintf("%.2f", j.GetDealEffi()))
		f.SetCellValue(SheetName, fmt.Sprintf("S%d", i+2), j.GetCreateTime())
	}
	f.SetActiveSheet(index)
	fn = l.GetFileName(1)
	fp = l.GetFilePath(fn)
	if err := f.SaveAs(fp); err != nil {
		l.Logger.Error("write excel error:", err)
	}
	return
}

func (l *UserProfileExportLogic) BuildProgressFile(info []*assessservice.ProfileInfo) (fn, fp string) {
	f := excelize.NewFile()
	index := f.NewSheet(SheetName)
	// 设置表头
	for k, v := range Charater {
		// A1, B1, C1,D1
		if k < len(ProgressHeader) {
			f.SetCellValue(SheetName, fmt.Sprintf("%s1", v), ProgressHeader[k])
		} else {
			break
		}
	}
	// A2,B2,C2,D2
	for i, j := range info {
		f.SetCellValue(SheetName, fmt.Sprintf("A%d", i+2), j.GetAccoutId())
		f.SetCellValue(SheetName, fmt.Sprintf("B%d", i+2), j.GetAccountName())
		f.SetCellValue(SheetName, fmt.Sprintf("C%d", i+2), j.GetProvider())
		f.SetCellValue(SheetName, fmt.Sprintf("D%d", i+2), j.GetAlgoName())
		f.SetCellValue(SheetName, fmt.Sprintf("E%d", i+2), j.GetAlgoOrderId())
		f.SetCellValue(SheetName, fmt.Sprintf("F%d", i+2), j.GetSecId())
		f.SetCellValue(SheetName, fmt.Sprintf("G%d", i+2), j.GetSecName())
		f.SetCellValue(SheetName, fmt.Sprintf("H%d", i+2), j.GetIndustry())
		f.SetCellValue(SheetName, fmt.Sprintf("I%d", i+2), getFundType(j.GetFundType()))
		f.SetCellValue(SheetName, fmt.Sprintf("J%d", i+2), getLiquidity(j.GetLiquidity()))

		f.SetCellValue(SheetName, fmt.Sprintf("K%d", i+2), fmt.Sprintf("%.2f%%", j.GetProgress()))
		f.SetCellValue(SheetName, fmt.Sprintf("L%d", i+2), fmt.Sprintf("%.2f", j.GetAlgoOrderFit()))
		f.SetCellValue(SheetName, fmt.Sprintf("M%d", i+2), fmt.Sprintf("%.2f", j.GetTradeVolFit()))
		f.SetCellValue(SheetName, fmt.Sprintf("N%d", i+2), j.GetCreateTime())
	}
	f.SetActiveSheet(index)
	fn = l.GetFileName(2)
	fp = l.GetFilePath(fn)
	if err := f.SaveAs(fp); err != nil {
		l.Logger.Error("write excel error:", err)
	}
	return
}

func (l *UserProfileExportLogic) BuildRiskFile(info []*assessservice.ProfileInfo) (fn, fp string) {
	f := excelize.NewFile()
	index := f.NewSheet(SheetName)
	// 设置表头
	for k, v := range Charater {
		// A1, B1, C1,D1
		if k < len(RiskHeader) {
			f.SetCellValue(SheetName, fmt.Sprintf("%s1", v), RiskHeader[k])
		} else {
			break
		}
	}
	// A2,B2,C2,D2
	for i, j := range info {
		f.SetCellValue(SheetName, fmt.Sprintf("A%d", i+2), j.GetAccoutId())
		f.SetCellValue(SheetName, fmt.Sprintf("B%d", i+2), j.GetAccountName())
		f.SetCellValue(SheetName, fmt.Sprintf("C%d", i+2), j.GetProvider())
		f.SetCellValue(SheetName, fmt.Sprintf("D%d", i+2), j.GetAlgoName())
		f.SetCellValue(SheetName, fmt.Sprintf("E%d", i+2), j.GetAlgoOrderId())
		f.SetCellValue(SheetName, fmt.Sprintf("F%d", i+2), j.GetSecId())
		f.SetCellValue(SheetName, fmt.Sprintf("G%d", i+2), j.GetSecName())
		f.SetCellValue(SheetName, fmt.Sprintf("H%d", i+2), j.GetIndustry())
		f.SetCellValue(SheetName, fmt.Sprintf("I%d", i+2), getFundType(j.GetFundType()))
		f.SetCellValue(SheetName, fmt.Sprintf("J%d", i+2), getLiquidity(j.GetLiquidity()))

		f.SetCellValue(SheetName, fmt.Sprintf("K%d", i+2), fmt.Sprintf("%.2f%%", j.GetMinJointRate()))
		f.SetCellValue(SheetName, fmt.Sprintf("L%d", i+2), fmt.Sprintf("%.2f%%", j.GetProfitRate()))
		f.SetCellValue(SheetName, fmt.Sprintf("M%d", i+2), fmt.Sprintf("%.2f%%", j.GetWithdrawRate()))
		f.SetCellValue(SheetName, fmt.Sprintf("N%d", i+2), j.GetCreateTime())
	}
	f.SetActiveSheet(index)
	fn = l.GetFileName(3)
	fp = l.GetFilePath(fn)
	if err := f.SaveAs(fp); err != nil {
		l.Logger.Error("write excel error:", err)
	}
	return
}

func (l *UserProfileExportLogic) BuildAssessFile(info []*assessservice.ProfileInfo) (fn, fp string) {
	f := excelize.NewFile()
	index := f.NewSheet(SheetName)
	// 设置表头
	for k, v := range Charater {
		// A1, B1, C1,D1
		if k < len(AssessHeader) {
			f.SetCellValue(SheetName, fmt.Sprintf("%s1", v), AssessHeader[k])
		} else {
			break
		}
	}
	// A2,B2,C2,D2
	for i, j := range info {
		f.SetCellValue(SheetName, fmt.Sprintf("A%d", i+2), j.GetAccoutId())
		f.SetCellValue(SheetName, fmt.Sprintf("B%d", i+2), j.GetAccountName())
		f.SetCellValue(SheetName, fmt.Sprintf("C%d", i+2), j.GetProvider())
		f.SetCellValue(SheetName, fmt.Sprintf("D%d", i+2), j.GetAlgoName())
		f.SetCellValue(SheetName, fmt.Sprintf("E%d", i+2), j.GetAlgoOrderId())
		f.SetCellValue(SheetName, fmt.Sprintf("F%d", i+2), j.GetSecId())
		f.SetCellValue(SheetName, fmt.Sprintf("G%d", i+2), j.GetSecName())
		f.SetCellValue(SheetName, fmt.Sprintf("H%d", i+2), j.GetIndustry())
		f.SetCellValue(SheetName, fmt.Sprintf("I%d", i+2), getFundType(j.GetFundType()))
		f.SetCellValue(SheetName, fmt.Sprintf("J%d", i+2), getLiquidity(j.GetLiquidity()))

		f.SetCellValue(SheetName, fmt.Sprintf("K%d", i+2), fmt.Sprintf("%.4f", j.GetVwapDev()))
		f.SetCellValue(SheetName, fmt.Sprintf("L%d", i+2), fmt.Sprintf("%.2f%%", j.GetAssessProfitRate()))
		f.SetCellValue(SheetName, fmt.Sprintf("M%d", i+2), j.GetCreateTime())
	}
	f.SetActiveSheet(index)
	fn = l.GetFileName(4)
	fp = l.GetFilePath(fn)
	if err := f.SaveAs(fp); err != nil {
		l.Logger.Error("write excel error:", err)
	}
	return
}

func (l *UserProfileExportLogic) BuildStabilityFile(info []*assessservice.ProfileInfo) (fn, fp string) {
	f := excelize.NewFile()
	index := f.NewSheet(SheetName)
	// 设置表头
	for k, v := range Charater {
		// A1, B1, C1,D1
		if k < len(StabilityHeader) {
			f.SetCellValue(SheetName, fmt.Sprintf("%s1", v), StabilityHeader[k])
		} else {
			break
		}
	}
	// A2,B2,C2,D2
	for i, j := range info {
		f.SetCellValue(SheetName, fmt.Sprintf("A%d", i+2), j.GetAccoutId())
		f.SetCellValue(SheetName, fmt.Sprintf("B%d", i+2), j.GetAccountName())
		f.SetCellValue(SheetName, fmt.Sprintf("C%d", i+2), j.GetProvider())
		f.SetCellValue(SheetName, fmt.Sprintf("D%d", i+2), j.GetAlgoName())
		f.SetCellValue(SheetName, fmt.Sprintf("E%d", i+2), j.GetAlgoOrderId())
		f.SetCellValue(SheetName, fmt.Sprintf("F%d", i+2), j.GetSecId())
		f.SetCellValue(SheetName, fmt.Sprintf("G%d", i+2), j.GetSecName())
		f.SetCellValue(SheetName, fmt.Sprintf("H%d", i+2), j.GetIndustry())
		f.SetCellValue(SheetName, fmt.Sprintf("I%d", i+2), getFundType(j.GetFundType()))
		f.SetCellValue(SheetName, fmt.Sprintf("J%d", i+2), getLiquidity(j.GetLiquidity()))

		f.SetCellValue(SheetName, fmt.Sprintf("K%d", i+2), fmt.Sprintf("%.2f", j.GetVwapStdDev()))
		f.SetCellValue(SheetName, fmt.Sprintf("L%d", i+2), fmt.Sprintf("%.2f", j.GetPfRateStdDev()))
		f.SetCellValue(SheetName, fmt.Sprintf("M%d", i+2), fmt.Sprintf("%.2f%%", j.GetMinJointRate()))
		f.SetCellValue(SheetName, fmt.Sprintf("N%d", i+2), fmt.Sprintf("%.2f", j.GetTradeVolFitStdDev()))
		f.SetCellValue(SheetName, fmt.Sprintf("O%d", i+2), fmt.Sprintf("%.2f", j.GetTimeFitStdDev()))
		f.SetCellValue(SheetName, fmt.Sprintf("P%d", i+2), j.GetCreateTime())
	}
	f.SetActiveSheet(index)
	fn = l.GetFileName(5)
	fp = l.GetFilePath(fn)
	if err := f.SaveAs(fp); err != nil {
		l.Logger.Error("write excel error:", err)
	}
	return
}

func (l *UserProfileExportLogic) GetFileName(profile int32) string {
	var fileName string
	timeStamp := time.Now().Format("200601021504")
	switch profile {
	case 1:
		fileName = "economy_" + timeStamp
	case 2:
		fileName = "progress_" + timeStamp
	case 3:
		fileName = "risk_" + timeStamp
	case 4:
		fileName = "assess_" + timeStamp
	case 5:
		fileName = "stability_" + timeStamp
	}
	return fileName + ".xlsx"
}

func (l *UserProfileExportLogic) GetFilePath(fileName string) string {
	pwd, _ := os.Getwd()
	if runtime.GOOS == "linux" {
		pwd += "/export/"
	} else if runtime.GOOS == "windows" {
		pwd += `\export\`
	}
	err := os.MkdirAll(pwd, 0766)
	if err != nil {
		l.Logger.Error("create dir err:", err)
		return ""
	}
	pwd += fileName
	l.Logger.Info("get file:", pwd)
	return pwd
}

// getFundType 市值类型转成字符串形式
func getFundType(fundType int32) string {
	switch fundType {
	case 1:
		return "超大"
	case 2:
		return "大"
	case 3:
		return "中等"
	case 4:
		return "小"
	default:
		return "-"
	}
}

// getLiquidity 流动性转成字符串类型
func getLiquidity(liq int32) string {
	switch liq {
	case 1:
		return "高"
	case 2:
		return "中"
	case 3:
		return "低"
	default:
		return "-"
	}
}
