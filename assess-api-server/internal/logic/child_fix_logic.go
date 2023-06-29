package logic

import (
	"algo_assess/assess-rpc-server/assessservice"
	"algo_assess/global"
	"context"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"github.com/spf13/cast"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"time"

	"algo_assess/assess-api-server/internal/svc"
	"algo_assess/assess-api-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChildFixLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChildFixLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChildFixLogic {
	return &ChildFixLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChildFixLogic) ChildFix(r *http.Request) (resp *types.BaseRsp, err error) {
	file, header, err := r.FormFile("file")
	if err != nil {
		l.Logger.Error("form file error:", err)
		return &types.BaseRsp{
			Code: 10000,
			Msg:  "文件解析失败",
		}, nil
	}

	var sourceFrom int32 = global.SourceFromFix
	var batchNo int64
	fileType := r.FormValue("fileType")
	strBatchNo := r.FormValue("batchNo")
	l.Logger.Info("get fileTYpe :", fileType)
	//fileType := r.Header.Get("fileType") // 头里加入key,区分是数据修复 还是订单导入
	if fileType == "import" { // 订单导入
		sourceFrom = global.SourceFromOrigin
		batchNo = cast.ToInt64(strBatchNo)
	}

	list := make([]*assessservice.ChildOrderPerf, 0)
	// 判断文件名后缀格式
	fileSuffix := path.Ext(header.Filename)
	l.Logger.Info("get filename:", header.Filename, ", suffix:", fileSuffix)
	if fileSuffix == ".csv" {
		reader := csv.NewReader(file)
		lines, err := reader.ReadAll()
		if err != nil {
			l.Logger.Error("read file error:", err)
			return &types.BaseRsp{
				Code: 10000,
				Msg:  "文件解析失败",
			}, nil
		}
		for k, line := range lines {
			if k == 0 { // 跳过页眉
				continue
			}
			//l.Logger.Info("get line :", line)
			// 这里要兼容绩效通用导入模板和总线的导入模板，单行10个字段为通用模板，单行超过30个字段是总线模板
			if len(line) < 17 {
				l.Logger.Info("format error, line:", line)
				continue
			} else if len(line) == 17 { // 通用模板
				pref := &assessservice.ChildOrderPerf{
					BatchNo: batchNo,
					//BatchName:      batchName,
					Id:             cast.ToUint32(strings.TrimSpace(line[0])),
					BusUserId:      strings.TrimSpace(line[4]),
					BusUuserId:     0,
					AlgoOrderId:    cast.ToUint32(strings.TrimSpace(line[1])),
					AlgorithmType:  cast.ToUint32(strings.TrimSpace(line[3])),
					AlgorithmId:    cast.ToUint32(strings.TrimSpace(line[2])),
					USecurityId:    0,
					SecurityId:     strings.TrimSpace(line[5]),
					Side:           cast.ToUint32(strings.TrimSpace(line[6])),
					OrderQty:       cast.ToUint64(strings.TrimSpace(line[7])),
					Price:          cast.ToUint64(strings.TrimSpace(line[8])),
					OrderType:      cast.ToUint32(strings.TrimSpace(line[9])),
					CumQty:         cast.ToUint64(strings.TrimSpace(line[12])),
					LastPx:         cast.ToUint64(strings.TrimSpace(line[10])),
					LastQty:        cast.ToUint64(strings.TrimSpace(line[11])),
					Charge:         cast.ToFloat64(strings.TrimSpace(line[13])),
					ArrivedPrice:   cast.ToUint64(strings.TrimSpace(line[14])),
					ChildOrdStatus: cast.ToUint32(strings.TrimSpace(line[15])),
					TransactTime:   cast.ToUint64(line[16]),
					SourceFrom:     sourceFrom,
				}
				//l.Logger.Infof("get ChildOrderPerf:%+v", pref)
				list = append(list, pref)
			} else if len(line) >= 35 { // 总线模板
				st, _ := time.ParseInLocation(global.TimeFormat, strings.TrimSpace(line[28]), time.Local)
				t := st.UnixMicro()
				pref := &assessservice.ChildOrderPerf{
					BatchNo: batchNo,
					//BatchName:      batchName,
					Id:             cast.ToUint32(strings.TrimSpace(line[0])),
					BusUserId:      strings.TrimSpace(line[1]),
					BusUuserId:     0,
					AlgoOrderId:    cast.ToUint32(strings.TrimSpace(line[5])),
					AlgorithmType:  cast.ToUint32(strings.TrimSpace(line[3])),
					AlgorithmId:    cast.ToUint32(strings.TrimSpace(line[4])),
					USecurityId:    0,
					SecurityId:     strings.TrimSpace(line[6]),
					Side:           cast.ToUint32(strings.TrimSpace(line[12])),
					OrderQty:       cast.ToUint64(strings.TrimSpace(line[8])),
					Price:          cast.ToUint64(strings.TrimSpace(line[9])),
					OrderType:      cast.ToUint32(strings.TrimSpace(line[11])),
					CumQty:         cast.ToUint64(strings.TrimSpace(line[15])),
					LastPx:         cast.ToUint64(strings.TrimSpace(line[16])),
					LastQty:        cast.ToUint64(strings.TrimSpace(line[15])),
					Charge:         cast.ToFloat64(strings.TrimSpace(line[34])),
					ArrivedPrice:   cast.ToUint64(strings.TrimSpace(line[33])),
					ChildOrdStatus: cast.ToUint32(strings.TrimSpace(line[21])),
					TransactTime:   cast.ToUint64(t),
					SourceFrom:     sourceFrom,
				}
				//l.Logger.Infof("get ChildOrderPerf:%+v", pref)
				list = append(list, pref)
			}
		}
	} else if fileSuffix == ".xml" {
		b, err := ioutil.ReadAll(file)
		if err != nil {
			l.Logger.Error("read file error:", err)
			return &types.BaseRsp{
				Code: 10000,
				Msg:  "文件读取错误",
			}, nil
		}
		var newTradeOrder NewTradeOrder
		xmlErr := xml.Unmarshal(b, &newTradeOrder)
		if xmlErr != nil {
			l.Logger.Error("Unmarshal xml error:", err)
			return &types.BaseRsp{
				Code: 10000,
				Msg:  "文件解析错误",
			}, nil
		}
		for _, child := range newTradeOrder.NewTradeOrderData {
			pref := &assessservice.ChildOrderPerf{
				BatchNo: batchNo,
				//BatchName:      batchName,
				Id:             child.Id,
				BusUserId:      child.BusUserId,
				BusUuserId:     child.BusUuserId,
				AlgoOrderId:    child.AlgoOrderId,
				AlgorithmType:  uint32(child.AlgorithmType),
				AlgorithmId:    uint32(child.AlgorithmId),
				USecurityId:    child.USecurityId,
				SecurityId:     child.SecurityId,
				Side:           uint32(child.Side),
				OrderQty:       child.OrderQty,
				Price:          uint64(child.Price),
				OrderType:      uint32(child.OrderQty),
				CumQty:         child.CumQty,
				LastPx:         uint64(child.LastPx),
				LastQty:        child.LastQty,
				Charge:         child.Charge,
				ArrivedPrice:   uint64(child.ArrivedPrice),
				ChildOrdStatus: uint32(child.ChildOrdStatus),
				TransactTime:   child.TransactTime,
				SourceFrom:     sourceFrom,
			}
			list = append(list, pref)
		}
	} else {
		l.Logger.Error("unsupported file format:", header.Filename)
		return &types.BaseRsp{
			Code: 10000,
			Msg:  fmt.Sprintf("暂不支持该文件类型：%s", fileSuffix),
		}, nil
	}
	if len(list) == 0 {
		return &types.BaseRsp{
			Code: 10000,
			Msg:  "文件解析，合规的数据为0",
		}, nil
	}

	in := &assessservice.ChildOrderPerfs{
		Parts: list,
	}
	if sourceFrom == global.SourceFromOrigin { // 订单导入||原始订单
		rsp, err := l.svcCtx.AssessClient.ImportChildOrder(l.ctx, in)
		if err != nil {
			l.Logger.Error("rpc call ImportChildOrder error:", err)
			return &types.BaseRsp{
				Code: 10000,
				Msg:  "子单导入失败",
			}, nil
		}
		resp = &types.BaseRsp{
			Code: int(rsp.GetCode()),
			Msg:  rsp.GetMsg(),
		}

	} else { // 子单修复
		rsp, err := l.svcCtx.AssessClient.PushChildOrder(l.ctx, in)
		if err != nil {
			l.Logger.Error("rpc PushChildOrder error:", err)
			return &types.BaseRsp{
				Code: 10000,
				Msg:  "子单推送失败",
			}, nil
		}
		resp = &types.BaseRsp{
			Code: int(rsp.GetCode()),
			Msg:  rsp.GetMsg(),
		}
	}

	return resp, nil
}

type NewTradeOrder struct {
	NewTradeOrderData []types.ChildOrder `xml:"NewTradeOrderData"`
}
