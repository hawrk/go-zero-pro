package logic

import (
	"algo_assess/models"
	"context"
	"strconv"

	"algo_assess/assess-rpc-server/internal/svc"
	"algo_assess/assess-rpc-server/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryAlgoOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryAlgoOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryAlgoOrderLogic {
	return &QueryAlgoOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  母单信息
func (l *QueryAlgoOrderLogic) QueryAlgoOrder(in *proto.ReqQueryAlgoOrder) (*proto.RespQueryAlgoOrder, error) {
	var algorithmId int32
	if in.GetAlgoName() != "" { // 算法名称不为空，需要返查一下算法ID
		a, err := l.svcCtx.AlgoInfoRepo.GetAlgoIdByAlgoName(l.ctx, in.GetAlgoName())
		if err != nil {
			l.Logger.Error("GetAlgoIdByAlgoName error:", err)
		}
		algorithmId = a
	}
	// 查母单修复表的数据
	count, orders, err := l.svcCtx.AlgoOrderRepo.QueryAlgoOrder(l.ctx, in, algorithmId)
	if err != nil {
		return &proto.RespQueryAlgoOrder{
			Code:  10000,
			Msg:   "查询失败",
			Total: count,
		}, nil
	}
	// 查算法基础表
	algo, err := l.svcCtx.AlgoInfoRepo.GetAlgoBase(l.ctx)
	if err != nil {
		l.Logger.Error("GetAlgoBase:", err)
		return &proto.RespQueryAlgoOrder{
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
	return &proto.RespQueryAlgoOrder{
		Code:  200,
		Msg:   "查询成功",
		Total: count,
		Parts: func(orders []*models.TbAlgoOrder) (ret []*proto.AlgoOrder) {
			for _, o := range orders {
				id++
				a := &proto.AlgoOrder{
					Id:            int64(id),
					Date:          int32(o.Date),
					BatchNo:       o.BatchNo,
					BatchName:     o.BatchName,
					BasketId:      int32(o.BasketId),
					AlgoId:        int32(o.AlgoId),
					AlgorithmId:   int32(o.AlgorithmId),
					AlgoName:      mAlgo[o.AlgorithmId],
					AlgorithmType: int32(o.AlgorithmType),
					UsecId:        int32(o.UsecId),
					UserId:        o.UserId,
					SecId:         o.SecId,
					AlgoOrderQty:  o.AlgoOrderQty,
					TransTime:     strconv.FormatInt(o.TransTime, 10),
					StartTime:     strconv.FormatInt(o.StartTime, 10),
					EndTime:       strconv.FormatInt(o.EndTime, 10),
					UnixTime:      o.UnixTime,
					FixFlag:       int32(o.Source),
					CreateTime:    o.CreateTime.Format("2006-01-02 15:04:05"),
				}
				ret = append(ret, a)
			}
			return ret
		}(orders),
	}, nil
}
