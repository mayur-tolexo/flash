package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mayur-tolexo/flash"
)

type TestService struct {
	flash.Server `version:"1" root:"/test/"`
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
	engine.Start(":7071")
}
