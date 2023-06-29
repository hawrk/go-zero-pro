package logic

import (
	"algo_assess/assess-rpc-server/internal/svc"
	"algo_assess/assess-rpc-server/proto"
	"algo_assess/global"
	"algo_assess/pkg/tools"
	"context"
	"github.com/Shopify/sarama"
	"github.com/spf13/cast"
	pt "google.golang.org/protobuf/proto"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type ImportAlgoOrdrLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewImportAlgoOrdrLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImportAlgoOrdrLogic {
	return &ImportAlgoOrdrLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ImportAlgoOrdr 订单导入:母单信息导入
func (l *ImportAlgoOrdrLogic) ImportAlgoOrdr(in *proto.AlgoOrderPerfs) (*proto.PushDataRsp, error) {
	l.Logger.Info("in ImportAlgoOrdr, get req len:", len(in.GetParts()))
	var algoIds []int
	for _, v := range in.GetParts() {
		algoIds = append(algoIds, int(v.Id))
	}
	// 1. 根据传入母单查DB数据
	ids, err := l.svcCtx.AlgoOrderRepo.BatchQueryAlgoOrder(l.ctx, algoIds)
	if err != nil {
		l.Logger.Error("BatchQueryAlgoOrder error:", err)
		return &proto.PushDataRsp{
			Code: 10000,
			Msg:  "母单信息查询失败",
		}, nil
	}
	// DB的母单列列转换成map格式
	m := make(map[int64]struct{}, len(ids))
	for _, v := range ids {
		m[v] = struct{}{}
	}
	var timeList []int64
	// 落表
	for _, part := range in.GetParts() {
		t := time.UnixMicro(int64(part.GetTransactTime())).Format(global.TimeFormatMinInt)
		transactAt := tools.TimeMoveForward(t)
		timeList = append(timeList, transactAt) // 跨天标识用
		if _, exist := m[int64(part.GetId())]; !exist {
			algoOrder := &global.MAlgoOrder{
				BatchNo:         part.BatchNo, // 与总线保持一致
				AlgoId:          int(part.Id),
				UserId:          part.BusUserId,
				BasketId:        int(part.BasketId),
				AlgorithmId:     int(part.AlgorithmId),
				UsecId:          int(part.USecurityId),
				SecId:           part.SecurityId,
				AlgoOrderQty:    int64(part.AlgoOrderQty) / 100,
				UnixTime:        cast.ToString(part.GetTransactTime()), // 精确到分钟的时间戳
				UnixTimeMillSec: int64(part.GetTransactTime()),         // 原始时间戳
				TransTime:       transactAt,
				StartTime:       int64(part.StartTime),
				EndTime:         int64(part.EndTime),
				SourceFrom:      int(part.SourceFrom), // 修复标识
			}
			if err := l.svcCtx.AlgoOrderRepo.CreateAlgoOrder(l.ctx, algoOrder); err != nil {
				l.Logger.Error("CreateAlgoOrder error:", err)
				continue
			}
		}

		// 发送Kafka
		bytes, ptErr := pt.Marshal(part)
		if ptErr != nil {
			l.Logger.Error("Marshal msg error:", ptErr)
			continue
			//return &proto.PushDataRsp{
			//	Code: 10000,
			//	Msg:  "proto格式错误",
			//}, nil
		}
		//l.Logger.Infof("push ProtoMessage:%+v", part)
		msg := &sarama.ProducerMessage{
			Topic:     l.svcCtx.Config.Kafka.AlgoTopic, // 统一放入正常队列
			Value:     sarama.ByteEncoder(bytes),
			Partition: 0, // 指定0分区，需要在初始化时指定sarama.NewManualPartitioner
		}
		_, _, kafkaErr := l.svcCtx.SyncProducer.SendMessage(msg)
		if kafkaErr != nil {
			l.Logger.Errorf("send message(%s) err=%s \n", msg.Value, kafkaErr)
			continue
			//return &proto.PushDataRsp{
			//	Code: 10000,
			//	Msg:  "kafka推送错误",
			//}, nil
		}
	}
	start, end := tools.GetDoubleTimePoint(timeList)
	return &proto.PushDataRsp{
		Code:      200,
		Msg:       "success",
		StartTime: start,
		EndTime:   end,
	}, nil
}
