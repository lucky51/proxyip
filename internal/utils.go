package internal

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"io"
	"strconv"
)

// RenderProxyIPTable render proxy ip table
func RenderProxyIPTable(writer io.Writer,caption string,data []* HttpProxyIP)  {
	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{"IP", "Port", "Anonymity","Protocol","Location","ISP","ResponseSpeed","TTL","LastCheckedTime","State","Metadata"})

	table.SetCaption(true, caption)
	for _, proxyIP := range data {
		proxy:= *proxyIP
		table.Append(proxy.ToTableRow())
	}
	table.Render()
}
func RenderConfigTable(writer io.Writer,caption string,config ProxyPoolConfig)  {
	table :=tablewriter.NewWriter(writer)
	table.SetHeader([]string{"Alarm value","Delay of checker","Delay of collector","Get IP strategy"})
	table.Append([]string{
		strconv.Itoa(config.AlarmValue),
		fmt.Sprintf("%d s",int(config.DelayOfChecker.Seconds())),
		fmt.Sprintf("%d s",int(config.DelayOfCollector.Seconds())),
		config.GetProxyStrategy.Name(),
	})
	table.Render()
}