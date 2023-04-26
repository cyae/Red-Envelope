package dbtools

import (
	"fmt"
	"github.com/Shopify/sarama"
	redigo "github.com/garyburd/redigo/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"redEnv_v1/app/redEnv/conftools"
	"redEnv_v1/filepath"
	"strings"
)

// 数据库连接
var Db4Snatch *gorm.DB
var Db4Open *gorm.DB
var Db4Gwl *gorm.DB

// RedisPool redis连接池
var RedisPool *redigo.Pool

// 消息队列连接
var mq4Snatch sarama.AsyncProducer
var mq4Open sarama.AsyncProducer
var mq4Cache sarama.AsyncProducer

//初始化mysql连接、redis连接池
func init() {
	//连接mysql
	mysqlconf := conftools.GetMysqlConfig(fmt.Sprintf("%v%v", filepath.ConfRoot, filepath.MysqlConf))
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?%v", mysqlconf.User, mysqlconf.Password, mysqlconf.Host, mysqlconf.Port, mysqlconf.Db, mysqlconf.Param)
	var err error
	Db4Snatch, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	Db4Open, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	Db4Gwl, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println(err.Error())
	}

	//连接redis
	redisconf := conftools.GetRedisConfig(fmt.Sprintf("%v%v", filepath.ConfRoot, filepath.RedisConf))
	RedisPool = &redigo.Pool{
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.Dial("tcp", fmt.Sprintf("%v:%v", redisconf.Host, redisconf.Port))
			if err != nil {
				return nil, err
			}
			if strings.Compare(redisconf.Password, "") != 0 {
				_, err := c.Do("AUTH", redisconf.Password)
				if err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, nil
		},
		MaxIdle: redisconf.PoolSize,
	}

	//连接kafka
	mqAddr := conftools.GetMqAddr()
	kafkaConf := sarama.NewConfig()
	kafkaConf.Producer.RequiredAcks = sarama.NoResponse

	mq4Snatch, err = sarama.NewAsyncProducer([]string{mqAddr}, kafkaConf)
	mq4Open, err = sarama.NewAsyncProducer([]string{mqAddr}, kafkaConf)
	mq4Cache, err = sarama.NewAsyncProducer([]string{mqAddr}, kafkaConf)

	if err != nil {
		fmt.Println(err.Error())
	}
}
