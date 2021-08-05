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
		c:=crawler.NewJXLCrawler()

		defer c.Close()
		go c.Crawl()
		go func() {
			var proxies = make([]*internal.HttpProxyIP,0)
			for  proxyIp := range c.GetProxies() {
				proxies =append(proxies,proxyIp)
			}
			internal.RenderTable(os.Stdout,"crawl proxy ip list",proxies)
			forever<-true
		}()

		<-forever
		fmt.Println("blocking this,ctrl c exit.")
		select{}
	},
}

func init() {
	crawlCmd.Flags().IntVarP(&page,"page","p",1,"crawl page")
}

