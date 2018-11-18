package main

import "github.com/mayur-tolexo/flash"

type CoreService struct {
	ping flash.GET `url:"/ping"`
}

func (me *CoreService) Ping() string {
	return "pong"
}

func main() {
	engine := flash.Default()
	engine.AddService(&CoreService{})
	engine.Start(":7071")
}
