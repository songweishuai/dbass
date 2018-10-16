package http

import (
	"dbass/mysql"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func CreateHttpWeb() error {
	router := CreateRouter()
	s := http.Server{
		Addr:           "192.168.1.97:3233",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err := s.ListenAndServe()

	return err
}

func CreateRouter() *gin.Engine {
	router := gin.Default()

	//router.GET("/medias:ss/:name", mysql.ReadMedias)
	router.GET("/medias", mysql.ReadMedias)
	router.GET("/actors", mysql.ReadActors)
	router.GET("/musicShadow",mysql.ReadMusicShadow)

	return router
}
