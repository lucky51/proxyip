package cmd

import (
	"fmt"
	"github.com/lucky51/proxyip/crawler"
	"github.com/lucky51/proxyip/internal"
	"github.com/spf13/cobra"
	"net/url"
	"os"
)

var  page int
var crawlCmd =  &cobra.Command{
	Use:"crawl",
	Short:"crawl proxy ip",
	Run: func(cmd *cobra.Command, args []string) {
		var forever = make(chan bool)
		fmt.Println("instance crawler")
		c:=crawler.NewJXLCrawler(page)
		fmt.Println("starting crawl")
		go func() {
			// notify receiver goroutine exit
			defer c.Close()
			c.Crawl()
		}()
		fmt.Println("starting receive proxy ip from channel")
		go func() {
			var proxies = make([]*internal.HttpProxyIP,0)
			for  proxyIp := range c.GetProxies() {
				proxies =append(proxies,proxyIp)
			}
			internal.RenderProxyIPTable(os.Stdout,"crawl proxy ip list",proxies)
			forever<-true
		}()
		<-forever
		fmt.Println("ctrl c exit.")
		select{}
	},
}

var proxyUrl string
var validUrl string

var checkCmd = &cobra.Command{
	Use:"check",
	Short: "check http proxy ip",
	Run: func(cmd *cobra.Command, args []string) {
		if proxyUrl ==""{
			cmd.Println("error:Please input a valid proxy URL")
			return
		}
		_,err:=url.Parse(validUrl)
		if err!=nil{
			cmd.PrintErrln(err)
			return
		}
		result,s,err:=internal.CheckProxyUrl(proxyUrl,validUrl)
		fmt.Printf("result:%s \n http status code:%d \n error: %v \n",result,s,err)
	},
}

func init() {
	crawlCmd.Flags().IntVarP(&page,"page","p",1,"crawl page")
	checkCmd.Flags().StringVarP(&proxyUrl,"proxy","p","","proxy url ex. http://192.168.1.1:80 ")
	checkCmd.Flags().StringVarP(&validUrl,"url","u",internal.DefaultIPCheckUrl,"request a URL to validate  proxy IP,default http://icanhazip.com")
}

