package logic

import (
	mqservice "algo_assess/assess-mq-server/assessmqservice"
	"algo_assess/assess-rpc-server/assessservice"
	"algo_assess/global"
	mkservice "algo_assess/market-mq-server/marketservice"
	"algo_assess/pkg/tools"
	"context"
	"github.com/spf13/cast"
	"time"

	"algo_assess/assess-api-server/internal/svc"
	"algo_assess/assess-api-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GeneralLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGeneralLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GeneralLogic {
	return &GeneralLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GeneralLogic) General(req *types.GeneralReq) (resp *types.GeneralRsp, err error) {
	// 可以打印指针类型的结构体，但输出信息比较多
	//l.Logger.Info("spew out:", spew.Sdump(req))

	gRspCh := make(chan *assessservice.GeneralRsp)
	mRspCh := make(chan *mqservice.GeneralRsp)
	qRspCh := make(chan *mkservice.MarkRsp)

	// 后面所有服务都是根据这个时间格式来计算
	start := cast.ToInt64(time.Unix(req.StartTime, 0).Format(global.TimeFormatMinInt))
	end := cast.ToInt64(time.Unix(req.EndTime, 0).Format(global.TimeFormatMinInt))
	var userId string
	if req.UserId == "" {
		userId = "0"
	} else {
		userId = req.UserId
	}

	go func() {
		gReq := &assessservice.GeneralReq{
			AlgoId:          req.AlgoId,
			SecId:           req.SecId,
			UserId:          userId,
			TimeDemension:   req.TimeDemension,
			OrderStatusType: req.OrderStatusType,
			StartTime:       start,
			EndTime:         end,
		}
		rsp, err := l.svcCtx.AssessClient.GetGeneral(l.ctx, gReq)
		if err != nil {
			l.Logger.Error("assess rpc call error :", err)
			gRspCh <- &assessservice.GeneralRsp{}
			return
		}
		gRspCh <- rsp
	}()

	go func() {
		mReq := &mqservice.GeneralReq{
			AlgoId:          req.AlgoId,
			SecId:           req.SecId,
			UserId:          userId,
			TimeDemension:   req.TimeDemension,
			OrderStatusType: req.OrderStatusType,
			StartTime:       start,
			EndTime:         end,
		}
		rsp, err := l.svcCtx.AssessMQClient.GetMqGeneral(l.ctx, mReq)
		if err != nil {
			l.Logger.Error("assess mq rpc call error:", err)
			mRspCh <- &mqservice.GeneralRsp{}
			return
		}
		mRspCh <- rsp
	}()

	go func() {
		qReq := &mkservice.MarkReq{
			SecId:     req.SecId,
			StartTime: start,
			EndTime:   end,
			SecSource: req.SecSource,
		}
		rsp, err := l.svcCtx.MarketMQClient.GetMarketInfo(l.ctx, qReq)
		if err != nil {
			l.Logger.Error("market mq rpc call error:", err)
			qRspCh <- &mkservice.MarkRsp{}
			return
		}
		qRspCh <- rsp
	}()

	gRsp, mRsp, qRsp := <-gRspCh, <-mRspCh, <-qRspCh
	data := l.BuildGeneralRsp(req.StartTime, req.EndTime, gRsp, mRsp, qRsp)
	//l.Logger.Infof("get rsp: %+v", data)

	p := &types.GeneralRsp{
		Code: 200,
		Msg:  "success",
		Data: data,
	}
	return p, nil
}

func (l *GeneralLogic) BuildGeneralRsp(startTime, endTime int64,
	grsp *assessservice.GeneralRsp, mrsp *mqservice.GeneralRsp, qRsp *mkservice.MarkRsp) []types.GeneralData {
	data := make([]types.GeneralData, 0, global.MaxAssessTimeLen)
	m := make(map[string]types.GeneralData) // key -> 10:00   val -> rsp
	var tmpProgress float64                 // 用于处理成交进度， 当前时间点无交易，需要取上一个有交易的点填充
	//TODO: 并发拼接
	MakeAssessMqRsp(m, mrsp, qRsp)  // 拼还在计算中的缓存数据
	MakeAssessRPCRsp(m, grsp, qRsp) // 拼已落地DB的数据
	// 填充无交易数据的时间点
	end := cast.ToInt64(time.Unix(endTime, 0).Format(global.TimeFormatMinInt))
	for i := 0; i <= global.MaxAssessTimeLen; i++ {
		start := time.Unix(startTime, 0).Add(time.Minute * time.Duration(i)).Format(global.TimeFormatMinInt)
		// 正常交易时间 9:30 - 11:30    13:00-15:00, 需要剔掉中间不交易的时间
		tt := cast.ToInt(start[8:])
		if tt > 1130 {
			start = time.Unix(startTime, 0).Add(time.Minute * time.Duration(i+90)).Format(global.TimeFormatMinInt)
			//l.Logger.Info("start:", start)
		}
		if cast.ToInt64(start) > end {
			break
		}
		timePoint := tools.GetTimePoint(start)
		//l.Logger.Info("get time point:", timePoint)
		if _, exist := m[timePoint]; exist { // 当前时间点有交易数据
			// 这里需要再判断一下，如果有交易的数据，但是没有行情的数据，也不返回该时间点，否则前端会展示异常
			if m[timePoint].MkLastPrice <= 0.00 {
				l.Logger.Info("filter: timepoint:", timePoint, " has no market data, continue")
				continue
			}
			l.Logger.Infof("get timepoint :%s, assess:%+v", timePoint, m[timePoint])
			data = append(data, m[timePoint])
			tmpProgress = m[timePoint].DealProgress
		} else { // 无交易数据，填充0
			if qRsp.Attrs[cast.ToInt64(start)].GetLastPrice() <= 0 { // 行情都没有数据时，不返回该时间节点
				continue
			}
			data = append(data, MakeEmptyDataRsp(start, qRsp, tmpProgress))
		}
	}
	return data
}
