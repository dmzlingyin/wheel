package test

import (
	"fmt"
	"testing"
	"wheel/framework/zinx/ziface"
	"wheel/framework/zinx/znet"
)

type CustomRouter struct {
}

func (cr *CustomRouter) PreHandle(request ziface.IRequest) {
	_, _ = request.GetConnection().GetTCPConn().Write([]byte("preHandle\n"))
	fmt.Println("preHandle data: ", string(request.GetData()))
}

func (cr *CustomRouter) Handle(request ziface.IRequest) {
	_, _ = request.GetConnection().GetTCPConn().Write([]byte("handle\n"))
	fmt.Println("handle data: ", string(request.GetData()))
}

func (cr *CustomRouter) PostHandle(request ziface.IRequest) {
	_, _ = request.GetConnection().GetTCPConn().Write([]byte("postHandle\n"))
	fmt.Println("postHandle data: ", string(request.GetData()))
}

func TestServer(t *testing.T) {
	s := znet.NewServer("zinx v0.3")
	s.AddRouter(&CustomRouter{})
	s.Serve()
}
