package mysql

import (
	"dbass/myRedis"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

type Actor struct {
	Number int32
	Name   string
	Click  int32
}

type actorParameter struct {
	Ss     string
	Name   string
	TypeId string
}

func parseActorParameter(c *gin.Context, p *actorParameter) error {
	if c == nil {
		err := errors.New("gin.Context is nil")
		return err
	}

	if p == nil {
		err := errors.New("actorParameter p is nil")
		return err
	}

	err := c.Request.ParseForm()
	if err != nil {
		err1 := errors.New("gin.Contest parse error")
		return err1
	}

	form := c.Request.Form
	p.Ss = form.Get("ss")
	p.Name = form.Get("name")
	p.TypeId = form.Get("typeid")

	return nil
}

func createActorRedisKey(p *actorParameter) (string, error) {
	if p == nil {
		err := errors.New("parameter p is nil")
		return "", err
	}

	key, err := json.Marshal(*p)

	return string(key), err
}

func createSql(p *actorParameter) string {
	/*parse media search parameter*/
	//var p actorParameter
	//sqlSentence := "select media_no,media_langtype,media_style,media_svrgroup,media_name,media_stars,media_langid,media_actname1,media_actname2,media_actno1,media_click,media_dafen,media_carria from medias"
	sqlSentence := "select actor_no,actor_name,actor_click from actors"

	//err := c.Request.ParseForm()
	//if err != nil {
	//	sqlSentence += "limit 1000"
	//	return sqlSentence
	//}

	if p == nil {
		sqlSentence += "limit 1000"
		return sqlSentence
	}

	//form := c.Request.Form
	//p.ss = form.Get("ss")
	//p.name = form.Get("name")
	//p.typeId = form.Get("typeid")

	var isWhere = false
	if p.Ss != "" {
		s := fmt.Sprintf(" where match(`actor_jp`) against(\"%s\" in boolean mode)", p.Ss)
		sqlSentence += s
		isWhere = true
	}

	if p.Name != "" {
		var s string
		if isWhere {
			s = fmt.Sprintf(" and match(`actor_name`) against(\"+%s*\" in boolean mode)", p.Name)
		} else {
			s = fmt.Sprintf(" where match(`actor_name`) against(\"+%s*\" in boolean mode)", p.Name)
		}
		sqlSentence += s
		isWhere = true
	}

	if p.TypeId != "" {
		var s string
		if isWhere {
			s = fmt.Sprintf(" and actor_typeid=%s", p.TypeId)
		} else {
			s = fmt.Sprintf(" where actor_typeid=%s", p.TypeId)
		}
		sqlSentence += s
		isWhere = true
	}
	sqlSentence += " limit 1000"

	return sqlSentence
}

func ReadActors(c *gin.Context) {
	//parse parameter
	var p actorParameter
	parseActorParameter(c, &p)

	//create redis key
	key, err := createActorRedisKey(&p)

	//get actors info from redis
	if err == nil && key != "" {
		fmt.Println("key:", key)
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
		returnErrorMsg(c)
	}

	/*create mysql sentence*/
	s := createSql(&p)
	fmt.Println("sql:", s)

	/*perform*/
	rows, err := db.Query(s)
	if err != nil {
		returnErrorMsg(c)
	}
	defer rows.Close()

	//types ,err:= rows.ColumnTypes()
	//if err!=nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(*types[0])
	//fmt.Println(rows.ColumnTypes())

	//var actors =make([]struct{
	//	ActorName string
	//	ActorNo   int32
	//	ActorID   int32
	//},100)
	var actorNum int32 = 1000
	var actors = make([]Actor, 0, actorNum)
	var count int32 = 0
	for rows.Next() {
		err := rows.Scan(&actors[count].Number, &actors[count].Name, &actors[count].Click)
		if err != nil {
			continue
		}
		if count >= actorNum-1 {
			break
		}
		count++
	}
	if count <= 0 {
		returnErrorMsg(c)
	}

	message, err := json.Marshal(actors[0:count])
	if err != nil {
		returnErrorMsg(c)
	}
	c.JSON(200, gin.H{
		"status": "ok",
		"msg":    string(message),
		"method": c.Request.Method,
	})

	/*insert actor info to redis*/
	if key != "" {
		go myRedis.SetRedisInfo(key, string(message))
	}
}

func returnErrorMsg(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "error",
		"msg":    "read actor fail",
		"method": c.Request.Method,
	})
}
