package internal

import (
	"github.com/olekukonko/tablewriter"
	"io"
	"os"
)

// RenderTable render proxy ip table
func RenderTable(writer io.Writer,caption string,data []* HttpProxyIP)  {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"IP", "Port", "Protocol","Location","ISP","ResponseSpeed","TTL","LastCheckedTime","State"})

	table.SetCaption(true, caption)
	for _, proxyIP := range data {
		table.Append(proxyIP.ToTableRow())
	}
	table.Render()
}


func DelItem(vs []*HttpProxyIP, s *HttpProxyIP) []*HttpProxyIP{
	for i := 0; i < len(vs); i++ {
		if s == vs[i] {
			vs = append(vs[:i], vs[i+1:]...)
			i = i-1
		}
	}
	return vs
}

