package internal

import (
	"errors"
	"sync"
)

var PoolEmptyError = errors.New("proxy pool is empty")

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
	proxies []*HttpProxyIP
	s GetProxyIPStrategy
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
		proxies: make([]*HttpProxyIP,0),
	}
	return pool
}