package logic

import (
	"algo_assess/assess-api-server/internal/svc"
	"algo_assess/assess-api-server/internal/types"
	"algo_assess/assess-rpc-server/assessservice"
	"algo_assess/global"
	"algo_assess/pkg/tools"
	"context"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"time"
)

type AlgoFixLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAlgoFixLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlgoFixLogic {
	return &AlgoFixLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// AlgoFix 母单信息上传
func (l *AlgoFixLogic) AlgoFix(r *http.Request) (resp *types.BaseRsp, err error) {
	file, header, err := r.FormFile("file")
	if err != nil {
		l.Logger.Error("parse form file error:", err)
		return &types.BaseRsp{
			Code: 10000,
			Msg:  "文件解析失败",
		}, nil
	}
	// 判断文件名后缀是什么格式
	fileSuffix := path.Ext(header.Filename)
	l.Logger.Info("get filename:", header.Filename, ", suffix:", fileSuffix)

	var sourceFrom int32 = global.SourceFromFix
	var batchNo int64
	fileType := r.FormValue("fileType")
	l.Logger.Info("get fileTYpe :", fileType)
	//fileType := r.Header.Get("fileType") // 头里加入key,区分是数据修复 还是订单导入
	if fileType == "import" { // 订单导入
		sourceFrom = global.SourceFromOrigin
		batchNo = tools.GeneralID()
	}

	list := make([]*assessservice.AlgoOrderPerf, 0)
	if fileSuffix == ".csv" {
		reader := csv.NewReader(file)
		reader.Comma = ','
		reader.FieldsPerRecord = -1
		reader.LazyQuotes = true
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
			//l.Logger.Info("get line :", line, len(line))
			// 这里要兼容绩效通用导入模板和总线的导入模板，单行10个字段为通用模板，单行超过30个字段是总线模板
			if len(line) < 10 {
				l.Logger.Info("format error, line:", line)
				continue
			} else if len(line) == 10 { // 通用模板
				algoOrderPref := &assessservice.AlgoOrderPerf{
					BatchNo: batchNo,
					//BatchName:     batchName,
					Id:            cast.ToUint32(strings.TrimSpace(line[1])), // 母单号
					BasketId:      cast.ToUint32(strings.TrimSpace(line[0])),
					AlgorithmType: cast.ToUint32(strings.TrimSpace(line[3])),
					AlgorithmId:   cast.ToUint32(strings.TrimSpace(line[2])),
					USecurityId:   0,
					SecurityId:    strings.TrimSpace(line[5]),
					AlgoOrderQty:  cast.ToUint64(strings.TrimSpace(line[6])),
					TransactTime:  cast.ToUint64(line[7]),
					StartTime:     cast.ToUint64(line[8]),
					EndTime:       cast.ToUint64(line[9]),
					BusUserId:     strings.TrimSpace(line[4]),
					SourceFrom:    sourceFrom,
				}
				//l.Logger.Infof("get algoOrderPref:%+v", algoOrderPref)
				list = append(list, algoOrderPref)
			} else if len(line) > 25 { // 总线模板
				st, _ := time.ParseInLocation(global.TimeFormat, strings.TrimSpace(line[21]), time.Local)
				t := st.UnixMicro()
				algoOrderPref := &assessservice.AlgoOrderPerf{
					BatchNo: batchNo,
					//BatchName:     batchName,
					Id:            cast.ToUint32(strings.TrimSpace(line[0])),
					BasketId:      cast.ToUint32(strings.TrimSpace(line[2])),
					AlgorithmType: cast.ToUint32(strings.TrimSpace(line[4])),
					AlgorithmId:   cast.ToUint32(strings.TrimSpace(line[5])),
					USecurityId:   0,
					SecurityId:    strings.TrimSpace(line[6]),
					AlgoOrderQty:  cast.ToUint64(strings.TrimSpace(line[9])),
					TransactTime:  cast.ToUint64(t),
					StartTime:     34200,
					EndTime:       53820,
					BusUserId:     strings.TrimSpace(line[18]),
					SourceFrom:    sourceFrom,
				}
				//l.Logger.Infof("get algoOrderPref:%+v", algoOrderPref)
				list = append(list, algoOrderPref)
			}

		}
	} else if fileSuffix == ".xml" {
		// xml格式解析
		b, err := ioutil.ReadAll(file)
		if err != nil {
			l.Logger.Error(" read file error:", err)
			return &types.BaseRsp{
				Code: 10000,
				Msg:  "文件上传错误",
			}, nil
		}
		var newAlgoOrder NewAlgoOrder
		xmlErr := xml.Unmarshal(b, &newAlgoOrder)
		if xmlErr != nil {
			l.Logger.Error("Unmarshal error:", err)
			return &types.BaseRsp{
				Code: 10000,
				Msg:  "文件解析错误",
			}, nil
		}
		for _, algo := range newAlgoOrder.NewAlgoOrderData {
			algoOrderPref := &assessservice.AlgoOrderPerf{
				BatchNo: batchNo,
				//BatchName:     batchName,
				Id:            algo.Id,
				BasketId:      algo.BasketId,
				AlgorithmType: uint32(algo.AlgorithmType),
				AlgorithmId:   uint32(algo.AlgorithmId),
				USecurityId:   algo.USecurityId,
				SecurityId:    algo.SecurityId,
				AlgoOrderQty:  algo.AlgoOrderQty,
				TransactTime:  algo.TransactTime,
				StartTime:     algo.StartTime,
				EndTime:       algo.EndTime,
				BusUserId:     algo.BusUserId,
				SourceFrom:    sourceFrom,
			}
			list = append(list, algoOrderPref)
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
	importReq := assessservice.AlgoOrderPerfs{
		Parts: list,
	}
	if sourceFrom == global.SourceFromOrigin { // 订单导入|| 原始订单
		importRsp, err := l.svcCtx.AssessClient.ImportAlgoOrdr(l.ctx, &importReq)
		if err != nil {
			l.Logger.Error("rpc call ImportAlgoOrdr error:", err)
			return &types.BaseRsp{
				Code: 10000,
				Msg:  "母单导入失败",
			}, nil
		}
		resp = &types.BaseRsp{
			Code:      int(importRsp.GetCode()),
			Msg:       importRsp.GetMsg(),
			BatchNo:   batchNo, // 母单产生批次号，同步给子单
			StartTime: importRsp.GetStartTime(),
			EndTime:   importRsp.GetEndTime(),
		}

	} else { // 母单修复导入
		rsp, err := l.svcCtx.AssessClient.PushAlgoOrder(l.ctx, &importReq)
		if err != nil {
			l.Logger.Error("rpc PushAlgoOrder error:", err)
			return &types.BaseRsp{
				Code: 10000,
				Msg:  "母单推送失败",
			}, nil
		}
		resp = &types.BaseRsp{
			Code: int(rsp.GetCode()),
			Msg:  rsp.GetMsg(),
		}
	}

	return resp, nil
}

type NewAlgoOrder struct {
	NewAlgoOrderData []types.AlgoOrder `xml:"NewAlgoOrderData"`
}
