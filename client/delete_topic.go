// Package main
/*
 Author: hawrkchen
 Date: 2022/4/29 11:04
 Desc:   删除指定Topic 的所有数据， 慎用！！！！！！
*/
package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"os"
)

var topic = "algo_order"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: ./delete_topic  [topic_name]")
		return
	}
	fmt.Println("begin to delete topic :", os.Args[1])
	// 新建一个arama配置实例
	config := sarama.NewConfig()
	// WaitForAll waits for all in-sync replicas to commit before responding.
	config.Producer.RequiredAcks = sarama.WaitForAll // 等待ack确认，防止发送消息丢失
	//config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Partitioner = sarama.NewHashPartitioner
	config.Producer.Return.Successes = true

	t, err := sarama.NewClusterAdmin([]string{"192.168.1.85:9092"}, config)
	if err != nil {
		log.Println("error clusterAdmin：", err)
		return
	}
	if err := t.DeleteTopic(os.Args[1]); err != nil {
		log.Println("errr delete topic:", err)
		return
	}
	fmt.Println(" done!")
}
