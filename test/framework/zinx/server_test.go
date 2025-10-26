package test

import (
	"testing"
	"wheel/framework/zinx/znet"
)

func TestServer(t *testing.T) {
	s := znet.NewServer("zinx v0.1")
	s.Serve()
}
