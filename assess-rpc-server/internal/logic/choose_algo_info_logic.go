package logic

import (
	"context"
	"errors"

	"algo_assess/assess-rpc-server/internal/svc"
	"algo_assess/assess-rpc-server/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChooseAlgoInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChooseAlgoInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChooseAlgoInfoLogic {
	return &ChooseAlgoInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ChooseAlgoInfo 算法条件筛选
func (l *ChooseAlgoInfoLogic) ChooseAlgoInfo(in *proto.ChooseAlgoReq) (*proto.ChooseAlgoRsp, error) {
	l.Logger.Info("ChooseAlgoInfo get req:", in)
	var rsp proto.ChooseAlgoRsp
	if in.GetChooseType() == 1 { // 查厂商列表
		out, err := l.svcCtx.AlgoInfoRepo.GetAlgoProvider(l.ctx)
		if err != nil {
			l.Logger.Error("get algo info provider error:", err)
			return &proto.ChooseAlgoRsp{}, nil
		}
		rsp.Provider = out
	} else if in.GetChooseType() == 2 { // 算法类型名称列表
		//if in.GetProvider() == "" {
		//	return &proto.ChooseAlgoRsp{
		//		Code: 204,
		//		Msg:  errors.New("field provider not set").Error(),
		//	}, errors.New("field provider not set")
		//}
		out, err := l.svcCtx.AlgoInfoRepo.GetAlgoTypeName(l.ctx, in.GetProvider())
		if err != nil {
			l.Logger.Error("get algo info provider error:", err)
			return &proto.ChooseAlgoRsp{}, nil
		}
		rsp.AlgoTypeName = out
	} else if in.GetChooseType() == 3 { // 算法ID列表
		//if in.GetProvider() == "" || in.GetAlgoTypeName() == "" {
		//	return &proto.ChooseAlgoRsp{
		//		Code: 204,
		//		Msg:  errors.New("field provider|algo_type_name not set").Error(),
		//	}, errors.New("field provider|algo_type_name not set")
		//}
		out, err := l.svcCtx.AlgoInfoRepo.GetAlgoName(l.ctx, in.GetProvider(), in.GetAlgoTypeName())
		if err != nil {
			l.Logger.Error("get algo info provider error:", err)
			return &proto.ChooseAlgoRsp{}, nil
		}
		rsp.AlgoName = out
	} else if in.GetChooseType() == 4 {
		if in.GetAlgoName() == "" {
			return &proto.ChooseAlgoRsp{
				Code: 204,
				Msg:  errors.New("field algo_name not set").Error(),
			}, errors.New("field algo_name not set")
		}
		out, err := l.svcCtx.AlgoInfoRepo.GetAlgoIdByAlgoName(l.ctx, in.GetAlgoName())
		if err != nil {
			l.Logger.Error("get algo info algoid error:", err)
			return &proto.ChooseAlgoRsp{}, nil
		}
		rsp.AlgoId = out
	} else if in.GetChooseType() == 5 { // 取所有算法类型名称
		out, err := l.svcCtx.AlgoInfoRepo.GetAllAlgoTypeName(l.ctx)
		if err != nil {
			l.Logger.Error("get algo info provider error:", err)
			return &proto.ChooseAlgoRsp{}, nil
		}
		rsp.AlgoTypeName = out
	} else if in.GetChooseType() == 6 { // 算法类型名称反查算法类型
		if in.GetAlgoTypeName() == "" {
			return &proto.ChooseAlgoRsp{
				Code: 204,
				Msg:  errors.New("field algo_type_name not set").Error(),
			}, errors.New("field algo_type_name not set")
		}
		out, err := l.svcCtx.AlgoInfoRepo.GetAlgoTypeIdByAlgoTypeName(l.ctx, in.GetAlgoTypeName())
		if err != nil {
			l.Logger.Error("get algo info algo type error:", err)
			return &proto.ChooseAlgoRsp{}, nil
		}
		rsp.AlgoType = out
	} else if in.GetChooseType() == 7 { // 根据算法类型查询所有算法名称
		if in.GetAlgoTypeName() == "" {
			return &proto.ChooseAlgoRsp{
				Code: 204,
				Msg:  errors.New("field algo_type_name not set").Error(),
			}, errors.New("field algo_type_name not set")
		}
		out, err := l.svcCtx.AlgoInfoRepo.GetAlgoName(l.ctx, "", in.GetAlgoTypeName())
		if err != nil {
			l.Logger.Error("get algo info algo type error:", err)
			return &proto.ChooseAlgoRsp{}, nil
		}
		rsp.AlgoName = out
	} else if in.GetChooseType() == 8 { // 根据算法类型查询当天有交易的算法列表(查汇总表）
		if in.GetDate() < 20220101 {
			return &proto.ChooseAlgoRsp{
				Code: 204,
				Msg:  errors.New("field date not set").Error(),
			}, errors.New("field date not set")
		}
		out, err := l.svcCtx.SummaryRepo.GetAlgoNameByRanking(l.ctx, in.GetDate(), "admin")
		if err != nil {
			l.Logger.Error("get algo name error:", err)
			return &proto.ChooseAlgoRsp{}, nil
		}
		rsp.AlgoName = out
	} else if in.GetChooseType() == 9 { // 选择算法名称时，返回其对应的算法厂商和算法类型名称
		out, err := l.svcCtx.AlgoInfoRepo.GetAlgoInfoByAlgoName(l.ctx, in.GetAlgoName())
		if err != nil {
			l.Logger.Error("GetAlgoInfoByAlgoName error:", err)
			return &proto.ChooseAlgoRsp{}, nil
		}
		rsp.Provider = append(rsp.Provider, out.Provider)
		rsp.AlgoTypeName = append(rsp.AlgoTypeName, out.AlgoTypeName)
		rsp.AlgoName = append(rsp.AlgoName, out.AlgoName)
		rsp.AlgoId = int32(out.AlgoId)
	} else {
		l.Logger.Error("unsupported choose algo type:", in.GetChooseType())
		return &proto.ChooseAlgoRsp{}, nil
	}

	return &rsp, nil
}
