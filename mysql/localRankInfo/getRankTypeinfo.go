package localRankInfo

import (
	"dbass/error"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/json"
	)

func GetLocalRankTypeInfo(c *gin.Context) {
	rankType := []struct {
		ModuleId int
		TypeId   int
		TypeName string
	}{
		{4, 1, "总排行"},
		{4, 2, "国语"},
		{4, 3, "粤语"},
		{4, 4, "闽南语"},
		{4, 9, "新歌"},
		{4, 10, "周"},
		{4, 11, "月"},
	}

	data, err := json.Marshal(rankType)
	if err != nil {
		myError.ReturnErrorMsg(c, err)
		return
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"data":    string(data),
		"method":  c.Request.Method,
		"remoute": c.Request.RemoteAddr,
	})
}
