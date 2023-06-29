package logic

import (
	"algo_assess/assess-rpc-server/internal/svc"
	"algo_assess/assess-rpc-server/proto"
	"algo_assess/global"
	"algo_assess/pkg/tools"
	"context"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type PushAlgoOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPushAlgoOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PushAlgoOrderLogic {
	return &PushAlgoOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// PushAlgoOrder 母单绩效数据修复
func (l *PushAlgoOrderLogic) PushAlgoOrder(in *proto.AlgoOrderPerfs) (*proto.PushDataRsp, error) {
	parts := in.GetParts()
	//var date string // 格式 "202305056"
	var algoIds []int
	for _, v := range parts {
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

	for _, part := range parts {
		t := time.UnixMicro(int64(part.GetTransactTime())).Format(global.TimeFormatMinInt)
		transactAt := tools.TimeMoveForward(t)
		//date = cast.ToString(transactAt)[:8]

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
		if _, exist := m[int64(part.GetId())]; !exist { // 母单号不存在，先补一份到DB
			if err := l.svcCtx.AlgoOrderRepo.CreateAlgoOrder(l.ctx, algoOrder); err != nil {
				l.Logger.Error("CreateAlgoOrder error:", err)
				continue
			}
		} else { // 母单号已存在， 根据需求，该母单数据可能会被修改，所以需要更新母单数据
			l.Logger.Info("algoID exist, update")
			if err := l.svcCtx.AlgoOrderRepo.UpdateAlgoOrder(l.ctx, algoOrder); err != nil {
				l.Logger.Error("UpdateAlgoOrder error:")
				continue
			}
		}
	}
	// 这部分逻辑放到子单里
	/*
		// 更新完母单表之后，需要从数据库里所有了关联的算法ID
		var algorithmIds []int32
		mId := make(map[int32]struct{}) // 辅助algorithmIds 去重
		for _, v := range parts {
			if _, exist := mId[int32(v.AlgorithmId)]; !exist {
				algorithmIds = append(algorithmIds, int32(v.AlgorithmId))
				mId[int32(v.AlgorithmId)] = struct{}{}
			}
		}
		//把所有该算法ID的母单号找出来
		out, err := l.svcCtx.AlgoOrderRepo.GetAlgoByAlgorithms(l.ctx, algorithmIds, tools.TimeDay2string(date))
		if err != nil {
			l.Logger.Error("get algoOrder List error:", err)
			return nil, nil
		}

		for _, v := range out {
			o := &proto.AlgoOrderPerf{
				Id:       uint32(v.AlgoId),
				BasketId: uint32(v.BasketId),
				//AlgorithmType: 0,
				AlgorithmId:  uint32(v.AlgorithmId),
				USecurityId:  uint32(v.UsecId),
				SecurityId:   v.SecId,
				AlgoOrderQty: uint64(v.AlgoOrderQty) * 100,
				TransactTime: uint64(v.UnixTime),
				StartTime:    uint64(v.StartTime),
				EndTime:      uint64(v.EndTime),
				BusUserId:    v.UserId,
				BatchNo:      v.BatchNo,
				BatchName:    v.BatchName,
				SourceFrom:   1,
			}
			// 发送Kafka
			bytes, ptErr := pt.Marshal(o)
			if ptErr != nil {
				l.Logger.Error("Marshal msg error:", ptErr)
				continue
			}
			//l.Logger.Infof("push ProtoMessage:%+v", o)
			msg := &sarama.ProducerMessage{
				Topic:     l.svcCtx.Config.Kafka.AlgoTopic, // 统一放入正常队列
				Value:     sarama.ByteEncoder(bytes),
				Partition: 0, // 指定0分区，需要在初始化时指定sarama.NewManualPartitioner
			}
			_, _, kafkaErr := l.svcCtx.SyncProducer.SendMessage(msg)
			if kafkaErr != nil {
				l.Logger.Error("send message error:", msg.Value, kafkaErr)
				continue
			}
		}
	*/

	return &proto.PushDataRsp{
		Code: 200,
		Msg:  "处理成功",
	}, nil
}
