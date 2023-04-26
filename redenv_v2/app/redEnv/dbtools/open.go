package dbtools

import (
	"encoding/json"
	"errors"
	"github.com/garyburd/redigo/redis"
	"gorm.io/gorm"
)

// OpenGet 根据uid与eid查询是否有该红包 若有则返回opened, val
func OpenGet(uid, eid int) (int, int, error, bool) {
	conn := RedisPool.Get()
	defer conn.Close()
	var opened int
	var val int

	//查询缓存
	vb, err := redis.Bytes(conn.Do("HGet", uid, "envs"))

	// 根据缓存是否命中
	if err == nil {
		var envs []Env
		json.Unmarshal(vb, &envs)
		flag := false
		for _, env := range envs {
			if env.Id == eid {
				flag = true
				opened = env.Opened
				val = env.Val
				break
			}
		}
		if flag == false {
			return opened, val, errors.New(""), true
		} else {
			return opened, val, nil, true
		}
	} else {
		var rec Record
		rs := Db4Open.Where("id = ?", eid).Select("opened", "val").Find(&rec)
		if rs.RowsAffected != 0 {
			err = nil
			opened = rec.Opened
			val = rec.Val
			return rec.Opened, rec.Val, nil, false
		}
		return 0, 0, err, true
	}
}

func OpenWrite(uid, eid, val int) {
	conn := RedisPool.Get()
	DelCache(uid, conn)

	Db4Open.Transaction(func(tx *gorm.DB) error {
		rec := Record{Id: eid}
		if err := tx.Where("id = ?", eid).Select("opened").Find(&rec).Error; err != nil {
			return err
		}

		if rec.Opened == 1 {
			return errors.New("")
		}

		if err := tx.Model(&rec).Where("id = ?", eid).Update("opened", 1).Error; err != nil {
			return err
		}
		usr := User{
			Id:    uid,
			Money: 0,
			Cnt:   0,
		}
		if err := tx.Model(&usr).Where("id = ?", uid).Update("money", gorm.Expr("money + ?", val)).Error; err != nil {
			return err
		}

		return nil
	})

	go SaveToCache(uid, conn)
}
