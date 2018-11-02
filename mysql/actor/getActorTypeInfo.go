package actor

import (
	"dbass/error"
	"dbass/mysql"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"fmt"
)

type actorType struct {
	ActorTypeID   int
	ActorTypeName string
}

func GetActorType(c *gin.Context) {
	sql := "select * from actortypes order by actortype_id"

	db, err := mysql.GetDbInstance()
	if err != nil {
		myError.ReturnErrorMsg(c, err)
		return
	}

	rows, err := db.Query(sql)
	if err != nil {
		myError.ReturnErrorMsg(c, err)
		return
	}
	defer rows.Close()

	fmt.Println(sql)

	count := 0
	var actorTypes []actorType
	var m actorType
	for rows.Next() {
		err := rows.Scan(&m.ActorTypeID, &m.ActorTypeName)
		if err != nil {
			continue
		}
		count++
		actorTypes = append(actorTypes, m)
	}

	data, err := json.Marshal(actorTypes)
	if err != nil {
		myError.ReturnErrorMsg(c, err)
		return
	}

	c.JSON(200, gin.H{
		"status": "ok",
		"data":   string(data),
		"method": c.Request.Method,
	})
}
