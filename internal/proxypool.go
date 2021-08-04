package internal

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

var PoolEmptyError = errors.New("proxy pool is empty")

type Crawler interface {
	GetProxies() <- chan *HttpProxyIP
	Crawl() error
	Close()
}


type GetProxyIPStrategy interface {
	Get([]*HttpProxyIP) (*HttpProxyIP,error)
}

type DefaultGetProxyIPStrategy struct {}

func (s *DefaultGetProxyIPStrategy) Get(proxies []*HttpProxyIP) ( *HttpProxyIP,error) {
	if proxies==nil ||len(proxies)==0{
		return nil, PoolEmptyError
	}
	return proxies[0],nil
}


type PollingGetProxyIPStrategy struct {
	counter int
}
// Get get proxy IP by polling strategy
func (s *PollingGetProxyIPStrategy) Get(proxies []*HttpProxyIP)(*HttpProxyIP,error)  {
	if proxies==nil ||len(proxies)==0{
		return nil, PoolEmptyError
	}
	if len(proxies)==s.counter{
		s.counter =0
	}
	defer func() {
		s.counter=s.counter+1
	}()
	return proxies[s.counter],nil
}

type proxyPool struct {
	_lock sync.Mutex
	_available int
	proxies []*HttpProxyIP
	s GetProxyIPStrategy
}

func (pool* proxyPool) StartCollector(ctx context.Context,c Crawler,delay time.Duration)  {
	fmt.Println("starting proxy ip collector...")
	if delay==0{
		delay = time.Second *5
	}
	var forever = make(chan bool)
	go func() {
		defer c.Close()
		crawLoop:
		for {
			select {
				case t:=<-time.After(delay):
					fmt.Println("collector:started by:",t.Format(time.Kitchen))
					if len(pool.proxies) < pool._available{
						c.Crawl()
					}else{
						fmt.Printf("ignore crawl, pool available ip %d < (%d) \n",len(pool.proxies),pool._available)
					}
					case <-ctx.Done():
						fmt.Println("notify exit collector.")
					break crawLoop
			}
		}
		forever<-true
	}()
	go func() {
		for ip := range c.GetProxies() {
			existsIP:=false
			for i := 0; i < len(pool.proxies); i++ {
				if pool.proxies[i].IP ==ip.IP {
					existsIP =true
					break
				}
			}
			if !existsIP{
				pool.Set(ip)
			}
		}
	}()

	<-forever
}
// StartChecker start a proxy ip checker ,this started by another go routine  in most cases
func (pool *proxyPool) StartChecker(ctx context.Context,delay time.Duration)  {
	fmt.Println("starting proxy ip checker...")
	for  {
		if delay ==0{
			delay =time.Second *20
		}
		select {
			case t:= <-time.After(delay):
				pool._lock.Lock()
				for _,proxy := range pool.proxies {
					fmt.Println("starting check proxy ip:",proxy.IP)

					result,err:=CheckProxyIP(proxy.HttpProtocol,proxy.IP,proxy.Port,"http://icanhazip.com")
					proxy.LastCheckedTime =t
					if err!=nil{
						proxy.LastCheckedState = ProxyIPStatusError
						fmt.Printf("checked result:%v \n",err)
					}else{
						proxy.LastCheckedState =ProxyIPStatusOk
						fmt.Printf("checked result:%s \n",result)
					}
				}
				// TODO: 清除无效的IP
				// var availablePIP
				for _,p := range pool.proxies {
					if p.LastCheckedState == ProxyIPStatusOk{

					}
				}
				pool._lock.Unlock()
				RenderTable(os.Stdout,fmt.Sprintf("proxy ip checker,started by:%s",time.Now().Format(time.Kitchen)),pool.proxies)
			case <-ctx.Done():
				fmt.Println("proxy ip checker is canceled")
				return
		}

	}
}
func (pool *proxyPool) Get() (*HttpProxyIP,error) {
	pool._lock.Lock()
	defer pool._lock.Unlock()
	return pool.s.Get(pool.proxies)
}

func (pool* proxyPool) Set(proxy *HttpProxyIP)  {
	pool._lock.Lock()
	defer pool._lock.Unlock()
	pool.proxies =append(pool.proxies,proxy)
}
func NewProxyPool(strategy GetProxyIPStrategy) *proxyPool {
	pool:=&proxyPool{
		s: strategy,
		_available: 20,
		proxies: make([]*HttpProxyIP,0),
	}
	return pool
}