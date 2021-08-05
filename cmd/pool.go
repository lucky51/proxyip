package cmd

import (
	"context"
	"fmt"
	"github.com/lucky51/proxyip/crawler"
	"github.com/lucky51/proxyip/internal"
	"github.com/spf13/cobra"
	"time"
)
// serve http service
var serve bool
var port int

var poolCmd = &cobra.Command{
	Use: "pool",
	Short:"create proxy ip pool",
	Run: func(cmd *cobra.Command, args []string) {
		var forever = make(chan bool)
		c:=crawler.NewJXLCrawler()
		pool:=internal.NewProxyPool(&internal.PollingGetProxyIPStrategy{})
		ctx,_:=context.WithCancel(context.Background())
		go pool.StartChecker(ctx,time.Second *120)
		go pool.StartCollector(ctx,c,time.Second *5)
		if serve{
			fmt.Println("starting create http router")
			router:=internal.CreateRouter(pool)
			fmt.Println("starting listen port:",port)
			router.Run(fmt.Sprintf(":%d",port))
		}
		<-forever
		fmt.Println("blocking this,ctrl c exit.")
		select{}
	},
}

func init() {
	poolCmd.Flags().BoolVarP(&serve,"serve","s",true,"serve http")
	poolCmd.Flags().IntVarP(&port,"port","p",8081,"port")
}