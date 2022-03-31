package logic

import (
	"algo_assess/models"
	"algo_assess/pkg/tools"
	"algo_assess/rpc/assess/internal/svc"
	"algo_assess/rpc/assess/proto"
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetDemoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDemoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDemoLogic {
	return &GetDemoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetDemoLogic) GetDemo(in *proto.DemoReq) (*proto.DemoRsp, error) {
	// todo: add your logic here and delete this line
	l.Logger.Info(" into rpc getdemo, req:", in.Hello)
	// db test model 模式
	//algoInfo, err := l.svcCtx.AlgoInfoModel.FindOne(l.ctx,1 )
	//if err != nil && err != model.ErrNotFound {
	//	l.Logger.Error("query db error")
	//}
	//l.Logger.Infof("get algoInfo :%+v", algoInfo)

	// db test 用原生的 sql
	//db := l.svcCtx.DB
	//algoInfo := make([]model.TbAlgoInfo, 0)
	//querySql := "select algo_name, provider, provider_name from " + "tb_algo_info" + " where uuser_id =1;"
	//if err := db.QueryRowsPartial(&algoInfo,querySql); err != nil {
	//	l.Logger.Error("query error:", err)
	//}
	//l.Logger.Infof("get userinfo :%+v", algoInfo)

	// db gorm
	var algoInfo []models.TbAlgoInfo
	db := l.svcCtx.DB
	db.Where(&models.TbAlgoInfo{UuserId: 1}).Find(&algoInfo)
	//l.Logger.Infof("get userinfo:%+v", algoInfo)
	for _, v := range algoInfo {
		fmt.Println("去掉空格前：", v.AlgoName, len(v.AlgoName))
		tmp := tools.RMu0000(v.AlgoName)
		fmt.Println("去掉空格后：", tmp, len(tmp))
	}

	//redis test
	//redis := l.svcCtx.RedisClient
	////redis.Set("hawrk:huayunsoft", "chen")
	//redis.Setex("hawrkhuayun", "name:chen", 30)

	resp := &proto.DemoRsp{Reply: "from rpc logic"}

	return resp, nil
}
