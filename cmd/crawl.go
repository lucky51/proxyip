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

		pool:=internal.NewProxyPool(&internal.PollingGetProxyIPStrategy{})
		go func() {
			defer c.Close()
			c.Crawl(page)
		}()
		go func() {
			for proxy := range c.GetProxies() {
				pool.Set(proxy)
			}
			forever<-true
		}()
		fmt.Println("starting get proxy ip from pool")
		<-forever
		var data = make([]*internal.HttpProxyIP,4)
		for i:=0;i<4;i++ {
			item,err:=pool.Get()
			if err!=nil{
				break
			}
			data[i] = item
		}
		internal.RenderTable(os.Stdout,"polling get proxy ip",data)

	},
}

func init() {
	crawlCmd.Flags().IntVarP(&page,"page","p",1,"crawl page")
}

