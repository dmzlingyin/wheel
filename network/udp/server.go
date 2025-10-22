package main

import (
	"fmt"
	"net"
)

func main() {
	ln, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	fmt.Println("udp 监听地址: 0.0.0.0:3000")

	for {
		var buf [1024]byte
		n, addr, err := ln.ReadFromUDP(buf[:])
		if err != nil {
			fmt.Println("读取 udp 数据错误: ", err)
			continue
		}
		fmt.Printf("收到来自 %v 的 %d 字节数据: %s", addr, n, string(buf[:n]))
	}
}
