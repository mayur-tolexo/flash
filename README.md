
[![Godocs](https://img.shields.io/badge/golang-documentation-blue.svg)](https://www.godoc.org/github.com/mayur-tolexo/flash)
[![Go Report Card](https://goreportcard.com/badge/github.com/mayur-tolexo/flash)](https://goreportcard.com/report/github.com/mayur-tolexo/flash)
[![Open Source Helpers](https://www.codetriage.com/mayur-tolexo/flash/badges/users.svg)](https://www.codetriage.com/mayur-tolexo/flash)
[![Release](https://img.shields.io/github/release/mayur-tolexo/flash.svg?style=flat-square)](https://github.com/mayur-tolexo/flash/releases)


# flash
Flash Restful API framework in Golang.
Wrapper on [gin-gonic](https://github.com/gin-gonic) engine to structure the api in better way and easy to maintain.


## Contents

- [Installation](#installation)
- [Example](#example)
- [Configuration](#configuration)
- [API Examples](#api-examples)
  - [Using GET,POST,PUT,PATCH,DELETE and OPTIONS](#using-get-post-put-patch-delete-and-options)
  - [Rest API definations are same as gin-gonic](https://github.com/gin-gonic/gin#api-examples)
  - [Using Configuration](#using-configuration)

## Installation

To install flash package, you need to install Go and set your Go workspace first.
As flash works on [gin-gonic](https://github.com/gin-gonic/gin), So all gin-gonic [Prerequisite](https://github.com/gin-gonic/gin#prerequisite) are required.

1. Download and install it:

```sh
$ go get -u github.com/mayur-tolexo/flash
```

2. Import it in your code:

```go
import "github.com/mayur-tolexo/flash"
```


### Example

```sh
# create a main file
$ cat main.go
```
```
import "github.com/mayur-tolexo/flash"

//PingService service struct containing only one api
type PingService struct {
	flash.Server `v:"1" root:"/test/"`
	ping         flash.GET `url:"/ping"`
}

//Ping api defination
func (*PingService) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func main() {
	router := flash.Default()
	router.AddService(&PingService{})
	router.Start(":8080")
}
```
```sh
# now run the main service
$ go run main.go
```
now open http://localhost:8080/v1/test/ping

### Configuration

| Tag          | Usage            
| ----------   |-----------------
| prefix, pre  | Url prefix as in: http://abc.com/[prefix]/v1/root/url                 
| root         | Url root as in: http://abc.com/prefix/v1/[root]/url                                  
| url          | Url path as in in: http://abc.com/prefix/v1/root/[url]                            
| version, v   | Url version as in: http://abc.com/prefix/v[1]/root/url
---

## API Examples

### Using GET, POST, PUT, PATCH, DELETE and OPTIONS

```go
//TestService service struct
type TestService struct {
	flash.Server    `v:"1" root:"/test/"`
	getAPI          flash.GET `url:"/"`
	postAPI         flash.POST `url:"/"`
	putAPI          flash.PUT `url:"/"`
	patchAPI        flash.PATCH `url:"/"`
	putAPI          flash.PUT `url:"/"`
	deleteAPI       flash.DELETE `url:"/"`
	optionsAPI      flash.OPTIONS `url:"/"`
}

func main() {
	router := flash.Default()
	router.AddService(&TestService{})
	router.Start(":8080")
}
```

### Using Configuration

```
import "github.com/mayur-tolexo/flash"

//PingService service struct containing only one api
type PingService struct {
	flash.Server `v:"1" root:"test" pre:"abc/"`
	ping         flash.GET `url:"/ping"`
	ping2         flash.GET `url:"/ping" v:"2"`
}

//Ping api defination
func (*PingService) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

//Ping2 api defination
func (*PingService) Ping2(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong from version 2",
	})
}

func main() {
	router := flash.Default()
	router.AddService(&PingService{})
	router.Start(":8080")
}
```
First API url http://localhost:8080/abc/v1/test/ping  
Second API url http://localhost:8080/abc/v2/test/ping
