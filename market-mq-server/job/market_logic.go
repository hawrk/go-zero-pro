// Package job
/*
 Author: hawrkchen
 Date: 2022/5/24 19:02
 Desc:
*/
package job

import (
	"algo_assess/global"
	"algo_assess/market-mq-server/internal/config"
	"algo_assess/market-mq-server/internal/svc"
	"algo_assess/pkg/tools"
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type MarketJober interface {
	SyncRedisToDb() // 同步redis 的数据到 DB， 并删除 Redis的数据
	SayHello()
}

type MarketJob struct {
	s *svc.ServiceContext
	logx.Logger
}

func NewMarketJob(c config.Config) MarketJober {
	return &MarketJob{
		s:      svc.NewServiceContext(c),
		Logger: logx.WithContext(context.Background()),
	}
}

// SayHello debug function
func (m *MarketJob) SayHello() {
	m.Logger.Info("say hello")
}

// SyncRedisToDb 每日行情数据落地，缓存数据直接清除
func (m *MarketJob) SyncRedisToDb() {
	m.Logger.Info("start backup redis data.....")
	t := time.Now()
	m.dealSzJob()
	m.dealShJob()
	m.dealCacheJob()
	//m.DealKafkaJob()
	m.Logger.Info("finished backup redis data....,latency:", time.Since(t))
}

func (m *MarketJob) dealSzJob() {
	sli, err := m.s.RedisClient.Keys("sz:*")
	if err != nil {
		m.Logger.Error("get sz keys pattern fail:", err)
		return
	}
	m.Logger.Info("get sz key len:", len(sli))
	for _, hkey := range sli {
		// 根据 key 读取redis 数据
		m.Logger.Info("proc key:", hkey)
		data, err := m.s.RedisClient.Hgetall(hkey) // sz:000001 一个key下一般有242条记录，可以批量插入
		if err != nil {
			m.Logger.Error("get redis key fail:", hkey)
			continue
		}
		batchData := make([]*global.QuoteLevel2Data, 0, len(data))
		for _, v := range data {
			var out global.QuoteLevel2Data
			if err := json.Unmarshal(tools.String2Bytes(v), &out); err != nil {
				m.Logger.Error("Unmarshal error:", err)
				continue
			}
			batchData = append(batchData, &out)
		}
		//m.Logger.Info("proc key:", hkey, "filed len:", len(batchData))
		// 批量插入
		if err := m.s.MarketLevelRepo.CreateMarketLevelBatch(context.Background(), batchData); err != nil {
			m.Logger.Error("batch create sz market data error:", err)
			continue
		}
		// 再删 redis
		_, err = m.s.RedisClient.Del(hkey)
		if err != nil {
			m.Logger.Error(" del key fail:", err)
		}
	}
}

func (m *MarketJob) dealShJob() {
	sh, err := m.s.RedisClient.Keys("sh:*")
	if err != nil {
		m.Logger.Error("get sh keys pattern fail:", err)
		return
	}
	m.Logger.Info("get sh key len:", len(sh))
	for _, hkey := range sh {
		// 根据 key 读取redis 数据
		//m.Logger.Info("proc key:", hkey)
		data, err := m.s.RedisClient.Hgetall(hkey)
		if err != nil {
			m.Logger.Error("get redis key fail:", hkey)
			continue
		}
		batchData := make([]*global.QuoteLevel2Data, 0, len(data))
		for _, v := range data {
			var out global.QuoteLevel2Data
			if err := json.Unmarshal([]byte(v), &out); err != nil {
				m.Logger.Error("unmarshal error:", err)
				continue
			}
			batchData = append(batchData, &out)
		}
		// 批量插入
		if err := m.s.ShMarketLevelRepo.CreateMarketLevelBatch(context.Background(), batchData); err != nil {
			m.Logger.Error("batch create sh market data error:", err)
			continue
		}
		// 再删 redis
		_, err = m.s.RedisClient.Del(hkey)
		if err != nil {
			m.Logger.Error(" del key fail:", err)
		}
	}
}

// dealCacheJob 干掉 level2 开头的所有key
func (m *MarketJob) dealCacheJob() {
	//TODO: 模糊匹配会比较慢，考虑用scan进行优化
	/*
		level2, err := m.s.RedisClient.Keys("level2:*")
		if err != nil {
			m.Logger.Error("get level2 keys pattern fail:", err)
			return
		}
		m.Logger.Info("get level2 key len:", len(level2))
		// 直接删除
		_, err = m.s.RedisClient.Del(level2...)
		if err != nil {
			m.Logger.Error("del level2 key fail:", err)
		}
	*/
	var cursor uint64 // 起始游标
	count := 0        // 计算已删除的条数
	for {
		keys, cur, err := m.s.RedisClient.Scan(cursor, "level2:*", 1000)
		if err != nil {
			m.Logger.Error("redis scan err:", err, ",current cursor:", cursor)
			continue
		}
		count += len(keys)
		//m.Logger.Info("get level2 key len:", len(keys))
		// 直接删除
		_, err = m.s.RedisClient.Del(keys...)
		if err != nil {
			m.Logger.Error("del level2 key fail:", err)
		}
		cursor = cur
		if cursor == 0 {
			break
		}
	}
	m.Logger.Info("del level2 count:", count)

}

func (m *MarketJob) DealKafkaJob() {
	config := sarama.NewConfig()
	// WaitForAll waits for all in-sync replicas to commit before responding.
	config.Producer.RequiredAcks = sarama.WaitForAll // 等待ack确认，防止发送消息丢失
	//  config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Partitioner = sarama.NewHashPartitioner
	config.Producer.Return.Successes = true
	topic := m.s.Config.AlgoPlatFormSHMarketConf.Topic
	brokers := m.s.Config.AlgoPlatFormSHMarketConf.Brokers
	m.Logger.Info("clear topic:", topic)
	t, err := sarama.NewClusterAdmin(brokers, config)
	if err != nil {
		m.Logger.Info("error clusterAdmin：", err)
		return
	}
	if err := t.DeleteTopic(topic); err != nil {
		m.Logger.Info("err delete topic:", err)
		return
	}
}
