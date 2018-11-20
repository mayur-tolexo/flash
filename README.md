
[![Godocs](https://img.shields.io/badge/golang-documentation-blue.svg)](https://www.godoc.org/github.com/mayur-tolexo/flash)
[![Go Report Card](https://goreportcard.com/badge/github.com/mayur-tolexo/flash)](https://goreportcard.com/report/github.com/mayur-tolexo/flash)
[![Open Source Helpers](https://www.codetriage.com/mayur-tolexo/flash/badges/users.svg)](https://www.codetriage.com/mayur-tolexo/flash)
[![Release](https://img.shields.io/github/release/mayur-tolexo/flash.svg?style=flat-square)](https://github.com/mayur-tolexo/flash/releases)


# flash
Flash Restful API framework in Golang.
Wrapper on [gin-gonic](https://github.com/gin-gonic) engine to struct the api in better way and easy to understand.


#### Q: What all configurations are available in flash?

| Tag          | Usage            
| ----------   |-----------------
| prefix, pre  | Url prefix as in: http://abc.com/[prefix]/v1/root/url                 
| root         | Url root as in: http://abc.com/prefix/v1/[root]/url                                  
| url          | Url path as in in: http://abc.com/prefix/v1/root/[url]                            
| version, v   | Url version as in: http://abc.com/prefix/v[1]/root/url
---

### Example
```

type TestService struct {
	flash.Server `v:"1" root:"/test/"`
	ping         flash.GET `url:"/ping"`
	ping2        flash.GET `url:"/ping" version:"2"`
}

func (*TestService) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (*TestService) Ping2(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong with version",
	})
}

func main() {
	engine := flash.Default()
	engine.AddService(&TestService{})
	engine.Start(":8080")
}
```

