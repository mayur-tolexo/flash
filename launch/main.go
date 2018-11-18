package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mayur-tolexo/flash"
)

type CoreService struct {
	ping flash.GET `url:"/ping"`
}

func (me *CoreService) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong bro",
	})
}

func main() {
	engine := flash.Default()
	engine.AddService(&CoreService{})
	engine.Start(":7071")
}
