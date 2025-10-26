package znet

import (
	"fmt"
	"net"
	"wheel/framework/zinx/ziface"
)

type Connection struct {
	Conn         *net.TCPConn
	ConnID       uint32
	isClosed     bool
	handler      ziface.HandleFunc
	ExitBuffChan chan bool
}

func NewConnection(conn *net.TCPConn, connID uint32, handler ziface.HandleFunc) *Connection {
	return &Connection{
		Conn:         conn,
		ConnID:       connID,
		handler:      handler,
		isClosed:     false,
		ExitBuffChan: make(chan bool, 1),
	}
}

func (c *Connection) Start() {
	go c.startReader()

	for {
		select {
		case <-c.ExitBuffChan:
			return
		}
	}
}

func (c *Connection) startReader() {
	fmt.Println("start reader ...")
	defer c.Stop()

	for {
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("reader error:", err)
			c.ExitBuffChan <- true
			continue
		}
		if err = c.handler(c.Conn, buf, cnt); err != nil {
			fmt.Println("reader error:", err)
			c.ExitBuffChan <- true
			return
		}
	}
}

func (c *Connection) Stop() {
	if c.isClosed {
		return
	}
	c.isClosed = true
	c.Conn.Close()
	c.ExitBuffChan <- true
	close(c.ExitBuffChan)
}

func (c *Connection) GetTCPConn() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}
