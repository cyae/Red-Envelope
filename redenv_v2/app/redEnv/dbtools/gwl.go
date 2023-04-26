package dbtools

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

func GwlGet(uid int) (int, []Env, error) {
	conn := RedisPool.Get()
	defer conn.Close()

	var money int
	var envs []Env
	var err error
	money, err = redis.Int(conn.Do("HGet", uid, "money"))

	//缓存未命中
	if err != nil {
		var recs []Record
		usr := User{}
		rs := Db4Gwl.Where("id = ?", uid).Select("money").Find(&usr)
		if rs.RowsAffected != 0 {
			Db4Gwl.Order("stime").Where("uid = ?", uid).Select("id", "val", "stime", "opened").Find(&recs)
			money = usr.Money
			for _, rec := range recs {
				envs = append(envs, Env{
					Id:     rec.Id,
					Opened: rec.Opened,
					Val:    rec.Val,
					Stime:  rec.Stime,
				})
			}
			MqSaveToCache(uid)
			return money, envs, nil
		} else {
			return money, envs, err
		}
	} else {
		vb, err := redis.Bytes(conn.Do("HGet", uid, "envs"))
		if err != nil {
			fmt.Println("缓存恰好失效，权当没有这个人吧")
			return money, envs, err
		}
		json.Unmarshal(vb, &envs)
		//排序
		//你以为哥不会写快排吗
		quickSort(envs, 0, len(envs))

		return money, envs, nil
	}
}

func quickSort(envs []Env, l int, r int) {
	if l >= r {
		return
	}
	k := envs[l].Stime
	i := l
	j := r
	for i < j {
		for j > i && envs[j].Stime >= k {
			j--
		}
		swap(envs, i, j)
		for i < j && envs[i].Stime <= k {
			i++
		}
		swap(envs, i, j)
	}
	quickSort(envs, l, i-1)
	quickSort(envs, i+1, r)
}

func swap(envs []Env, i int, j int) {
	t := envs[i]
	envs[i] = envs[j]
	envs[j] = t
}
