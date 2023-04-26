package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"redEnv_v1/app/redEnv/dbtools"
	"redEnv_v1/app/redEnv/statuscode"
)

func GwlHandler(c *gin.Context) {
	available := TokenBucket.TakeAvailable(1)
	if available <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": statuscode.FlowOverrun,
		})
		return
	}

	type Jsin struct {
		Uid int `json:"uid"`
	}
	var jsin Jsin
	err2 := c.ShouldBindJSON(&jsin)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	money, envs, err := dbtools.GwlGet(jsin.Uid)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": statuscode.NoSuchUser,
			"msg": "cant find this user",
		})
		return
	}

	var envList []gin.H
	for _, env := range envs {
		if env.Opened == 1 {
			envList = append(envList, gin.H{
				"envelope_id": env.Id,
				"value": env.Val,
				"opened": true,
				"snatch_time": env.Stime,
			})
		} else {
			envList = append(envList, gin.H{
				"envelope_id": env.Id,
				"opened": false,
				"snatch_time": env.Stime,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": statuscode.OK,
		"msg": "success",
		"data": gin.H{
			"amount": money,
			"envelope_list": envList,
		},
	})
}