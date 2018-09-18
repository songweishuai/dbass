package mysql

import (
	"dbass/myRedis"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
)

type media struct {
	Name       string `json:"Name"`
	Number     int32  `json:"SerialNo"`
	ActorName1 string `json:"ActorName"`
}

type mediaParameter struct {
	Ss        string
	Name      string
	No        string
	Stroke    string
	StrokeNum string
	Language  string
	Hot       string
	New       string
	ActorId   string
	Len       string
}

func parseMediaParameter(c *gin.Context, p *mediaParameter) error {
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
	p.No = form.Get("name") /*can recv json array*/
	p.Stroke = form.Get("stroke")
	p.Language = form.Get("language")
	p.Hot = form.Get("hot")
	p.New = form.Get("new")
	p.ActorId = form.Get("actorId")
	p.Len = form.Get("len")

	return nil
}

func createMediaRedisKey(p *mediaParameter) (string, error) {
	if p == nil {
		err := errors.New("parameter p is nil")
		return "", err
	}

	key, err := json.Marshal(*p)

	return string(key), err
}

func CreateSql(c *gin.Context) string {
	/*parse media search parameter*/
	var p mediaParameter
	//sqlSentence := "select media_no,media_langtype,media_style,media_svrgroup,media_name,media_stars,media_langid,media_actname1,media_actname2,media_actno1,media_click,media_dafen,media_carria from medias"
	sqlSentence := "select media_name,media_no,media_actname1 from medias"

	err := c.Request.ParseForm()
	if err != nil {
		sqlSentence += "limit 10"
		return sqlSentence
	}

	fmt.Println(c.Request.Form)

	form := c.Request.Form

	p.Ss = form.Get("ss")
	p.Name = form.Get("name")
	p.No = form.Get("name") /*can recv json array*/
	p.Stroke = form.Get("stroke")
	p.Language = form.Get("language")
	p.Hot = form.Get("hot")
	p.New = form.Get("new")
	p.ActorId = form.Get("actorId")
	p.Len = form.Get("len")

	//fmt.Println(p.ss, p.name, p.no, p.stroke, p.language, p.hot, p.new, p.actorId)

	var isWhere = false
	if p.Len != "" {
		s := fmt.Sprintf(" where media_namelen=%s", p.Len)
		isWhere = true
		sqlSentence += s
	}

	if p.Language != "" {
		var s string
		if isWhere {
			s = fmt.Sprintf(" and media_langid=%s", p.Language)
		} else {
			s = fmt.Sprintf(" where media_langid=%s", p.Language)
		}
		isWhere = true
		sqlSentence += s
	}

	if p.New != "" {
		var s string
		if isWhere {
			s = fmt.Sprintf(" and media_isnew=%s", p.New)

		} else {
			s = fmt.Sprintf(" where media_isnew=%s", p.New)
		}
		isWhere = true
		sqlSentence += s
	}

	if p.StrokeNum != "" {
		var s string
		if isWhere {
			s = fmt.Sprintf(" and media_stroke=%s", p.StrokeNum)
		} else {
			s = fmt.Sprintf(" wher media_stroke=%s", p.StrokeNum)
		}
		isWhere = true
		sqlSentence += s
	}

	if p.Hot != "" {
		var s string
		if isWhere {
			s = fmt.Sprintf(" and media_stars=%s", p.Hot)
		} else {
			s = fmt.Sprintf(" where media_stars=%s", p.Hot)
		}
		isWhere = true
		sqlSentence += s
	}

	if p.ActorId != "" {
		var s string
		if isWhere {
			s = fmt.Sprintf(" and media_actno1=%s", p.ActorId)
		} else {
			s = fmt.Sprintf(" where media_actno1=%s", p.ActorId)
		}
		isWhere = true
		sqlSentence += s
	}

	if p.Ss != "" {
		var s string
		if isWhere {
			//s = fmt.Sprintf(" and media_jp like \"%s%s\"", p.ss, "%")
			s = fmt.Sprintf(" and match(`media_jp`) against(\"%s*\" in boolean mode)", p.Ss)
		} else {
			//s = fmt.Sprintf(" where media_jp like \"%s%s\"", p.ss, "%")
			s = fmt.Sprintf(" where match(`media_jp`) against(\"%s*\" in boolean mode)", p.Ss)
		}
		isWhere = true
		sqlSentence += s
	}

	if p.Stroke != "" {
		var s string
		if isWhere {
			//s = fmt.Sprintf(" and media_stroks like \"%s%s\"", p.stroke, "%")
			s = fmt.Sprintf(" and match(`media_stroks`) against(\"%s\" in boolean mode)", p.Stroke)
		} else {
			//s = fmt.Sprintf(" where media_stroks like \"%s%s\"", p.stroke, "%")
			s = fmt.Sprintf(" where match(`media_stroks`) against(\"%s\" in boolean mode)", p.Stroke)
		}
		isWhere = true
		sqlSentence += s
	}

	sqlSentence += " limit 10"
	fmt.Println(sqlSentence)

	return sqlSentence
}

func ReadMedias(c *gin.Context) {
	//parse parameter
	var p mediaParameter
	err := parseMediaParameter(c, &p)

	/*create redis media key*/
	key, err := createMediaRedisKey(&p)
	if err == nil && key != "" {
		/*read medias info from redis*/
		str, err := myRedis.GetRedisInfo(key)
		if err == nil && str != "" {
			/*return search msg*/
			c.JSON(200, gin.H{
				"status": "ok",
				"msg":    str,
				"method": c.Request.Method,
			})
			return
		}
	}

	/*mysql sentence*/
	s := CreateSql(c)

	/*get db*/
	db, err := GetDbInstance()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	/*perform mysql sentence*/
	row, err := db.Query(s)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer row.Close()

	/*read medias*/
	var mediaArray []media
	i := 0
	for row.Next() {
		var m media
		err := row.Scan(&m.Name, &m.Number, &m.ActorName1)
		if err != nil {
			continue
		}
		//fmt.Println(m.name, m.number, m.actorName1)
		mediaArray = append(mediaArray, m)
		i++
	}
	//fmt.Println("i :", i, mediaArray)

	//mediaArray := make([]media, 12)
	//i := 0
	//for row.Next() {
	//	row.Scan(&mediaArray[i].Name, &mediaArray[i].Number, &mediaArray[i].ActorName1)
	//	if err != nil {
	//		continue
	//	}
	//	i++
	//}

	/*to json*/
	message, err := json.Marshal(mediaArray[0:i])
	if err != nil {
		c.JSON(200, gin.H{
			"statue": "error",
			"msg":    "read medias error",
			"method": c.Request.Method,
		})
		return
	}
	//fmt.Printf("%s\n", message)

	/*insert medias info to redis*/
	go myRedis.SetRedisInfo(key, string(message))

	/*return search msg*/
	c.JSON(200, gin.H{
		"status": "ok",
		"msg":    string(message),
		"method": c.Request.Method,
	})

	//bytes.NewBuffer()
	//c.String(200,string(message))
}

//type media struct {
//	MediaName     string
//	MediaSerialNu string
//}
//func ReadMedias(c *gin.Context) {
//
//	mediaArray := make([]media, 10)
//	for i := 0; i < 10; i++ {
//		mediaArray[i].MediaName = "media" + strconv.Itoa(i)
//		mediaArray[i].MediaSerialNu = strconv.Itoa(7000000 + i)
//
//		fmt.Println(mediaArray[i].MediaName)
//		fmt.Println(mediaArray[i].MediaSerialNu)
//	}
//
//	message, err := json.Marshal(mediaArray)
//	if err != nil {
//		c.JSON(200, gin.H{
//			"statue": "error",
//			"msg":    "read medias error",
//			"method": c.Request.Method,
//		})
//		return
//	}
//	fmt.Println(string(message))
//
//	c.JSON(200, gin.H{
//		"status": "ok",
//		"msg":    string(message),
//		"method": c.Request.Method,
//	})
//
//	//c.AsciiJSON(200, gin.H{
//	//	"status": "ok",
//	//	"msg":    string(message),
//	//	"method": c.Request.Method,
//	//})
//}
