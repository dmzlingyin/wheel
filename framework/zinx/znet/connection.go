package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"wheel/framework/zinx/zconf"
	"wheel/framework/zinx/ziface"
)

type Connection struct {
	TcpServer    ziface.IServer
	Conn         *net.TCPConn
	ConnID       uint32
	isClosed     bool
	ExitBuffChan chan bool
	MsgHandler   ziface.IMsgHandler
	msgChan      chan []byte // 无缓冲
	msgBuffChan  chan []byte // 有缓冲
	property     map[string]any
	lock         sync.RWMutex
}

func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) *Connection {
	c := &Connection{
		TcpServer:    server,
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		ExitBuffChan: make(chan bool, 1),
		MsgHandler:   msgHandler,
		msgChan:      make(chan []byte),
		msgBuffChan:  make(chan []byte, zconf.GlobalObject.MaxMsgChanLen),
		property:     make(map[string]any),
	}
	c.TcpServer.GetConnMgr().Add(c)
	return c
}

func (c *Connection) Start() {
	go c.startReader()
	go c.startWriter()

	c.TcpServer.CallOnConnStart(c)

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
		if zconf.GlobalObject.WorkerPoolSize > 0 {
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			go c.MsgHandler.DoMsgHandler(&req)
		}
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
		case data, ok := <-c.msgBuffChan:
			if ok {
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("write msg error:", err)
					c.ExitBuffChan <- true
					return
				}
			} else {
				break
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
	c.TcpServer.CallOnConnStop(c)
	c.Conn.Close()
	c.ExitBuffChan <- true
	c.TcpServer.GetConnMgr().Remove(c)
	close(c.ExitBuffChan)
	close(c.msgBuffChan)
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

func (c *Connection) SendBuffMsg(msgID uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection closed")
	}

	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgID, data))
	if err != nil {
		return err
	}

	c.msgBuffChan <- msg

	return nil
}

func (c *Connection) SetProperty(key string, value any) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.property[key] = value
}

func (c *Connection) GetProperty(key string) (any, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if value, ok := c.property[key]; ok {
		return value, nil
	}
	return nil, errors.New("property not found")
}

func (c *Connection) RemoveProperty(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.property, key)
}
