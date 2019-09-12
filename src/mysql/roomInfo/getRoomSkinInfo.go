package roomInfo

import (
	"dbass/error"
	"dbass/mysql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/json"
)

type roomSkinInfo struct {
	SkinId    int
	SkinNname string
}

type roomSkinInfoParamter struct {
	RoomIP string
}

func paraseRoomSkinParamter(c *gin.Context, p *roomSkinInfoParamter) error {
	if c == nil {
		return errors.New("this gin.Context is nil")
	}

	//fmt.Println("c.Request.RemoteAddr:",c.Request.RemoteAddr)
	//fmt.Println("c.Request.Host:",c.Request.Host)
	//fmt.Println("c.Request.RequestURI:",c.Request.RequestURI)

	err := c.Request.ParseForm()
	if err != nil {
		return err
	}

	form := c.Request.Form
	p.RoomIP = form.Get("ip")
	fmt.Println("roomip:", p.RoomIP)
	return nil
}

func createRoomSkinSql(p *roomSkinInfoParamter) (string, error) {
	sql := "select skin_id,skin_name from skins s inner join rooms r on r.room_ip='%s' and s.skin_id=r.room_skin"

	if p == nil {
		return "", errors.New("this roomSkinInfoParamter* is nil")
	}

	sql = fmt.Sprintf(sql, p.RoomIP)

	return sql, nil
}

func GetRoomSkinInfo(c *gin.Context) {
	var p roomSkinInfoParamter

	err := paraseRoomSkinParamter(c, &p)
	if err != nil {
		fmt.Println(err)
		myError.ReturnErrorMsg(c, err)
		return
	}

	sql, err := createRoomSkinSql(&p)
	if err != nil {
		fmt.Println(err)
		myError.ReturnErrorMsg(c, err)
		return
	}

	db, err := mysql.GetDbInstance()
	if err != nil {
		fmt.Println(err)
		myError.ReturnErrorMsg(c, err)
		return
	}

	rows, err := db.Query(sql)
	if err != nil {
		fmt.Println(err)
		myError.ReturnErrorMsg(c, err)
		return
	}
	defer rows.Close()

	var info roomSkinInfo
	count := 0
	for rows.Next() {
		rows.Scan(&info.SkinId, &info.SkinNname)
		count++
		break
	}

	if count != 1 {
		myError.ReturnErrorMsg(c, err)
		return
	}

	msg, err := json.Marshal(info)
	if err != nil {
		fmt.Println(err)
		myError.ReturnErrorMsg(c, err)
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data":    string(msg),
		"method": c.Request.Method,
	})
}
