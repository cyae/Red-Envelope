package tools

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

// GetPort 通过yaml返回服务器运行端口
func GetPort(filepath string) int {
	type Config struct {
		Port int `yaml:"port"`
	}
	var conf Config
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = yaml.Unmarshal(file, &conf)
	if err != nil {
		fmt.Println(err.Error())
	}
	return conf.Port
}

func GetLimiter(filepath string) (time.Duration, int64) {
	type Config struct {
		BucketFillDuring int64 `yaml:"bucket_fill_during"`
		BucketMax int64 `yaml:"bucket_max"`
	}
	var conf Config
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = yaml.Unmarshal(file, &conf)
	if err != nil {
		fmt.Println(err.Error())
	}
	return time.Duration(conf.BucketFillDuring) * time.Millisecond, conf.BucketMax
}