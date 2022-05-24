// Package main
/*
 Author: hawrkchen
 Date: 2022/5/5 16:06
 Desc:
*/
package main

import (
	pb "algo_assess/assess-mq-server/proto/order"
	"encoding/csv"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"github.com/spf13/cast"
	"log"
	"os"
	"time"
)

var (
	CsvFile   = "order.csv"
	TopicName = "hawrk_demo"
	Producer  sarama.SyncProducer

	timeLayoutStr = "2006-01-02 15:04"
)

func init() {
	// 新建一个arama配置实例
	config := sarama.NewConfig()
	// WaitForAll waits for all in-sync replicas to commit before responding.
	config.Producer.RequiredAcks = sarama.WaitForAll // 等待ack确认，防止发送消息丢失
	// config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Partitioner = sarama.NewHashPartitioner
	config.Producer.Return.Successes = true

	// 新建一个同步生产者
	producer, err := sarama.NewSyncProducer([]string{"192.168.1.85:9092"}, config)
	if err != nil {
		fmt.Println("producer close, err:", err)
		return
	}
	Producer = producer
}

func main() {
	data := ReadCsv(CsvFile)
	Send2Kafka(data)

	defer Producer.Close()
}

func ReadCsv(file string) []pb.ChildOrderPerf {
	f, err := os.Open(file)
	if err != nil {
		log.Println("open file fail:%v", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	//
	ReadAll, err := reader.ReadAll()
	if err != nil {
		log.Println("read file fail %v:", err)
	}
	var data []pb.ChildOrderPerf
	for k, line := range ReadAll {
		if k == 0 { // 跳过页眉
			continue
		}
		if len(line) < 15 {
			continue
		}
		log.Println("get line:", line)
		st, _ := time.ParseInLocation(timeLayoutStr, line[14], time.Local)
		t := st.UnixMicro()
		order := pb.ChildOrderPerf{
			Id:             cast.ToUint32(line[0]),
			AlgoOrderId:    cast.ToUint32(line[1]),
			AlgorithmType:  cast.ToUint32(line[2]),
			AlgorithmId:    cast.ToUint32(line[3]),
			USecurityId:    cast.ToUint32(line[4]),
			SecurityId:     line[5],
			OrderQty:       cast.ToUint64(line[6]),
			Price:          cast.ToUint64(line[7]),
			OrderType:      cast.ToUint32(line[8]),
			CumQty:         cast.ToUint64(line[9]),
			LastPx:         cast.ToUint64(line[10]),
			LastQty:        cast.ToUint64(line[11]),
			ArrivedPrice:   cast.ToUint64(line[12]),
			ChildOrdStatus: cast.ToUint32(line[13]),
			TransactTime:   cast.ToUint64(t),
		}
		log.Printf("get order:%+v", order)
		data = append(data, order)
	}
	return data
}

func Send2Kafka(s []pb.ChildOrderPerf) {
	for _, v := range s {
		msg := v
		b, _ := proto.Marshal(&msg)

		smsg := &sarama.ProducerMessage{
			Topic: TopicName,
			Key:   sarama.StringEncoder("hawrk"),
			Value: sarama.StringEncoder(b),
		}
		// 发送消息
		partition, offset, err := Producer.SendMessage(smsg)
		log.Println("partition:", partition, ", offset:", offset, ", err:", err)
	}
}
