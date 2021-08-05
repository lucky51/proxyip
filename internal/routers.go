package internal

import "github.com/gin-gonic/gin"

func CreateRouter(pool *proxyPool) *gin.Engine {
	router:=gin.Default()
	router.GET("/",GetIPListHandler(pool))
	router.GET("/ip",GetIPHandler(pool))
	return router
}