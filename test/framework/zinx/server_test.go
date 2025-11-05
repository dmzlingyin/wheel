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

func startConn(conn ziface.IConnection) {
	fmt.Println("startConn")
}

func stopConn(conn ziface.IConnection) {
	fmt.Println("stopConn")
}

func TestServer(t *testing.T) {
	s := znet.NewServer()
	s.AddRouter(0, &CustomRouter{})
	s.SetOnConnStart(startConn)
	s.SetOnConnStop(stopConn)
	s.Serve()
}
