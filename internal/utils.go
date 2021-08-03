package internal

import (
	"github.com/olekukonko/tablewriter"
	"io"
	"os"
)

// RenderTable render proxy ip table
func RenderTable(writer io.Writer,caption string,data []* HttpProxyIP)  {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"IP", "Port", "Protocol","Location","ISP","ResponseSpeed","TTL"})
	table.SetCaption(true, caption)
	for _, proxyIP := range data {
		table.Append(proxyIP.ToTableRow())
	}
	table.Render()
}
