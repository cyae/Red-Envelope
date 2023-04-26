package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type SnatchReq struct {
	Uid int `json:"uid"`
}

type Output struct {
	TotalTime int64 `yaml:"total_time(ms)"`
	AvgQps float64 `yaml:"avg_qps"`
	AvgReqTime float64 `yaml:"avg_req_time(ms)"`
	FailedTime int `yaml:"failed_time"`
}

var ReqBytesPool [][][]byte
var respTime int64 = 0
var failedReq int = 0

var Url string
var wg sync.WaitGroup
var ch chan int64
var ch2 chan int


func main() {
	rand.Seed(time.Now().UnixNano())
	// test
	type InitData struct {
		N int `yaml:"n"`
		Ts int `yaml:"ts"`
		Url string `yaml:"url"`
	}
	var initData InitData
	yamlBytes, err := ioutil.ReadFile("init.yml")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = yaml.Unmarshal(yamlBytes, &initData)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	Url = initData.Url

	tn := initData.N / initData.Ts

	//首先构造发送请求的对象池
	for i := 0; i < initData.Ts; i++ {
		var line [][]byte
		for j := 0; j < tn; j++ {
			tmpReq := SnatchReq{Uid: rand.Intn(2000)}
			vb, err := json.Marshal(tmpReq)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			line = append(line, vb)
		}
		ReqBytesPool = append(ReqBytesPool, line)
	}
/*	fmt.Println(cap(ReqBytesPool), len(ReqBytesPool))
	fmt.Println(cap(ReqBytesPool[0]), len(ReqBytesPool[0]))
	fmt.Println(cap(ReqBytesPool[0][0]), len(ReqBytesPool[0][0]))*/
	//构造接受请求的对象池

	ch = make(chan int64, initData.N) //接受每个goroutine 请求本身的总耗时
	ch2 = make(chan int, initData.N)
	//开启线程
	t1 := time.Now().UnixMilli()
	for i := 0; i < initData.Ts; i++ {
		wg.Add(1)
		go SendReq(tn, i)
	}

	go AddTime()
	go AddFailed()
	wg.Wait()
	close(ch)
	close(ch2)
	t2 := time.Now().UnixMilli()

	res := Output{
		TotalTime:  t2 - t1,
		AvgQps:     float64(initData.N) / float64(t2 - t1) * 1000,
		AvgReqTime: float64(respTime) / float64(initData.N),
		FailedTime: failedReq,
	}

	vb, err := yaml.Marshal(res)

	err = ioutil.WriteFile("./test_result.yml", vb, 0666) //写入文件(字节数组)

	if err != nil {
		fmt.Println(err.Error())
	}
	/*fmt.Println("time :", t2 - t1, "ms")
	fmt.Println("total QPS =", float64(initData.N) / float64(t2 - t1) * 1000)
	fmt.Println("req average resp time :", float64(respTime) / float64(initData.N), "ms")
	fmt.Println("failed Time :", failedReq)*/
}

func SendReq(n, tid int) {

	for i := 0; i < n; i++ {
		req, err := http.NewRequest("POST", Url, bytes.NewBuffer(ReqBytesPool[tid][i]))

		client := http.Client{}

		t1 := time.Now().UnixMilli()
		_, err = client.Do(req)
		t2 := time.Now().UnixMilli()

		if err != nil {
			fmt.Println(err.Error())
			ch2 <- 1
		}
		ch <- t2 - t1
	}
	wg.Done()
}

func AddTime() {
	for i := range ch {
		respTime += i
	}
}

func AddFailed() {
	for i := range ch2 {
		failedReq += i
	}
}