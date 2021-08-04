package internal

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)


type ProxyIPStatus  uint8

const (
	ProxyIPStatusOk = 200
	ProxyIPStatusError =201
)

type HttpProxyIP struct {
	Port             int
	IP               string
	Anonymity        string
	HttpProtocol     string
	Location         string
	ISP              string
	ResponseSpeed    string
	TTL              string
	LastValidateTime string
	Metadata         map[string]string
	LastCheckedTime 	 time.Time
	LastCheckedState ProxyIPStatus
}
func (ip * HttpProxyIP) ToTableRow() []string {
	state:=""
	if ip.LastCheckedState == ProxyIPStatusOk{
		state="ok"
	}else{
		state="error"
	}
	return []string{ip.IP,strconv.Itoa(ip.Port),ip.Anonymity,ip.Location,ip.ISP,ip.ResponseSpeed,ip.TTL,ip.LastCheckedTime.Format(time.Kitchen),state}
}
func CheckProxyIP(protocol string,ip string,port int,validUrl string) (string,error)  {
	if port==0{
		port =80
	}
	if protocol==""{
		protocol="http"
	}
	proxyUrlStr:=fmt.Sprintf("%s://%s:%d",protocol,ip,port)
	proxyUrl,err:=url.Parse(proxyUrlStr)
	if err!=nil {
		return "",err
	}
	newTrans:=&http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
		ResponseHeaderTimeout: time.Second *10,
	}
	client:=http.Client{
		Transport:newTrans,
		Timeout: time.Second*10,
	}
	resp,err:=client.Get(validUrl)
	if err!=nil{
		return "", err
	}
	defer resp.Body.Close()
	body,err:=ioutil.ReadAll(resp.Body)
	if err!=nil{
		return "", err
	}
	return string(body),nil
}
