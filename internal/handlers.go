package internal

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

//var pool = NewProxyPool(&PollingGetProxyIPStrategy{})

func GetIPListHandler(pool *proxyPool) gin.HandlerFunc{
	return func(c *gin.Context) {
		ips,err:=pool.GetAll()
		if err!=nil{
			RenderProxyIPTable(c.Writer,"proxy ip pool",make([]*HttpProxyIP,0))
		}else{
			RenderProxyIPTable(c.Writer,"proxy ip pool",ips)
		}
	}
}

func GetIPHandler(pool *proxyPool) gin.HandlerFunc {
	return func(context *gin.Context) {
		ip,err:=pool.Get()
		if err!=nil{
			context.AbortWithError(http.StatusOK,err)
			return
		}
		context.JSON(200,gin.H{
			"proxy":ip,
		})
	}
}
// GetConfigurations get proxy pool configurations
func GetConfigurations(pool *proxyPool)gin.HandlerFunc {
	return func(c *gin.Context) {
		 config:=*pool.c
		 RenderConfigTable(c.Writer,"proxy ip pool configurations",config)
	}
}

func IndexHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type","text/html; charset=utf-8")
		indexBody:=`
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Proxy IP Pool</title>
    <style>
        /* For mobile phones: */
        [class*="col-"] {
            width: 100%;
        }

        @media only screen and (min-width: 600px) {

            /* For tablets: */
            .col-m-12 {
                width: 8.33%;
            }

            .col-m-2 {
                width: 16.66%;
            }

            .col-m-3 {
                width: 25%;
            }

            .col-m-4 {
                width: 33.33%;
            }

            .col-m-5 {
                width: 41.66%;
            }

            .col-m-6 {
                width: 50%;
            }

            .col-m-7 {
                width: 58.33%;
            }

            .col-m-8 {
                width: 66.66%;
            }

            .col-m-9 {
                width: 75%;
            }

            .col-m-10 {
                width: 83.33%;
            }

            .col-m-11 {
                width: 91.66%;
            }

            .col-m-12 {
                width: 100%;
            }
        }

        @media only screen and (min-width: 768px) {

            /* For desktop: */
            .col-1 {
                width: 8.33%;
            }

            .col-2 {
                width: 16.66%;
            }

            .col-3 {
                width: 25%;
            }

            .col-4 {
                width: 33.33%;
            }

            .col-5 {
                width: 41.66%;
            }

            .col-6 {
                width: 50%;
            }

            .col-7 {
                width: 58.33%;
            }

            .col-8 {
                width: 66.66%;
            }

            .col-9 {
                width: 75%;
            }

            .col-10 {
                width: 83.33%;
            }

            .col-11 {
                width: 91.66%;
            }

            .col-12 {
                width: 100%;
            }
        }

        a {
            display: inline-block;
            height: 100px;
            background-color: #0094ff;
            color: white;
            font-size: large;
            text-align: center;
            line-height: 100px;
            margin-bottom: 20px;
        }
        a:hover{
            border-bottom: yellowgreen;
            color:red;
        }
    </style>
</head>

<body>
    <h1>Proxy IP Pool</h1>
    <hr>
    <a href="/proxies" target="_blank" class="col-m-12 col-3">Proxy IP table</a>
    <a href="/ip" target="_blank" class="col-m-12 col-3">Get Proxy IP by the current strategy</a>
    <a href="/configuration" target="_blank" class="col-m-12 col-3">Get current IP Pool configurations table</a>
</body>

</html>
`
		io.WriteString(c.Writer,indexBody)
	}
}