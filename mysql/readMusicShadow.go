package mysql

import (
	"dbass/myRedis"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
)

type MusicShadow struct {
	Id   int
	No   int
	Path string
	Type string
}

func createMusicShadowRedisKey() (string, error) {

	return "MusicShadow", nil
}

func createMusicShadowSql() string {
	return "select id,shadow_no,savepath,music_type from cloud_musicshadow"
}

func ReadMusicShadow(c *gin.Context) {
	//read music shadow info from redis
	key, err := createMusicShadowRedisKey()
	if err == nil && key != "" {
		fmt.Println(key)
		str, err := myRedis.GetRedisInfo(key)
		if err == nil && str != "" {
			c.JSON(200, gin.H{
				"status": "ok",
				"msg":    string(str),
				"method": c.Request.Method,
			})
			return
		}
	}

	/*get mysql instance*/
	db, err := GetDbInstance()
	if err != nil {
		returnMusicShadowErrorMsg(c)
		return
	}

	/*create mysql sentence*/
	s := createMusicShadowSql()
	fmt.Println("sql:", s)

	/*perform*/
	rows, err := db.Query(s)
	if err != nil || rows == nil {
		returnMusicShadowErrorMsg(c)
		return
	}
	defer rows.Close()
	fmt.Printf("rows == %p\n", rows)

	//read music shadow info
	count := 0
	shadow := make([]MusicShadow, 0, 1000)
	for rows.Next() {
		var m MusicShadow
		err := rows.Scan(&m.Id, &m.No, &m.Path, &m.Type)
		if err != nil {
			fmt.Println(err)
			continue
		}
		count++
		shadow = append(shadow, m)
		if count >= 1000 {
			break
		}
	}
	fmt.Println("count=", count)

	message, err := json.Marshal(shadow[0:count])
	if err != nil {
		returnMusicShadowErrorMsg(c)
	}
	c.JSON(200, gin.H{
		"status": "success",
		"msg":    string(message),
		"method": c.Request.Method,
	})

	/*insert medias info to redis*/
	go myRedis.SetRedisInfo(key, string(message))
}

func returnMusicShadowErrorMsg(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "error",
		"msg":    "read muisc shadow fail",
		"method": c.Request.Method,
	})
}
