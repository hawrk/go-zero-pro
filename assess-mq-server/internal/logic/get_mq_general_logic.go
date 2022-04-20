package logic

import (
	"algo_assess/global"
	"context"
	"github.com/jinzhu/copier"

	"algo_assess/assess-mq-server/internal/svc"
	"algo_assess/assess-mq-server/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMqGeneralLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMqGeneralLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMqGeneralLogic {
	return &GetMqGeneralLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取绩效概况
func (l *GetMqGeneralLogic) GetMqGeneral(in *proto.GeneralReq) (*proto.GeneralRsp, error) {
	l.Logger.Info("into mq GetGeneral:req:", in)
	// 读取本地缓存数据
	global.GlobalAssess.RWMutex.RLock()
	infos := make([]*proto.AssessInfo, 0, len(global.GlobalAssess.CalAlgo))
	for _, val := range global.GlobalAssess.CalAlgo {
		var info proto.AssessInfo
		copier.Copy(&info, val)
		infos = append(infos, &info)
	}
	global.GlobalAssess.RWMutex.RUnlock()
	rsp := &proto.GeneralRsp{
		Code: 0,
		Msg:  "success",
		Info: infos,
	}
	return rsp, nil
}
