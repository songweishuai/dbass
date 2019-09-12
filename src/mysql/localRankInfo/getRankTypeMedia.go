package localRankInfo

import (
	"dbass/error"
	"dbass/myRedis"
	"dbass/mysql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
)

//var rankTypeInfo = map[int]string{
//	1:  "总排行",
//	2:  "国语",
//	3:  "粤语",
//	4:  "闽南语",
//	9:  "新歌",
//	10: "周",
//	11: "月",
//}
type parameter struct {
	TypeID string `json:"LocalRankTypeID"`
}

func (p *parameter) parseParameter(c *gin.Context) {
	err := c.Request.ParseForm()
	if err != nil {
		return
	}

	form := c.Request.Form
	p.TypeID = form.Get("typeID")
}

func (p *parameter) createRedisKey() string {
	if p.TypeID == "" {
		return ""
	}
	key, err := json.Marshal(p)
	if err != nil {
		return ""
	}
	return string(key)
}

func (p *parameter) createSql() string {
	sql := "select media_no,media_name,media_langtype,media_langid,media_tag1,media_actname1,media_actname2," +
		"media_yuan,media_ban,media_svrgroup,media_file,media_style,media_volume,media_stars,media_actno1,media_actno2," +
		"media_dafen,media_climax,media_climaxinfo,media_yinyi,media_light from medias %s limit 100"
	fmt.Println("p.TypeID:", p.TypeID)
	switch p.TypeID {
	case "1":
		sql = fmt.Sprintf(sql, "order by media_click desc")
	case "2":
		sql = fmt.Sprintf(sql, "where media_lang='国语' order by media_click desc")
	case "3":
		sql = fmt.Sprintf(sql, "where media_lang='粤语' order by media_click desc")
	case "4":
		sql = fmt.Sprintf(sql, "where media_lang='闽南语' order by media_click desc")
	case "9":
		sql = fmt.Sprintf(sql, "where media_isnew=1 order by media_click desc")
	case "10":
		sql = fmt.Sprintf(sql, "where order by media_clickw desc")
	case "11":
		sql = fmt.Sprintf(sql, "where order by media_clickm desc")
	default:
		return ""
	}
	return sql
}

func GetLocalRankTypeMedia(c *gin.Context) {
	var p parameter

	//parse request parameter
	p.parseParameter(c)

	//create redis key
	key := p.createRedisKey()
	redisMsg, err := myRedis.GetRedisInfo(key)
	if redisMsg != "" && err == nil {
		c.JSON(200, gin.H{
			"status": "ok",
			"data":   redisMsg,
			"method": c.Request.Method,
		})
		return
	}

	//create sql
	sql := p.createSql()
	fmt.Println("sql:", sql)

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

	num := 100
	count := 0
	medias := make([]mysql.Media, 0, num)
	var m mysql.Media
	for rows.Next() {
		err := rows.Scan(&m.MediaNo, &m.MediaName, &m.MediaLangtype, &m.MediaLangid, &m.MediaTag1, &m.MediaActName1, &m.MediaActName2,
			&m.MediaYuan, &m.MediaBan, &m.MediaSvrGroup, &m.MediaFile, &m.MediaStyle, &m.MediaVolume, &m.MediaStars, &m.MediaActNo1,
			&m.MediaActNo2, &m.MediaDafen, &m.MediaClimax, &m.MediaClimaxInfo, &m.MediaYinYi, &m.MediaLight)
		if err != nil {
			continue
		}
		count++
		medias = append(medias, m)
		if count >= num {
			break
		}
	}

	data, err := json.Marshal(medias[0:count])
	if err != nil {
		myError.ReturnErrorMsg(c, err)
		return
	}

	c.JSON(200, gin.H{
		"status": "ok",
		"data":   string(data),
		"method": c.Request.Method,
	})

	if key != "" {
		go myRedis.SetRedisInfo(key, string(data))
	}
}
