package dbtools

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

// QueryFromCache 从缓存中查询红包列表
func QueryFromCache(uid int, conn redis.Conn) ([]Env, error) {
	var envs []Env
	vb, err := redis.Bytes(conn.Do("HGet", uid, "envs"))
	if err != nil {
		return envs, err
	}
	err = json.Unmarshal(vb, &envs)
	if err != nil {
		fmt.Println(err.Error())
	}
	return envs, nil
}

// SaveToCache 从mysql查询再写入到缓存中
func SaveToCache(uid int, conn redis.Conn) {
	var usr User
	var recs []Record
	Db4Gwl.Where("id = ?", uid).Find(&usr)
	Db4Gwl.Where("uid = ?", uid).Find(&recs)
	var envs []Env
	for _, val := range recs {
		envs = append(envs, Env{
			Id:     val.Id,
			Opened: val.Opened,
			Val:    val.Val,
			Stime:  val.Stime,
		})
	}

	vb, err := json.Marshal(envs)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	conn.Do("HSet", uid, "money", usr.Money)
	conn.Do("HSet", uid, "cnt", usr.Cnt)
	conn.Do("HSet", uid, "envs", vb)
}

func DelCache(uid int, conn redis.Conn) {
	conn.Do("del", uid)
}
