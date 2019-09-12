package myError

import "github.com/gin-gonic/gin"


func ReturnErrorMsg(c *gin.Context, err error) {
	c.JSON(200, gin.H{
		"status": "error",
		"method": c.Request.Method,
		"URI":    c.Request.RequestURI,
		"Remote": c.Request.RemoteAddr,
		"Err":    err,
	})
}
