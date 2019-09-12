package serverInfo

import (
	"dbass/error"
	"dbass/mysql"
	"encoding/json"
		"github.com/gin-gonic/gin"
)

type serverInfo struct {
	ServerId int
	GroupId  int
	Type     int
	IP       string
}

func GetServerInfo(c *gin.Context) {
	sql := "select server_id,server_grpid,server_weight,server_ip from servers"
	db, err := mysql.GetDbInstance()
	if err != nil {
		myError.ReturnErrorMsg(c,err)
		return
	}

	rows, err := db.Query(sql)
	if err != nil {
		myError.ReturnErrorMsg(c, err)
		return
	}
	defer rows.Close()

	var info = make([]serverInfo, 0, 20)
	var val serverInfo
	count := 0
	for rows.Next() {
		err := rows.Scan(&val.ServerId, &val.GroupId, &val.Type, &val.IP)
		if err != nil {
			continue
		}
		count++
		info = append(info, val)
	}

	if count <= 0 {
		myError.ReturnErrorMsg(c,err)
	}

	msg, err := json.Marshal(info)
	if err != nil {
		myError.ReturnErrorMsg(c,err)
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"method": c.Request.Method,
		"msg":    string(msg),
	})
}
