package internal

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type ProxyIPStatus uint8

const (
	ProxyIPStatusOk    = 200
	ProxyIPStatusError = 201
)

const DefaultIPCheckUrl  = "http://icanhazip.com"

// HttpProxyIP proxy ip item
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
	LastCheckedTime  time.Time
	LastCheckedState ProxyIPStatus
}

// ToTableRow map to a table row
func (ip *HttpProxyIP) ToTableRow() []string {
	state := ""
	lastCheckedTime := ip.LastCheckedTime.Format("15:04:05.000")
	if ip.LastCheckedState == ProxyIPStatusOk {
		state = "ok"
	} else if ip.LastCheckedState == 00 {
		state = "--"
		lastCheckedTime = "--"
	}else{
		state="error"
	}
	return []string{
		ip.IP,
		strconv.Itoa(ip.Port),
		ip.Anonymity,
		ip.HttpProtocol,
		ip.Location,
		ip.ISP,
		ip.ResponseSpeed,
		ip.TTL,
		lastCheckedTime,
		state,
		fmt.Sprintf("t=%s,p=%s",ip.Metadata["totals"],ip.Metadata["page"]),
	}
}
func CheckProxyUrl(proxyUrlStr,validUrl string) (string,int,error) {
	proxyUrl, err := url.Parse(proxyUrlStr)
	if err != nil {
		return "",0, err
	}
	newTrans := &http.Transport{
		Proxy:                 http.ProxyURL(proxyUrl),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		MaxIdleConnsPerHost: 50,
		ResponseHeaderTimeout: 10 * time.Second,
	}
	client := http.Client{
		Transport: newTrans,
		Timeout:   time.Second*10,
	}
	request, _ := http.NewRequest("GET", validUrl, nil)
	request.Header.Add("accept", "text/plain")
	resp, err := client.Do(request)
	var status =0
	if resp!=nil{
		status =resp.StatusCode
	}
	if err != nil {
		return "", status,err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", status,err
	}
	return string(body), status,nil
}
// CheckProxyIP check proxy ip status
func CheckProxyIP(protocol string, ip string, port int, validUrl string) (string,int, error) {
	if port == 0 {
		port = 80
	}
	if protocol == "" {
		protocol = "http"
	}
	proxyUrlStr := fmt.Sprintf("%s://%s:%d", protocol, ip, port)
	return CheckProxyUrl(proxyUrlStr,validUrl)
}
