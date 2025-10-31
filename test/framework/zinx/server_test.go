package test

import (
	"fmt"
	"testing"
	"wheel/framework/zinx/ziface"
	"wheel/framework/zinx/znet"
)

type CustomRouter struct {
	znet.BaseRouter
}

func (cr *CustomRouter) Handle(request ziface.IRequest) {
	fmt.Println("handle data: ", string(request.GetData()))
}

func TestServer(t *testing.T) {
	s := znet.NewServer()
	s.AddRouter(0, &CustomRouter{})
	s.Serve()
}
