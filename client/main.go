// Package client
/*
 Author: hawrkchen
 Date: 2022/4/1 9:59
 Desc:
*/
package main

import (
	pb "algo_assess/assess-mq-server/proto/order"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
)

func main() {
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
	defer producer.Close()

	//1 部分成交
	msg1 := &pb.ChildOrderPerf{
		Id:             2022000001,
		AlgorithmType:  10003,
		AlgorithmId:    5,
		USecurityId:    8868,
		SecurityId:     "00700",
		OrderQty:       1000,
		Price:          300,
		OrderType:      1,
		CumQty:         700,
		LastPx:         380,
		LastQty:        100,
		ChildOrdStatus: 6,
		TransactTime:   1648516164,
	}
	byte1, _ := proto.Marshal(msg1)
	// 定义一个生产消息，包括Topic、消息内容、
	smsg1 := &sarama.ProducerMessage{
		Topic: "hawrk_demo",
		Key:   sarama.StringEncoder("hawrk"),
		Value: sarama.StringEncoder(byte1),
	}
	// 发送消息
	partition, offset, err := producer.SendMessage(smsg1)
	fmt.Println("patition:", partition, "offset:", offset)

	//2 部分成交
	msg2 := &pb.ChildOrderPerf{
		Id:             2022000001,
		AlgorithmType:  10003,
		AlgorithmId:    5,
		USecurityId:    8868,
		SecurityId:     "00700",
		OrderQty:       1000,
		Price:          300,
		OrderType:      1,
		CumQty:         600,
		LastPx:         338,
		LastQty:        400,
		ChildOrdStatus: 6,
		TransactTime:   1648516164,
	}
	byte2, _ := proto.Marshal(msg2)
	// 定义一个生产消息，包括Topic、消息内容、
	smsg2 := &sarama.ProducerMessage{
		Topic: "hawrk_demo",
		Key:   sarama.StringEncoder("hawrk"),
		Value: sarama.StringEncoder(byte2),
	}
	// 发送消息
	partition2, offset2, err := producer.SendMessage(smsg2)
	fmt.Println("patition:", partition2, "offset:", offset2)

	//3 撤销
	msg3 := &pb.ChildOrderPerf{
		Id:             2022000001,
		AlgorithmType:  10003,
		AlgorithmId:    5,
		USecurityId:    8868,
		SecurityId:     "00700",
		OrderQty:       1000,
		Price:          300,
		OrderType:      1,
		CumQty:         700,
		LastPx:         0,
		LastQty:        0,
		ChildOrdStatus: 8,
		TransactTime:   1648516164,
	}
	byte3, _ := proto.Marshal(msg3)
	// 定义一个生产消息，包括Topic、消息内容、
	smsg3 := &sarama.ProducerMessage{
		Topic: "hawrk_demo",
		Key:   sarama.StringEncoder("hawrk"),
		Value: sarama.StringEncoder(byte3),
	}
	// 发送消息
	partition3, offset3, err := producer.SendMessage(smsg3)
	fmt.Println("patition3:", partition3, "offset3:", offset3)

	//4 拒单
	msg4 := &pb.ChildOrderPerf{
		Id:             2022000002,
		AlgorithmType:  10003,
		AlgorithmId:    5,
		USecurityId:    8868,
		SecurityId:     "00700",
		OrderQty:       660,
		Price:          66,
		OrderType:      1,
		CumQty:         0,
		LastPx:         0,
		LastQty:        0,
		ChildOrdStatus: 5,
		TransactTime:   1648516164,
	}
	byte4, _ := proto.Marshal(msg4)
	// 定义一个生产消息，包括Topic、消息内容、
	smsg4 := &sarama.ProducerMessage{
		Topic: "hawrk_demo",
		Key:   sarama.StringEncoder("hawrk"),
		Value: sarama.StringEncoder(byte4),
	}
	// 发送消息
	partition4, offset4, err := producer.SendMessage(smsg4)
	fmt.Println("patition4:", partition4, "offset4:", offset4)

	//5 另一单 部分成交
	msg5 := &pb.ChildOrderPerf{
		Id:             2022000003,
		AlgorithmType:  10003,
		AlgorithmId:    5,
		USecurityId:    8868,
		SecurityId:     "00700",
		OrderQty:       50,
		Price:          5,
		OrderType:      1,
		CumQty:         30,
		LastPx:         6,
		LastQty:        30,
		ChildOrdStatus: 6,
		TransactTime:   1648516164,
	}
	byte5, _ := proto.Marshal(msg5)
	// 定义一个生产消息，包括Topic、消息内容、
	smsg5 := &sarama.ProducerMessage{
		Topic: "hawrk_demo",
		Key:   sarama.StringEncoder("hawrk"),
		Value: sarama.StringEncoder(byte5),
	}
	// 发送消息
	partition5, offset5, err := producer.SendMessage(smsg5)
	fmt.Println("patition:5", partition5, "offset5:", offset5)

	//6 另一支股票部分成交
	msg6 := &pb.ChildOrderPerf{
		Id:             2022000004,
		AlgorithmType:  10005,
		AlgorithmId:    6,
		USecurityId:    88988,
		SecurityId:     "00701",
		OrderQty:       500,
		Price:          99,
		OrderType:      1,
		CumQty:         300,
		LastPx:         98,
		LastQty:        300,
		ChildOrdStatus: 6,
		TransactTime:   1648516164,
	}
	byte6, _ := proto.Marshal(msg6)
	// 定义一个生产消息，包括Topic、消息内容、
	smsg6 := &sarama.ProducerMessage{
		Topic: "hawrk_demo",
		Key:   sarama.StringEncoder("hawrk"),
		Value: sarama.StringEncoder(byte6),
	}
	// 发送消息
	partition6, offset6, err := producer.SendMessage(smsg6)
	fmt.Println("patition6:", partition6, "offset6:", offset6)

	//6 另一个时间点 一支股票部分成交
	msg7 := &pb.ChildOrderPerf{
		Id:             2022000005,
		AlgorithmType:  10008,
		AlgorithmId:    9,
		USecurityId:    88966,
		SecurityId:     "00705",
		OrderQty:       10000,
		Price:          6000,
		OrderType:      1,
		CumQty:         1000,
		LastPx:         6500,
		LastQty:        1000,
		ChildOrdStatus: 6,
		TransactTime:   1648517195,
	}
	byte7, _ := proto.Marshal(msg7)
	// 定义一个生产消息，包括Topic、消息内容、
	smsg7 := &sarama.ProducerMessage{
		Topic: "hawrk_demo",
		Key:   sarama.StringEncoder("hawrk"),
		Value: sarama.StringEncoder(byte7),
	}
	// 发送消息
	partition7, offset7, err := producer.SendMessage(smsg7)
	fmt.Println("patition7:", partition7, "offset7:", offset7)

}
