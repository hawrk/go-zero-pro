package logic

import (
	"algo_assess/assess-rpc-server/assessservice"
	"algo_assess/assess-rpc-server/proto"
	"algo_assess/global"
	"algo_assess/pkg/tools"
	"bytes"
	"context"
	"encoding/csv"
	"github.com/spf13/cast"
	"net/http"
	"time"

	"algo_assess/assess-api-server/internal/svc"
	"algo_assess/assess-api-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExportTradeOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExportTradeOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExportTradeOrderLogic {
	return &ExportTradeOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExportTradeOrderLogic) ExportTradeOrder(req *types.TradeOrderReq, w http.ResponseWriter) error {
	l.Logger.Infof("get req:%+v", req)
	var fileName string
	var content [][]string
	var start, end int32
	if req.StartTime == 0 || req.EndTime == 0 {
		start = 0
		end = 0
	} else {
		start = cast.ToInt32(time.Unix(req.StartTime, 0).Format(global.TimeFormatDay))
		end = cast.ToInt32(time.Unix(req.EndTime, 0).Format(global.TimeFormatDay))
	}
	if req.ExportType == 1 { // 导出母单
		in := &assessservice.ReqQueryAlgoOrder{
			PageId:    1,
			PageNum:   100000000,
			StartTime: start,
			EndTime:   end,
			UserId:    req.UserId,
			AlgoId:    int32(req.AlgoOrderId),
		}
		aRsp, err := l.svcCtx.AssessClient.QueryAlgoOrder(l.ctx, in)
		if err != nil {
			l.Logger.Error("rpc QueryAlgoOrder error:", err)
			return err
		}
		fileName, content = l.BuildAlgoOrderCSV(aRsp.GetParts())

	} else if req.ExportType == 2 { // 导出子单
		cReq := &assessservice.ReqQueryChildOrder{
			PageId:       1,
			PageNum:      100000000,
			StartTime:    start,
			EndTime:      end,
			UserId:       req.UserId,
			AlgoOrderId:  req.AlgoOrderId,
			ChildOrderId: req.ChildOrderId,
		}
		cRsp, err := l.svcCtx.AssessClient.QueryChildOrder(l.ctx, cReq)
		if err != nil {
			l.Logger.Error("rpc QueryChildOrder error:", err)
			return err
		}
		fileName, content = l.BuildChildOrderCSV(cRsp.GetParts())
	}

	header := w.Header()
	header.Add("Content-Type", "application/octet-stream")
	header.Add("Content-Disposition", "filename="+fileName)
	buff := new(bytes.Buffer)
	buff.WriteString("\xEF\xBB\xBF")
	wStr := csv.NewWriter(buff)
	wStr.WriteAll(content)
	wStr.Flush()
	_, _ = w.Write(buff.Bytes())

	return nil
}

func (l *ExportTradeOrderLogic) BuildAlgoOrderCSV(in []*proto.AlgoOrder) (string, [][]string) {
	var content [][]string
	// set header
	content = append(content, global.AlgoHeader)
	// set content
	curDay := time.Now().Format(global.TimeFormatDay)
	for _, v := range in {
		sli := make([]string, 0, 10)
		sli = append(sli, cast.ToString(v.BasketId), cast.ToString(v.AlgoId), cast.ToString(v.AlgorithmId), cast.ToString(v.AlgorithmType),
			v.UserId, v.SecId, cast.ToString(v.AlgoOrderQty*100), l.getTransTime(curDay, v.TransTime, cast.ToString(v.UnixTime)), v.StartTime, v.EndTime)
		content = append(content, sli)
	}
	return "algo_order", content
}

func (l *ExportTradeOrderLogic) BuildChildOrderCSV(in []*proto.ChildOrder) (string, [][]string) {
	var content [][]string
	content = append(content, global.ChildOrderHeader)
	curDay := time.Now().Format(global.TimeFormatDay)
	for _, v := range in {
		sli := make([]string, 0, 18)
		sli = append(sli, cast.ToString(v.ChildOrderId), cast.ToString(v.AlgoOrderId), cast.ToString(v.AlgorithmId), cast.ToString(v.AlgorithmType),
			v.UserId, v.SecurityId, cast.ToString(v.TradeSide), cast.ToString(v.OrderQty*100), cast.ToString(v.Price), cast.ToString(v.OrderType),
			cast.ToString(v.LastPx), cast.ToString(v.LastQty*100), cast.ToString(v.ComQty*100), cast.ToString(v.TotalFee), cast.ToString(v.ArrivedPrice),
			cast.ToString(v.OrdStatus), l.getTransTime(curDay, v.TransactAt, v.TransactTime))
		content = append(content, sli)
	}
	return "child_order", content
}

func (l *ExportTradeOrderLogic) getTransTime(curDay, curTransTime string, unixTime string) string {
	if l.svcCtx.Config.WorkControl.EnableFakeDay {
		newDay := curDay + curTransTime[8:12]
		return cast.ToString(tools.TimeStr2TimeMicro(newDay))
	}
	return unixTime
}
