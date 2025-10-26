package znet

import (
	"fmt"
	"net"
	"wheel/framework/zinx/ziface"
)

type Server struct {
	Name      string
	IPVersion string
	Addr      string
	Port      int
}

func NewServer(name string) ziface.IServer {
	return &Server{
		Name:      name,
		IPVersion: "tcp4",
		Addr:      "0.0.0.0",
		Port:      8888,
	}
}

func (s *Server) Start() {
	fmt.Printf("server name[%s] ip[%s] port[%d]\n", s.Name, s.IPVersion, s.Port)

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
			dealConn := NewConnection(conn, cid, Callback)
			go dealConn.Start()
		}
	}()
}

func Callback(conn *net.TCPConn, buf []byte, cnt int) error {
	if _, err := conn.Write(buf[:cnt]); err != nil {
		return err
	}
	return nil
}

func (s *Server) Stop() {
	fmt.Println("server stop")
}

func (s *Server) Serve() {
	// 此处可能会有其他初始化操作
	s.Start()

	select {}
}
