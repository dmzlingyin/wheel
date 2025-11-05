package znet

import (
	"fmt"
	"net"
	"wheel/framework/zinx/utils"
	"wheel/framework/zinx/ziface"
)

type Server struct {
	Name        string
	IPVersion   string
	Addr        string
	Port        int
	Router      ziface.IRouter
	msgHandler  ziface.IMsgHandler
	connManager ziface.IConnManager
	OnConnStart func(ziface.IConnection)
	OnConnStop  func(ziface.IConnection)
}

func NewServer() ziface.IServer {
	return &Server{
		Name:        utils.GlobalObject.Name,
		IPVersion:   "tcp4",
		Addr:        utils.GlobalObject.Host,
		Port:        utils.GlobalObject.TcpPort,
		msgHandler:  NewMsgHandle(),
		connManager: NewConnManager(),
	}
}

func (s *Server) Start() {
	fmt.Printf("server name[%s] ip[%s] port[%d]\n", s.Name, s.IPVersion, s.Port)
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d,  MaxPacketSize: %d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPacketSize)

	go func() {
		// 启动工作池
		s.msgHandler.StartWorkPool()

		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.Addr, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err:", err)
			return
		}
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen tcp err:", err)
			return
		}
		defer listener.Close()

		var cid uint32

		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("accept err:", err)
				continue
			}

			// 判断连接数量
			if s.connManager.Len() >= utils.GlobalObject.MaxConn {
				conn.Close()
				continue
			}

			dealConn := NewConnection(s, conn, cid, s.msgHandler)
			go dealConn.Start()
			cid++
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("server stop")
	s.connManager.CleanConn()
}

func (s *Server) Serve() {
	// 此处可能会有其他初始化操作
	s.Start()

	select {}
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.msgHandler.AddRouter(msgID, router)
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.connManager
}

func (s *Server) SetOnConnStart(f func(ziface.IConnection)) {
	s.OnConnStart = f
}

func (s *Server) SetOnConnStop(f func(ziface.IConnection)) {
	s.OnConnStop = f
}

func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		s.OnConnStart(conn)
	}
}

func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		s.OnConnStop(conn)
	}
}
