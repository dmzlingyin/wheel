package ziface

import "net"

type IConnection interface {
	Start()
	Stop()
	GetTCPConn() *net.TCPConn
	GetRemoteAddr() net.Addr
	GetConnID() uint32
}

type HandleFunc func(*net.TCPConn, []byte, int) error
