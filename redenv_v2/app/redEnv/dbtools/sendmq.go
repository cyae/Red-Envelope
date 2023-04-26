package dbtools

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"strconv"
)

func MqSnatch(uid int, eid int, val int, stime int64, N int) {
	type Jsout struct {
		Uid   int   `json:"uid"`
		Eid   int   `json:"eid"`
		Val   int   `json:"val"`
		Stime int64 `json:"stime"`
		N     int   `json:"n"`
	}

	jsout := Jsout{
		Uid:   uid,
		Eid:   eid,
		Val:   val,
		Stime: stime,
		N:     N,
	}

	vb, _ := json.Marshal(jsout)

	msg := &sarama.ProducerMessage{
		Topic: "snatch",
		Value: sarama.ByteEncoder(vb),
	}
	mq4Snatch.Input() <- msg
}

func MqOpen(uid, eid, val int) {
	type Jsout struct {
		Uid int `json:"uid"`
		Eid int `json:"eid"`
		Val int `json:"val"`
	}

	jsout := Jsout{
		Uid: uid,
		Eid: eid,
		Val: val,
	}

	vb, _ := json.Marshal(jsout)

	msg := &sarama.ProducerMessage{
		Topic: "open",
		Value: sarama.ByteEncoder(vb),
	}

	mq4Open.Input() <- msg
}

func MqSaveToCache(uid int) {
	msg := &sarama.ProducerMessage{
		Topic: "cache",
		Value: sarama.ByteEncoder(strconv.Itoa(uid)),
	}

	mq4Cache.Input() <- msg
}
