# HTTP IP proxy pool


```shell
// start a proxy pool and serve http
go run main.go pool --serve -p=8081
```


```shell
//get all available proxies
curl http://localhost:8081

//get a proxy ip by polling strategy
curl http://localhost:8081/ip
```

![get all proxies](https://github.com/lucky51/proxyip/blob/main/screenshots/getall.jpg)
