package logic

import (
	"algo_assess/global"
	"algo_assess/mornano-rpc-server/internal/svc"
	"algo_assess/mornano-rpc-server/proto"
	"algo_assess/pkg/tools"
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAlgoChooseListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAlgoChooseListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAlgoChooseListLogic {
	return &GetAlgoChooseListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetAlgoChooseList 算法选择框数据
func (l *GetAlgoChooseListLogic) GetAlgoChooseList(in *proto.AlgoChooseReq) (*proto.AlgoChooseRsp, error) {
	// todo: add your logic here and delete this line
	l.Logger.Info("in GetAlgoChooseList, req:", in)
	// 1. 根据用户ID查用户表找到algo_group id
	groupId, err := l.svcCtx.UserInfoRepo.GetGroupIdByUserId(l.ctx, in.GetUserId())
	if err != nil {
		l.Logger.Error("GetGroupIdByUserId error:", err)
		return &proto.AlgoChooseRsp{}, err
	}
	if groupId == 0 { // 如果group id为0,即无该数据
		l.Logger.Error("current no group id")
		return &proto.AlgoChooseRsp{}, nil
	}
	l.Logger.Info("get groupId:", groupId)
	// 2. 根据group id在algo_group表中找到属性值并转成二进制
	property, err := l.svcCtx.AlgoGroupRepo.GetAlgoPropertyById(l.ctx, groupId)
	if err != nil {
		l.Logger.Error("GetAlgoPropertyById error:", err)
		return &proto.AlgoChooseRsp{}, err
	}
	p := tools.RMu0000(property)
	l.Logger.Info("get property:", p)
	b := tools.Hex2Binary(p)
	l.Logger.Info("get binary:", b)
	ids := l.GetAuthAlgoList(b)
	l.Logger.Info("get algo id list:", ids)
	// 3. 在算法表中找到算法名称和厂商名称
	if in.GetSelectType() == 1 { // 取厂商列表
		pn, err := l.svcCtx.AlgoInfoRepo.GetAlgoProviders(l.ctx, ids)
		if err != nil {
			l.Logger.Error("GetAlgoProviders error:", err)
			return &proto.AlgoChooseRsp{}, err
		}
		l.Logger.Info("get provider:", pn)
		return &proto.AlgoChooseRsp{
			Code:         200,
			Msg:          "success",
			Provider:     pn,
			AlgoTypeName: nil,
			AlgoName:     nil,
		}, nil
	} else if in.GetSelectType() == 2 { // 取算法类型名称
		if in.GetProvider() == "" {
			return &proto.AlgoChooseRsp{
				Code: 204,
				Msg:  errors.New("field provider not set").Error(),
			}, errors.New("field provider not set")
		}
		names, err := l.svcCtx.AlgoInfoRepo.GetAlgoTypeNames(l.ctx, ids, in.GetProvider())
		if err != nil {
			l.Logger.Error("GetAlgoTypeNames error:", err)
			return &proto.AlgoChooseRsp{}, err
		}
		l.Logger.Info("get algo_type_name :", names)
		return &proto.AlgoChooseRsp{
			Code:         200,
			Msg:          "success",
			Provider:     nil,
			AlgoTypeName: names,
			AlgoName:     nil,
		}, nil
	} else if in.GetSelectType() == 3 { // 取算法名称
		var algoType int
		if in.GetAlgoTypeName() == global.AlgoTypeNameT0 {
			algoType = 1
		} else if in.GetAlgoTypeName() == global.AlgoTypeNameSplit {
			algoType = 2
		}
		n, err := l.svcCtx.AlgoInfoRepo.GetAlgoNameByIds(l.ctx, ids, in.GetProvider(), algoType)
		if err != nil {
			l.Logger.Error("GetAlgoNameByIds", err)
			return &proto.AlgoChooseRsp{}, err
		}
		l.Logger.Info("get algo_name:", n)
		return &proto.AlgoChooseRsp{
			Code:         200,
			Msg:          "success",
			Provider:     nil,
			AlgoTypeName: nil,
			AlgoName:     n,
		}, nil
	} else if in.GetSelectType() == 9 { // 根据算法名称，查回厂商名称和算法类型名称
		if in.GetAlgoName() == "" {
			return &proto.AlgoChooseRsp{
				Code: 204,
				Msg:  errors.New("field algo_name not set").Error(),
			}, errors.New("field algo_name not set")
		}
		p, a, err := l.svcCtx.AlgoInfoRepo.GetAlgoInfoByName(l.ctx, in.GetAlgoName())
		if err != nil {
			l.Logger.Error("GetAlgoInfoByName", err)
			return &proto.AlgoChooseRsp{}, err
		}
		return &proto.AlgoChooseRsp{
			Code:         200,
			Msg:          "success",
			Provider:     []string{p},
			AlgoTypeName: []string{a},
			AlgoName:     []string{in.GetAlgoName()},
		}, nil
	}

	return &proto.AlgoChooseRsp{
		Code:         400,
		Msg:          "unsupported select_type",
		Provider:     nil,
		AlgoTypeName: nil,
		AlgoName:     nil,
	}, nil
}

// GetAuthAlgoList 根据解析出来的二进制，从右往左，对应的位数为1时即有该算法ID的权限
func (l *GetAlgoChooseListLogic) GetAuthAlgoList(b string) []int {
	n := len(b)
	var list []int
	for i := n - 1; i >= 0; i-- {
		if b[i] == '1' {
			index := n - i
			list = append(list, index)
		}
	}
	return list
}
