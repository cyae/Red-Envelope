package dbtools

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SaveToCache(uid int, conn redis.Conn) {
	var usr User
	var recs []Record
	Db.Where("id = ?", uid).Find(&usr)
	Db.Where("uid = ?", uid).Find(&recs)
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

func SnatchWrite(uid int, eid int, val int, stime int64, N int) {
	conn := RedisPool.Get()
	defer conn.Close()

	//缓存失效
	DelCache(uid, conn)

	//mysql写入 事务
	Db.Transaction(func(tx *gorm.DB) error {
		usr := User{
			Id:    uid,
			Money: 0,
			Cnt:   0,
		}
		// select for update
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", uid).Select("cnt").Find(&usr).Error; err != nil {
			return err
		}
		if usr.Cnt >= N {
			return errors.New("cnt >= N")
		}
		if err := tx.Model(&usr).Where("id = ?", uid).Update("cnt", gorm.Expr("cnt + ?", 1)).Error; err != nil {
			return err
		}
		rec := Record{
			Id:     eid,
			Uid:    uid,
			Val:    val,
			Stime:  stime,
			Opened: 0,
		}
		if err := tx.Create(&rec).Error; err != nil {
			return err
		}
		return nil
	})

	//写入缓存
	SaveToCache(uid, conn)
}

func OpenWrite(uid, eid, val int) {
	conn := RedisPool.Get()
	DelCache(uid, conn)

	Db.Transaction(func(tx *gorm.DB) error {
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

	SaveToCache(uid, conn)
}
