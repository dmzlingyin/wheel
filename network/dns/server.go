package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

var cache = make(map[string][]net.IP)

// 假设客户端发送的请求包完整无误
func processRequest(socket *net.UDPConn, msg []byte, addr *net.UDPAddr) {
	// 解析域名
	domain, offset := decodeDomainName(msg, 12)
	ips, ok := cache[domain]
	if !ok {
		ips = query(msg)
		// 写入缓存
		cache[domain] = ips
	}
	// 构造响应包
	res := buildResponse(ips, msg, offset+4)
	_, err := socket.WriteToUDP(res, addr)
	if err != nil {
		fmt.Println("发送错误: ", err)
	}
}

func query(msg []byte) []net.IP {
	// 模拟查询
	return []net.IP{
		net.IPv4(192, 168, 1, 1),
		net.IPv4(192, 168, 1, 2),
	}
}

func buildResponse(ips []net.IP, msg []byte, offset int) []byte {
	var buf bytes.Buffer
	// header
	_ = binary.Write(&buf, binary.BigEndian, msg[:2])
	_ = binary.Write(&buf, binary.BigEndian, uint16(0x8180))
	_ = binary.Write(&buf, binary.BigEndian, uint16(1))
	_ = binary.Write(&buf, binary.BigEndian, uint16(len(ips)))
	_ = binary.Write(&buf, binary.BigEndian, uint16(0))
	_ = binary.Write(&buf, binary.BigEndian, uint16(0))
	// 写入 question
	_ = binary.Write(&buf, binary.BigEndian, msg[12:offset])
	// 写入 answers
	for _, ip := range ips {
		ip = ip.To4()
		if ip == nil {
			continue
		}
		// 指针
		_ = binary.Write(&buf, binary.BigEndian, uint16(0xc00c))
		// type
		_ = binary.Write(&buf, binary.BigEndian, uint16(1))
		// class
		_ = binary.Write(&buf, binary.BigEndian, uint16(1))
		// ttl
		_ = binary.Write(&buf, binary.BigEndian, uint32(300))
		// rdLength
		_ = binary.Write(&buf, binary.BigEndian, uint16(4))
		// rData
		_ = binary.Write(&buf, binary.BigEndian, ip)
	}
	return buf.Bytes()
}

func main() {
	socket, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	if err != nil {
		panic(err)
	}
	defer socket.Close()

	fmt.Println("udp 监听地址: 0.0.0.0:3000")

	for {
		var buf [512]byte
		n, addr, err := socket.ReadFromUDP(buf[:])
		if err != nil {
			fmt.Println("读取 udp 数据错误: ", err)
			continue
		}
		go processRequest(socket, buf[:n], addr)
	}
}
