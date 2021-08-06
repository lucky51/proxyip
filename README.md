# HTTP IP proxy pool

http ip proxy pool
## Usage
```shell
//install
go install github.com/lucky51/proxyip

// start a proxy pool and serve http
proxyip pool --serve -p=8081

// you can clone this project from github
git clone github.com/lucky51/proxyip
cd proxyip
go run main.go pool --serve -p=8081

//other command 
// run a crawler without proxy pool,page is optional ,defalut 1
proxyip crawl -p=1 
```

```shell
//get all available proxies
curl http://localhost:8081

//get a proxy ip by polling strategy
curl http://localhost:8081/ip
```

![home page](https://github.com/lucky51/proxyip/blob/main/screenshots/home.jpg)

![get all proxies](https://github.com/lucky51/proxyip/blob/main/screenshots/getall.jpg)

## Related Projects

* [Gin Web Framework](https://github.com/gin-gonic/gin) 
* [ASCII Table Writer](https://github.com/olekukonko/tablewriter)