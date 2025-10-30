package znet

import (
	"fmt"
	"io"
	"net"
	"wheel/framework/zinx/ziface"
)

type Connection struct {
	Conn         *net.TCPConn
	ConnID       uint32
	isClosed     bool
	ExitBuffChan chan bool
	Router       ziface.IRouter
}

func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	return &Connection{
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		ExitBuffChan: make(chan bool, 1),
		Router:       router,
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
		dp := NewDataPack()
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConn(), headData); err != nil {
			fmt.Println("read head error:", err)
			c.ExitBuffChan <- true
			return
		}
		// 拆包
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error:", err)
			c.ExitBuffChan <- true
			return
		}

		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConn(), data); err != nil {
				fmt.Println("read msg data error:", err)
				c.ExitBuffChan <- true
				return
			}
		}
		msg.SetData(data)

		req := Request{conn: c, msg: msg}
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
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
