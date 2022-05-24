// Package main
/*
 Author: hawrkchen
 Date: 2022/4/29 10:21
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
	//config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Partitioner = sarama.NewHashPartitioner
	config.Producer.Return.Successes = true

	// 新建一个同步生产者
	producer, err := sarama.NewSyncProducer([]string{"192.168.1.85:9092"}, config)
	if err != nil {
		fmt.Println("producer close, err:", err)
		return
	}
	defer producer.Close()

	msg := &pb.AlgoOrderPerf{
		Id:            10086,
		AlgorithmType: 5,
		AlgorithmId:   1,
		USecurityId:   433,
		SecurityId:    "000038",
		AlgoOrderQty:  5000000,          // 5W 笔
		TransactTime:  1652665200000000, //微秒
	}

	byte, _ := proto.Marshal(msg)
	// 定义一个生产消息，包括Topic、消息内容、
	smsg := &sarama.ProducerMessage{
		Topic:     "sh_quote",
		Key:       sarama.StringEncoder("hawrk1"), // hash 放到0 partition
		Value:     sarama.StringEncoder(byte),
		Partition: 0, // 生产的消息不能指定Partition, 消息投送成功才会产生
	}
	// 发送消息
	partition, offset, err := producer.SendMessage(smsg)
	fmt.Println("patition:", partition, "offset:", offset)

	////////
	msg1 := &pb.AlgoOrderPerf{
		Id:            10088,
		AlgorithmType: 5,
		AlgorithmId:   1,
		USecurityId:   433,
		SecurityId:    "000038",
		AlgoOrderQty:  1000000,          // 5W 笔
		TransactTime:  1652665200000000, //微秒
	}

	byte1, _ := proto.Marshal(msg1)
	// 定义一个生产消息，包括Topic、消息内容、
	smsg1 := &sarama.ProducerMessage{
		Topic:     "sh_quote",
		Key:       sarama.StringEncoder("hawrk1"),
		Value:     sarama.StringEncoder(byte1),
		Partition: 0,
	}
	// 发送消息
	partition1, offset1, err := producer.SendMessage(smsg1)
	fmt.Println("patition:", partition1, "offset:", offset1)

	/*
		//////////
		msg2 := &pb.AlgoOrderPerf{
			Id:            125,
			AlgorithmType: 22,
			AlgorithmId:   44,
			USecurityId:   433,
			SecurityId:    "0000047",
			AlgoOrderQty:  1000,
			TransactTime:  1650933000,
		}

		byte2, _ := proto.Marshal(msg2)
		// 定义一个生产消息，包括Topic、消息内容、
		smsg2 := &sarama.ProducerMessage{
			Topic:     "algo_order",
			Key:       sarama.StringEncoder("hawrk2"),
			Value:     sarama.StringEncoder(byte2),
			Partition: 2,
		}
		// 发送消息
		partition2, offset2, err := producer.SendMessage(smsg2)
		fmt.Println("patition:", partition2, "offset:", offset2)
	*/

}
