package actor

import (
	"dbass/error"
	"dbass/myRedis"
	"dbass/mysql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

type Actor struct {
	Number int32
	Name   string
	TypeId int32
}

type actorParameter struct {
	Ss     string
	Name   string
	TypeId string
}

func (p *actorParameter) parseActorParameter(c *gin.Context) error {
	if c == nil {
		err := errors.New("gin.Context is nil")
		return err
	}

	err := c.Request.ParseForm()
	if err != nil {
		return err
	}

	form := c.Request.Form
	p.Ss = form.Get("ss")
	p.Name = form.Get("name")
	p.TypeId = form.Get("typeid")

	return nil
}

func (p *actorParameter) createActorRedisKey() (string, error) {
	key, err := json.Marshal(p)

	return string(key), err
}

func (p *actorParameter) createSql() string {
	sqlSentence := "select actor_no,actor_name,actor_typeid from actors"

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
	p.parseActorParameter(c)

	//create redis key and get actors info from redis
	key, err := p.createActorRedisKey()
	if err == nil && key != "" {
		str, err := myRedis.GetRedisInfo(key)
		if err == nil && str != "" {
			c.JSON(200, gin.H{
				"status": "ok",
				"data":   str,
				"method": c.Request.Method,
			})
			return
		}
	}

	/*get mysql instance*/
	db, err := mysql.GetDbInstance()
	if err != nil {
		myError.ReturnErrorMsg(c, err)
		return
	}

	/*create mysql sentence*/
	s := p.createSql()
	fmt.Println("sql:", s)

	/*perform*/
	rows, err := db.Query(s)
	if err != nil {
		myError.ReturnErrorMsg(c, err)
		return
	}
	defer rows.Close()

	var actorNum int32 = 1000
	var actors = make([]Actor, 0, actorNum)
	var count int32 = 0
	var m Actor
	for rows.Next() {
		err := rows.Scan(&m.Number, &m.Name, &m.TypeId)
		if err != nil {
			continue
		}
		count++
		actors = append(actors, m)
		if count >= actorNum {
			break
		}
	}

	if count <= 0 {
		myError.ReturnErrorMsg(c, err)
		return
	}

	message, err := json.Marshal(actors[0:count])
	if err != nil {
		myError.ReturnErrorMsg(c, err)
		return
	}
	c.JSON(200, gin.H{
		"status": "ok",
		"data":   string(message),
		"method": c.Request.Method,
	})

	/*insert actor info to redis*/
	if key != "" {
		go myRedis.SetRedisInfo(key, string(message))
	}
}
