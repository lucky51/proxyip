package cmd

import (
	"fmt"
	"github.com/lucky51/proxyip/crawler"
	"github.com/lucky51/proxyip/internal"
	"github.com/spf13/cobra"
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

func init() {
	crawlCmd.Flags().IntVarP(&page,"page","p",1,"crawl page")
}

