package logic

import (
	"algo_assess/models"
	"context"
	"strconv"

	"algo_assess/assess-rpc-server/internal/svc"
	"algo_assess/assess-rpc-server/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryChildOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryChildOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryChildOrderLogic {
	return &QueryChildOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  子单信息
func (l *QueryChildOrderLogic) QueryChildOrder(in *proto.ReqQueryChildOrder) (*proto.RespQueryChildOrder, error) {
	var algorithmId int32
	if in.GetAlgoName() != "" {
		a, err := l.svcCtx.AlgoInfoRepo.GetAlgoIdByAlgoName(l.ctx, in.GetAlgoName())
		if err != nil {
			l.Logger.Error("GetAlgoIdByAlgoName error:", err)
		}
		algorithmId = a
	}
	// 查子单修复表的数据
	count, orders, err := l.svcCtx.OrderDetailRepo.QueryChildOrder(l.ctx, in, algorithmId)
	if err != nil {
		return &proto.RespQueryChildOrder{
			Code:  10000,
			Msg:   "查询失败",
			Total: count,
		}, err
	}
	// 查算法基础表
	algo, err := l.svcCtx.AlgoInfoRepo.GetAlgoBase(l.ctx)
	if err != nil {
		l.Logger.Error("GetAlgoBase:", err)
		return &proto.RespQueryChildOrder{
			Code:  10000,
			Msg:   "查询失败",
			Total: count,
		}, nil
	}
	mAlgo := make(map[int]string) // key -> algoId  value-> algoName
	for _, v := range algo {
		mAlgo[v.AlgoId] = v.AlgoName
	}
	id := (in.GetPageId() - 1) * in.GetPageNum()
	return &proto.RespQueryChildOrder{
		Code:  200,
		Msg:   "查询成功",
		Total: count,
		Parts: func(orders []*models.TbAlgoOrderDetail) (ret []*proto.ChildOrder) {
			for _, o := range orders {
				id++
				a := &proto.ChildOrder{
					Id:            uint64(id),
					Date:          int32(o.Date),
					BatchNo:       o.BatchNo,
					BatchName:     o.BatchName,
					ChildOrderId:  o.ChildOrderId,
					AlgoOrderId:   uint32(o.AlgoOrderId),
					AlgorithmType: uint32(o.AlgorithmType),
					AlgorithmId:   uint32(o.AlgorithmId),
					AlgoName:      mAlgo[int(o.AlgorithmId)],
					UserId:        o.UserId,
					UsecurityId:   uint32(o.UsecurityId),
					SecurityId:    o.SecurityId,
					TradeSide:     int32(o.TradeSide),
					OrderQty:      o.OrderQty,
					Price:         o.Price,
					OrderType:     int32(o.OrderType),
					LastPx:        o.LastPx,
					LastQty:       o.LastQty,
					ComQty:        o.ComQty,
					ArrivedPrice:  o.ArrivedPrice,
					TotalFee:      o.TotalFee / 10000,
					OrdStatus:     int32(o.OrdStatus),
					TransactTime:  strconv.FormatInt(o.TransactTime, 10),
					TransactAt:    strconv.FormatInt(o.TransactAt, 10),
					ProcStatus:    int32(o.ProcStatus),
					FixFlag:       int32(o.Source),
					CreateTime:    o.CreateTime.Format("2006-01-02 15:04:05"),
				}
				ret = append(ret, a)
			}
			return ret
		}(orders),
	}, nil
}
