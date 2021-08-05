package internal

import "github.com/gin-gonic/gin"

//var pool = NewProxyPool(&PollingGetProxyIPStrategy{})

func GetIPListHandler(pool *proxyPool) gin.HandlerFunc{
	return func(c *gin.Context) {
		ips,err:=pool.GetAll()
		if err!=nil{
			RenderTable(c.Writer,"proxy ip pool",make([]*HttpProxyIP,0))
		}else{
			RenderTable(c.Writer,"proxy ip pool",ips)
		}
	}
}

func GetIPHandler(pool *proxyPool) gin.HandlerFunc {
	return func(context *gin.Context) {
		ip,err:=pool.Get()
		if err!=nil{
			context.AbortWithError(500,err)
			return
		}
		context.JSON(200,gin.H{
			"proxy":ip,
		})
	}
}

