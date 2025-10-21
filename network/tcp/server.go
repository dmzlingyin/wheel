package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":3000")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	fmt.Println("监听地址: 0.0.0.0:3000")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("获取连接错误:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	fmt.Println("处理连接:", conn.RemoteAddr())
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("连接关闭:", conn.RemoteAddr())
			} else {
				fmt.Println("读取出错:", err)
			}
			return
		}
		fmt.Printf("收到数据: %s", line)
	}
}
