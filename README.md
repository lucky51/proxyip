# HTTP IP proxy pool


```shell
//install
go install github.com/lucky51/proxyip
proxyip pool --serve -p=8081
// start a proxy pool and serve http
// go run main.go pool --serve -p=8081
```


```shell
//get all available proxies
curl http://localhost:8081

//get a proxy ip by polling strategy
curl http://localhost:8081/ip
```

![get all proxies](https://github.com/lucky51/proxyip/blob/main/screenshots/getall.jpg)
