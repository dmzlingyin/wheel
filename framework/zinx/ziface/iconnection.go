package ziface

import "net"

type IConnection interface {
	Start()
	Stop()
	GetTCPConn() *net.TCPConn
	GetRemoteAddr() net.Addr
	GetConnID() uint32
	SendMsg(uint32, []byte) error
	SendBuffMsg(uint32, []byte) error
}
