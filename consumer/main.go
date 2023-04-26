package main

import (
	"Consumer/dbtools"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"os"
	"strconv"
	"sync"
)

// 这波我直接开摆了 个人建议不要尝试看懂这一大段代码
func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	addr := os.Getenv("MQ_ADDR")
	consumer, err := sarama.NewConsumer([]string{addr}, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	snatchPartitionList, err := consumer.Partitions("snatch")
	openPartitionList, err := consumer.Partitions("open")
	cachePartitionList, err := consumer.Partitions("cache")
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	//如果你有罪，那么请去自首，让法律制裁你，而不是来尝试看懂这段代码
	for snatchPartion := range snatchPartitionList {
		pc, err := consumer.ConsumePartition("snatch", int32(snatchPartion), sarama.OffsetNewest)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		go func(partitionConsumer sarama.PartitionConsumer) {
			for msg := range partitionConsumer.Messages() {
				type Jsin struct {
					Uid   int   `json:"uid"`
					Eid   int   `json:"eid"`
					Val   int   `json:"val"`
					Stime int64 `json:"stime"`
					N     int   `json:"n"`
				}
				var jsin Jsin
				fmt.Println("snatch value : ", msg.Value)
				json.Unmarshal(msg.Value, &jsin)
				fmt.Println("snatch json ; ", jsin)
				dbtools.SnatchWrite(jsin.Uid, jsin.Eid, jsin.Val, jsin.Stime, jsin.N)
			}
		}(pc)
	}

	for openPartion := range openPartitionList {
		pc, err := consumer.ConsumePartition("open", int32(openPartion), sarama.OffsetNewest)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		go func(partitionConsumer sarama.PartitionConsumer) {
			for msg := range partitionConsumer.Messages() {
				type Jsin struct {
					Uid int `json:"uid"`
					Eid int `json:"eid"`
					Val int `json:"val"`
				}
				var jsin Jsin
				fmt.Println("open value : ", msg.Value)
				json.Unmarshal(msg.Value, &jsin)
				dbtools.OpenWrite(jsin.Uid, jsin.Eid, jsin.Val)
			}
		}(pc)
	}

	for cachePartion := range cachePartitionList {
		pc, err := consumer.ConsumePartition("cache", int32(cachePartion), sarama.OffsetNewest)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		go func(partitionConsumer sarama.PartitionConsumer) {
			for msg := range partitionConsumer.Messages() {
				conn := dbtools.RedisPool.Get()
				uid, err := strconv.Atoi(string(msg.Value))
				if err != nil {
					fmt.Println(err.Error())
					fmt.Println("uid : ", msg.Value)
				}
				dbtools.SaveToCache(uid, conn)
				conn.Close()
			}
		}(pc)
	}
	wg.Wait()
}
