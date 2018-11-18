package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mayur-tolexo/flash"
)

type TestService struct {
	ping  flash.GET `url:"/ping"`
	ping2 flash.GET `url:"/ping2"`
}

func (*TestService) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func main() {
	engine := flash.Default()
	engine.AddService(&TestService{})
	engine.Start(":7071")
}
