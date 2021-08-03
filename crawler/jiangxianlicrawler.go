package crawler

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/lucky51/proxyip/internal"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
)

const CrawlUrlBase = "https://ip.jiangxianli.com/"
const TotalIPCountRegStr = `laypage\.render\(\{ elem: 'paginate'\, count: \"(\d+)\"`
var totalIPCountRegExp = regexp.MustCompile(TotalIPCountRegStr)

func getCrawlUrl(page int) string {
	if page==0 {
		page =1
	}
	return fmt.Sprintf("%s?page=%d", CrawlUrlBase,page)
}

type JXLCrawler struct {
	p chan *internal.HttpProxyIP
}

func (c *JXLCrawler) GetProxies()<-chan *internal.HttpProxyIP {
	return c.p
}

func (c*JXLCrawler) Close() {
	close(c.p)
}

func (c *JXLCrawler) Crawl(page int)error  {
	if page<1{
		page =1
	}
	crawlUrl:=getCrawlUrl(page)
	resp,err:=http.Get(crawlUrl)
	if err!=nil{
		return err
	}
	defer resp.Body.Close()
	body,err:= ioutil.ReadAll(resp.Body)
	if err!=nil{
		return err
	}
	doc,err:=goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err!=nil{
		return err

	}
	//var data = make([]*internal.HttpProxyIP,0)
	trs:= doc.Find("div.ip-tables > div.layui-form > table > tbody >tr")
	//total := totalIPCountRegExp.FindSubmatch(body)
	//captionStr:=fmt.Sprintf("proxy ip :%s ,%d",total[1],page)
	trs.Each(func(i int, selection *goquery.Selection) {
		tds:=selection.Find("td")
		if tds.Length()==1{
			return
		}
		currPortStr:=tds.Eq(1).Text()
		currPort,err:=strconv.Atoi(currPortStr)
		if err!=nil{
			fmt.Println(err)
			return
		}
		proxyItem:=&internal.HttpProxyIP{
			IP: tds.Eq(0).Text(),
			Port: currPort,
			Anonymity:tds.Eq(2).Text(),
			HttpProtocol:tds.Eq(3).Text(),
			Location:tds.Eq(4).Text(),
			ISP:tds.Eq(6).Text(),
			ResponseSpeed:tds.Eq(7).Text(),
			TTL:tds.Eq(8).Text(),
			LastValidateTime: tds.Eq(9).Text(),
		}
		c.p <-proxyItem
		//data=append(data,proxyItem)
	})
	//renderTable(os.Stdout,captionStr,data)
	return nil
}


func NewJXLCrawler() *JXLCrawler {
	c:=&JXLCrawler{p: make(chan *internal.HttpProxyIP,10)}
	return c
}
