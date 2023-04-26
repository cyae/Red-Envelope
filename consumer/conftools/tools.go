package conftools

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type MysqlConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Db       string
	Param    string
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	PoolSize int    `yaml:"pool_size"`
	Password string `yaml:"password"`
}

// GetMysqlConfig 通过yaml返回mysql配置信息
func GetMysqlConfig(filepath string) MysqlConfig {
	var conf MysqlConfig
	/*	file, err := ioutil.ReadFile(filepath)
		if err != nil {
			fmt.Println(err.Error())
		}
		err = yaml.Unmarshal(file, &conf)
		if err != nil {
			fmt.Println(err.Error())
		}*/

	conf.User = os.Getenv("MYSQL_USER")
	conf.Password = os.Getenv("MYSQL_PASSWORD")
	conf.Host = os.Getenv("MYSQL_HOST")
	if conf.User == "" || conf.Password == "" || conf.Host == "" {
		fmt.Println("cant get mysql env var")
	}
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println(err.Error())
	}
	type MysqlYaml struct {
		Port  string `yaml:"port"`
		Db    string `yaml:"db"`
		Param string `yaml:"param"`
	}

	var ymlConf MysqlYaml
	err = yaml.Unmarshal(file, &ymlConf)
	if err != nil {
		fmt.Println(err.Error())
	}

	conf.Db = ymlConf.Db
	conf.Param = ymlConf.Param
	conf.Port = ymlConf.Port
	return conf
}

// GetRedisConfig 通过yaml返回redis配置信息
func GetRedisConfig(filepath string) RedisConfig {
	var conf RedisConfig
	conf.Host = os.Getenv("REDIS_HOST")
	conf.Password = os.Getenv("REDIS_PASSWORD")

	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println(err.Error())
	}

	type RedisYaml struct {
		Port     string `yaml:"port"`
		PoolSize int    `yaml:"pool_size"`
	}

	var ymlConf RedisYaml
	err = yaml.Unmarshal(file, &ymlConf)
	if err != nil {
		fmt.Println(err.Error())
	}
	conf.Port = ymlConf.Port
	conf.PoolSize = ymlConf.PoolSize
	return conf
}
