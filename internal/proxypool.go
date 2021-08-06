package internal

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var PoolEmptyError = errors.New("proxy pool is empty")

const TimeLayout = "15:04:05.000"

type Crawler interface {
	GetProxies() <-chan *HttpProxyIP
	Crawl() error
	Close()
}
type GetProxyIPStrategy interface {
	Get([]*HttpProxyIP) (*HttpProxyIP, error)
	Name() string
}

type DefaultGetProxyIPStrategy struct{}

func (s *DefaultGetProxyIPStrategy) Name() string {
	return "Default(get first)"
}
func (s *DefaultGetProxyIPStrategy) Get(proxies []*HttpProxyIP) (*HttpProxyIP, error) {
	if proxies == nil || len(proxies) == 0 {
		return nil, PoolEmptyError
	}
	return proxies[0], nil
}

type PollingGetProxyIPStrategy struct {
	counter int
}

func (s *PollingGetProxyIPStrategy) Name() string {
	return "Polling"
}

// Get get proxy IP by polling strategy
func (s *PollingGetProxyIPStrategy) Get(proxies []*HttpProxyIP) (*HttpProxyIP, error) {
	if proxies == nil || len(proxies) == 0 {
		return nil, PoolEmptyError
	}
	j := 0
	for j < len(proxies) {
		s.counter = s.counter + 1
		if s.counter >= len(proxies) {
			s.counter = 0
		}
		if currIP := proxies[s.counter]; currIP.LastCheckedState == ProxyIPStatusOk {
			return currIP, nil
		}
		j = j + 1
	}
	return nil, PoolEmptyError
}

type ProxyPoolConfig struct {
	// collector will pull the latest Proxy IP when pool size less than this value
	AlarmValue int
	// delay of checker
	DelayOfChecker time.Duration
	// delay of collector
	DelayOfCollector time.Duration
	GetProxyStrategy GetProxyIPStrategy
}

type proxyPool struct {
	l       sync.RWMutex
	c       *ProxyPoolConfig
	proxies []*HttpProxyIP
}

// StartCollector start a collector to crawl proxy ip
func (pool *proxyPool) StartCollector(ctx context.Context, c Crawler) {
	fmt.Println("starting proxy ip collector...")
	var first = true
	delay := time.Duration(0)
	var forever = make(chan bool)
	go func() {
		defer c.Close()
	crawLoop:
		for {
			if first {
				delay = 0
			} else {
				delay = pool.c.DelayOfCollector
			}
			select {
			case t := <-time.After(delay):
				fmt.Println("collector:started by:", t.Format(TimeLayout))
				if len(pool.proxies) < pool.c.AlarmValue {
					fmt.Printf("pool available ip %d < (%d) \n", len(pool.proxies), pool.c.AlarmValue)
					c.Crawl()
				} else {
					fmt.Printf("ignore crawl, pool available ip %d >= (%d) \n", len(pool.proxies), pool.c.AlarmValue)
				}
			case <-ctx.Done():
				fmt.Println("notify exit collector.")
				break crawLoop
			}
			first = false
		}
		forever <- true
	}()
	go func() {
		for ip := range c.GetProxies() {
			currProxy := *ip
			existsIP := false
			for i := 0; i < len(pool.proxies); i++ {
				if pool.proxies[i].IP == currProxy.IP {
					existsIP = true
					break
				}
			}
			if !existsIP && len(pool.proxies) < pool.c.AlarmValue {
				pool.Set(&currProxy)
			}
		}
	}()

	<-forever
}

// StartChecker start a proxy ip checker ,this started by another go routine  in most cases
func (pool *proxyPool) StartChecker(ctx context.Context) {
	fmt.Println("starting proxy ip checker...")
	var first = true
	delay := time.Duration(0)
	for {
		if first {
			delay = time.Second * 2
		} else {
			delay = pool.c.DelayOfChecker
		}
		select {
		case t := <-time.After(delay):
			pool.l.RLock()
			for _, proxy := range pool.proxies {
				fmt.Println("starting check proxy ip:", proxy.IP)

				result, s, err := CheckProxyIP(proxy.HttpProtocol, proxy.IP, proxy.Port, "http://icanhazip.com")
				proxy.LastCheckedTime = t
				if err != nil {
					proxy.LastCheckedState = ProxyIPStatusError
					fmt.Println("checked result error,http status code:", s,"error:",err)
				} else {
					if s != http.StatusOK || len(result) > 15 {
						proxy.LastCheckedState = ProxyIPStatusError
						fmt.Printf("checked result:%s ,http status code:%d\n", "INVALID PROXY IP", s)
					} else {
						proxy.LastCheckedState = ProxyIPStatusOk
						fmt.Printf("checked result:%s \n", result)
					}
				}
				time.Sleep(time.Millisecond * 400)
			}
			// clear invalid proxy ip
			var temp = pool.proxies[:0]
			for _, proxy := range pool.proxies {
				if proxy.LastCheckedState == ProxyIPStatusOk {
					temp = append(temp, proxy)
				}
			}
			pool.l.RUnlock()
			pool.l.Lock()
			pool.proxies = temp
			pool.l.Unlock()
		case <-ctx.Done():
			fmt.Println("proxy ip checker is canceled")
			return
		}
		first = false
	}
}

func (pool *proxyPool) Get() (*HttpProxyIP, error) {
	pool.l.RLock()
	defer pool.l.RUnlock()
	return pool.c.GetProxyStrategy.Get(pool.proxies)
}

func (pool *proxyPool) GetAll() ([]*HttpProxyIP, error) {
	//pool.l.RLock()
	//defer pool.l.RUnlock()
	if pool.proxies == nil || len(pool.proxies) == 0 {
		return pool.proxies, PoolEmptyError
	}
	return pool.proxies, nil
}

func (pool *proxyPool) Set(proxy *HttpProxyIP) {
	pool.l.Lock()
	defer pool.l.Unlock()
	pool.proxies = append(pool.proxies, proxy)
}

// NewProxyPool create proxy pool with configuration
func NewProxyPool(config *ProxyPoolConfig) *proxyPool {
	if config.AlarmValue|0 == 0 {
		config.AlarmValue = 10
	}
	if config.GetProxyStrategy == nil {
		config.GetProxyStrategy = &DefaultGetProxyIPStrategy{}
	}
	if config.DelayOfChecker == 0 {
		config.DelayOfChecker = time.Second * 5
	}
	if config.DelayOfCollector == 0 {
		config.DelayOfCollector = time.Second * 5
	}
	pool := &proxyPool{
		c:       config,
		proxies: make([]*HttpProxyIP, 0),
	}
	return pool
}
