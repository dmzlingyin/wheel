package znet

import (
	"errors"
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
	MsgHandler   ziface.IMsgHandler
	msgChan      chan []byte
}

func NewConnection(conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) *Connection {
	return &Connection{
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		ExitBuffChan: make(chan bool, 1),
		MsgHandler:   msgHandler,
		msgChan:      make(chan []byte),
	}
}

func (c *Connection) Start() {
	go c.startReader()
	go c.startWriter()

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
		// 能否把 c 传进去, 在 unpack 方法里面读取数据呢?
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error:", err)
			c.ExitBuffChan <- true
			return
		}

		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err = io.ReadFull(c.GetTCPConn(), data); err != nil {
				fmt.Println("read msg data error:", err)
				c.ExitBuffChan <- true
				return
			}
		}
		msg.SetData(data)

		req := Request{conn: c, msg: msg}
		go c.MsgHandler.DoMsgHandler(&req)
	}
}

func (c *Connection) startWriter() {
	fmt.Println("start writer ...")
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("write msg error:", err)
				c.ExitBuffChan <- true
				return
			}
		case <-c.ExitBuffChan:
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

func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection closed")
	}

	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgID, data))
	if err != nil {
		return err
	}

	c.msgChan <- msg

	return nil
}
