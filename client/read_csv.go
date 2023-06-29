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
	CsvFile            = "order.csv"
	ChildTopicName     = "child_order_80"
	AlgoOrderTopicName = "algo_order_80"
	Producer           sarama.SyncProducer

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
	if len(os.Args) < 3 {
		fmt.Println("usage: ./read_csv [type] [file_name]")
		fmt.Println("[type]: 1-母单模板， 2-子单模板")
		fmt.Println("[file_name]: 文件名")
		return
	}
	if os.Args[1] == "1" {
		data := ReadAlgoOrderCsv(os.Args[2])
		if len(data) > 0 {
			SendAlgoOrderData2Kafka(data)
		}
	} else if os.Args[1] == "2" {
		data := ReadChildCsv(os.Args[2])
		if len(data) > 0 {
			SendChildData2Kafka(data)
		}
	}

	defer Producer.Close()
}

func ReadAlgoOrderCsv(file string) []pb.AlgoOrderPerf {
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
	var data []pb.AlgoOrderPerf
	for k, line := range ReadAll {
		if k == 0 { // 跳过页眉
			continue
		}
		if len(line) < 7 {
			fmt.Println("algo order template file error")
			return nil
		}
		log.Println("get line:", line)
		st, _ := time.ParseInLocation(timeLayoutStr, line[6], time.Local)
		t := st.UnixMicro()
		algoOrder := pb.AlgoOrderPerf{
			Id:            cast.ToUint32(line[0]),
			AlgorithmType: cast.ToUint32(line[1]),
			AlgorithmId:   cast.ToUint32(line[2]),
			USecurityId:   cast.ToUint32(line[3]),
			SecurityId:    line[4],
			AlgoOrderQty:  cast.ToUint64(line[5]),
			TransactTime:  cast.ToUint64(t),
		}
		log.Printf("get algo order:%+v", algoOrder)
		data = append(data, algoOrder)
	}
	return data
}

func ReadChildCsv(file string) []pb.ChildOrderPerf {
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
			fmt.Println("child order template file error")
			return nil
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

func SendChildData2Kafka(s []pb.ChildOrderPerf) {
	for _, v := range s {
		msg := v
		b, _ := proto.Marshal(&msg)

		smsg := &sarama.ProducerMessage{
			Topic: ChildTopicName,
			Key:   sarama.StringEncoder("hawrk"),
			Value: sarama.StringEncoder(b),
		}
		// 发送消息
		partition, offset, err := Producer.SendMessage(smsg)
		log.Println("partition:", partition, ", offset:", offset, ", err:", err)
	}
}

func SendAlgoOrderData2Kafka(s []pb.AlgoOrderPerf) {
	for _, v := range s {
		msg := v
		b, _ := proto.Marshal(&msg)

		smsg := &sarama.ProducerMessage{
			Topic: AlgoOrderTopicName,
			Key:   sarama.StringEncoder("hawrk"),
			Value: sarama.StringEncoder(b),
		}
		// 发送消息
		partition, offset, err := Producer.SendMessage(smsg)
		log.Println("partition:", partition, ", offset:", offset, ", err:", err)
	}
}
