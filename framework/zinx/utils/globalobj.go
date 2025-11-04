package utils

import (
	"encoding/json"
	"io"
	"os"
	"wheel/framework/zinx/ziface"
)

type GlobalObj struct {
	TcpServer ziface.IServer
	Host      string
	TcpPort   int
	Name      string

	Version          string
	MaxPacketSize    uint32
	MaxConn          int
	WorkerPoolSize   uint32
	MaxWorkerTaskLen uint32
}

func (g *GlobalObj) Reload() {
	f, err := os.Open("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	if err = json.Unmarshal(data, &GlobalObject); err != nil {
		panic(err)
	}
}

var GlobalObject *GlobalObj

func init() {
	GlobalObject = &GlobalObj{
		Name:             "ZinxServerApp",
		Version:          "V0.4",
		TcpPort:          8888,
		Host:             "0.0.0.0",
		MaxPacketSize:    4096,
		MaxConn:          12000,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
	}
	GlobalObject.Reload()
}
