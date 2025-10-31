package test

import (
	"net"
	"testing"
	"time"
	"wheel/framework/zinx/znet"
)

func TestClient(t *testing.T) {
	t.Log("开始连接服务器...")
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("开始打包数据...")
	// 构造数据包
	dp := znet.NewDataPack()
	data, err := dp.Pack(znet.NewMsgPackage(0, []byte("hello world")))
	if err != nil {
		t.Fatal(err)
	}

	t.Log("开始发送数据...")
	if _, err = conn.Write(data); err != nil {
		t.Fatal(err)
	}

	time.Sleep(5 * time.Second)
}
