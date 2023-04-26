package dbtools

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SnatchGet 通过uid查询用户已抢红包数量
func SnatchGet(uid int) (int, error, bool) {
	conn := RedisPool.Get()
	defer conn.Close()

	//查询缓存
	cnt, err := redis.Int(conn.Do("HGet", uid, "cnt"))

	//缓存未命中
	if err != nil {
		var usr User
		rs := Db4Snatch.Where("id = ?", uid).Select("cnt").Find(&usr)
		if rs.RowsAffected != 0 {
			err = nil
			cnt = usr.Cnt
			//直接发到消息队列写入缓存 但是抢到的这一次写入多余了 后期简单优化即可
			return usr.Cnt, nil, false
		}
		return 0, err, false
	}
	return cnt, err, true
}

// SnatchWrite 消息队列的消费者 更新mysql以及redis
func SnatchWrite(uid int, eid int, val int, stime int64, N int) {
	conn := RedisPool.Get()
	defer conn.Close()

	//缓存失效
	DelCache(uid, conn)

	//mysql写入 事务
	Db4Snatch.Transaction(func(tx *gorm.DB) error {
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
