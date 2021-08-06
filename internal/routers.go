package internal

import "github.com/gin-gonic/gin"

func CreateRouter(pool *proxyPool) *gin.Engine {
	router:=gin.Default()
	router.GET("/",HandleGetIPIndexPage())
	router.GET("/proxies",HandleGetIPList(pool))
	router.GET("/ip",HandleGetIP(pool))
	router.GET("/configuration",HandleGetConfigurations(pool))
	return router
}