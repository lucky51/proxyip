# HTTP IP proxy Pool


```shell
// start a proxy pool and serve http
go run main.go pool --serve -p=8081
```


```shell
//get all available proxy ip
curl http://localhost:8081

//get a proxy ip by polling strategy
curl http://localhost:8081/ip
```

