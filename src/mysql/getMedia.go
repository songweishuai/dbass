package mysql

import (
	"dbass/error"
	"dbass/myRedis"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

//type media struct {
//	Name       string `json:"Name"`
//	Number     int32  `json:"SerialNo"`
//	ActorName1 string `json:"ActorName"`
//}

type Media struct {
	MediaNo         int
	MediaName       string
	MediaLangtype   int
	MediaLangid     int
	MediaTag1       string
	MediaActName1   string
	MediaActName2   string
	MediaYuan       int
	MediaBan        int
	MediaSvrGroup   int
	MediaFile       string
	MediaStyle      string
	MediaVolume     int
	MediaStars      int
	MediaActNo1     int
	MediaActNo2     int
	MediaDafen      int
	MediaClimax     int
	MediaClimaxInfo string
	MediaYinYi      int
	MediaLight      int
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

func (p *mediaParameter) parseMediaParameter(c *gin.Context) error {
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
	p.No = form.Get("name") /*can recv json array*/
	p.Stroke = form.Get("stroke")
	p.Language = form.Get("language")
	p.Hot = form.Get("hot")
	p.New = form.Get("new")
	p.ActorId = form.Get("actorId")
	p.Len = form.Get("len")

	return nil
}

func (p *mediaParameter) createMediaRedisKey() (string, error) {
	key, err := json.Marshal(*p)

	return string(key), err
}

func (p *mediaParameter) CreateSql(c *gin.Context) string {
	/*parse media search parameter*/
	sqlSentence := "select media_no,media_name,media_langtype,media_langid,media_tag1,media_actname1,media_actname2," +
		"media_yuan,media_ban,media_svrgroup,media_file,media_style,media_volume,media_stars,media_actno1,media_actno2," +
		"media_dafen,media_climax,media_climaxinfo,media_yinyi,media_light" +
		" from medias"

	err := c.Request.ParseForm()
	if err != nil {
		sqlSentence += "limit 10"
		return sqlSentence
	}

	fmt.Println(c.Request.Form)

	var isWhere = false
	//if p.No {
	//	//select * from medias where media_no in (7300616,8434232)
	//}

	if p.Len != "" {
		var s string
		if isWhere {
			s = fmt.Sprintf(" and media_namelen=%s", p.Len)
		} else {
			s = fmt.Sprintf(" where media_namelen=%s", p.Len)
		}
		//s := fmt.Sprintf(" where media_namelen=%s", p.Len)
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
	err := p.parseMediaParameter(c)

	/*create redis media key*/
	key, err := p.createMediaRedisKey()
	if err == nil && key != "" {
		/*read medias info from redis*/
		str, err := myRedis.GetRedisInfo(key)
		if err == nil && str != "" {
			/*return search msg*/
			c.JSON(200, gin.H{
				"status": "ok",
				"data":   str,
				"method": c.Request.Method,
			})
			return
		}
	}

	/*mysql sentence*/
	s := p.CreateSql(c)

	/*get db*/
	db, err := GetDbInstance()
	if err != nil {
		myError.ReturnErrorMsg(c, err)
	}

	/*perform mysql sentence*/
	row, err := db.Query(s)
	if err != nil {
		myError.ReturnErrorMsg(c, err)
	}
	defer row.Close()

	/*read medias*/
	mediaNum := 300
	mediaArray := make([]Media, 0, mediaNum)
	var m Media
	i := 0
	for row.Next() {
		err := row.Scan(&m.MediaNo, &m.MediaName, &m.MediaLangtype, &m.MediaLangid, &m.MediaTag1, &m.MediaActName1, &m.MediaActName2,
			&m.MediaYuan, &m.MediaBan, &m.MediaSvrGroup, &m.MediaFile, &m.MediaStyle, &m.MediaVolume, &m.MediaStars, &m.MediaActNo1,
			&m.MediaActNo2, &m.MediaDafen, &m.MediaClimax, &m.MediaClimaxInfo, &m.MediaYinYi, &m.MediaLight)
		if err != nil {
			continue
		}
		i++
		mediaArray = append(mediaArray, m)
		if i >= mediaNum {
			break
		}
	}

	/*to json*/
	message, err := json.Marshal(mediaArray[0:i])
	if err != nil {
		myError.ReturnErrorMsg(c, err)
		return
	}


	/*return search msg*/
	c.JSON(200, gin.H{
		"status": "ok",
		"data":   string(message),
		"method": c.Request.Method,
	})

	/*insert medias info to redis*/
	if key != "" {
		go myRedis.SetRedisInfo(key, string(message))
	}
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
