package logic

import (
	"algo_assess/models"
	"context"
	"fmt"

	"algo_assess/assess-rpc-server/internal/svc"
	"algo_assess/assess-rpc-server/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type SelectOptimizeBaseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSelectOptimizeBaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SelectOptimizeBaseLogic {
	return &SelectOptimizeBaseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询一键优选基础数据
func (l *SelectOptimizeBaseLogic) SelectOptimizeBase(in *proto.SelectOptimizeBaseReq) (*proto.SelectOptimizeBaseRsp, error) {
	count, result, err := l.svcCtx.OptimizeBaseRepo.SelectAlgoOptimizeBase(in)
	//如果无数据就从最后一页开始查询
	if count > 0 && len(result) == 0 {
		//分页获取的数据不对，就从最后一页查
		fmt.Printf("count:%v,len(infos):%v\n", count, len(result))
		i := count / int64(in.GetLimit())
		j := count % int64(in.GetLimit())
		if j > 0 {
			i = i + 1
		}
		fmt.Printf("i:%v,j:%v\n", i, j)
		in.Page = int32(i)
		count, result, err = l.svcCtx.OptimizeBaseRepo.SelectAlgoOptimizeBase(in)
	}
	if err == nil {
		return &proto.SelectOptimizeBaseRsp{
			Code:  0,
			Msg:   "查询成功",
			Total: count,
			List:  ToProtoOptimizeBase(result),
		}, nil
	} else {
		return &proto.SelectOptimizeBaseRsp{
			Code:  20000,
			Msg:   "查询失败",
			Total: 0,
			List:  nil,
		}, err
	}
}

func ToProtoOptimizeBase(result []*models.TbAlgoOptimizeBase) (ret []*proto.OptimizeBase) {
	for _, o := range result {
		optimize := &proto.OptimizeBase{
			Id:           o.Id,
			ProviderId:   int32(o.ProviderId),
			ProviderName: o.ProviderName,
			SecId:        o.SecId,
			SecName:      o.SecName,
			AlgoId:       int32(o.AlgoId),
			AlgoType:     int32(o.AlgoType),
			AlgoName:     o.AlgoName,
			OpenRate:     o.OpenRate,
			IncomeRate:   o.IncomeRate,
			BasisPoint:   o.BasisPoint,
			CreateTime:   o.CreateTime.Format("2006-01-02 15:04:05"),
			UpdateTime:   o.UpdateTime.Format("2006-01-02 15:04:05"),
		}
		ret = append(ret, optimize)
	}
	return ret
}
