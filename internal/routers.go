package internal

import "github.com/gin-gonic/gin"

func CreateRouter(pool *proxyPool) *gin.Engine {
	router:=gin.Default()
	router.GET("/",IndexHandler())
	router.GET("/proxies",GetIPListHandler(pool))
	router.GET("/ip",GetIPHandler(pool))
	router.GET("/configuration",GetConfigurations(pool))
	return router
}