package znet

import (
	"fmt"
	"net"
	"wheel/framework/zinx/utils"
	"wheel/framework/zinx/ziface"
)

type Server struct {
	Name      string
	IPVersion string
	Addr      string
	Port      int
	Router    ziface.IRouter
}

func NewServer() ziface.IServer {
	return &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		Addr:      utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
	}
}

func (s *Server) Start() {
	fmt.Printf("server name[%s] ip[%s] port[%d]\n", s.Name, s.IPVersion, s.Port)
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d,  MaxPacketSize: %d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPacketSize)
	go func() {
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
			dealConn := NewConnection(conn, cid, s.Router)
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("server stop")
}

func (s *Server) Serve() {
	// 此处可能会有其他初始化操作
	s.Start()

	select {}
}

func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
}
