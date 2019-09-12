package http

import (

	"github.com/gin-gonic/gin"
	"mysql"
	"mysql/actor"
	"mysql/localRankInfo"
	"mysql/roomInfo"
	"mysql/serverInfo"
	"net/http"
	"time"
)

func CreateHttpWeb() error {
	router := CreateRouter()
	s := http.Server{
		Addr:           "192.168.1.96:3233",
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

	router.GET("/serverInfo", serverInfo.GetServerInfo)
	router.GET("/roomSkinInfo", roomInfo.GetRoomSkinInfo)

	//router.GET("/medias:ss/:name", mysql.ReadMedias)
	router.GET("/medias", mysql.ReadMedias)

	router.GET("/actors", actor.ReadActors)
	router.GET("/actorType", actor.GetActorType)

	router.GET("/musicShadow", mysql.ReadMusicShadow)

	router.GET("/localRankType", localRankInfo.GetLocalRankTypeInfo)
	router.GET("/localRankTypeMedia", localRankInfo.GetLocalRankTypeMedia)

	return router
}
