package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
)

func encodeDomainName(domain string) []byte {
	var buf []byte
	for _, item := range strings.Split(domain, ".") {
		buf = append(buf, byte(len(item)))
		buf = append(buf, item...)
	}
	buf = append(buf, 0)
	return buf
}

func buildQuery(domain string) []byte {
	// 写入 header
	var buf bytes.Buffer
	_ = binary.Write(&buf, binary.BigEndian, uint16(0x1234))
	_ = binary.Write(&buf, binary.BigEndian, uint16(0x0100)) // RD = 1 递归查询
	_ = binary.Write(&buf, binary.BigEndian, uint16(1))      // 单个查询
	_ = binary.Write(&buf, binary.BigEndian, uint16(0))
	_ = binary.Write(&buf, binary.BigEndian, uint16(0))
	_ = binary.Write(&buf, binary.BigEndian, uint16(0))
	// 写入 question
	_ = binary.Write(&buf, binary.BigEndian, encodeDomainName(domain))
	_ = binary.Write(&buf, binary.BigEndian, uint16(1))
	_ = binary.Write(&buf, binary.BigEndian, uint16(1))
	return buf.Bytes()
}

func parseResponse(msg []byte) []net.IP {
	// 针对我们发出的请求, 只有一个 question, 但结果可能包含多个 answer
	// 首先定位到第一个 answer: offset 12 + question 的长度
	_, offset := decodeDomainName(msg, 12)
	// 跳过 QType & QClass, 此时, offset 指向第一个 answer
	offset += 4
	// 获取 answer 数量
	anCount := binary.BigEndian.Uint16(msg[6:8])

	var ips []net.IP

	// 循环遍历每一个 answer
	for i := 0; i < int(anCount); i++ {
		_, offset = decodeDomainName(msg, offset)
		if offset+10 > len(msg) {
			panic("invalid dns response")
		}
		rType := binary.BigEndian.Uint16(msg[offset : offset+2])
		rClass := binary.BigEndian.Uint16(msg[offset+2 : offset+4])
		//ttl := binary.BigEndian.Uint32(msg[offset+4 : offset+8])
		rdLength := binary.BigEndian.Uint16(msg[offset+8 : offset+10])
		offset += 10

		if offset+int(rdLength) > len(msg) {
			panic("invalid dns response")
		}

		rData := msg[offset : offset+int(rdLength)]
		// 指向下一个 answer
		offset += int(rdLength)

		if rType == 1 && rClass == 1 && rdLength == 4 {
			ips = append(ips, net.IPv4(rData[0], rData[1], rData[2], rData[3]))
		}
	}

	return ips
}

// decodeDomainName 解析域名 返回 域名和偏移量(name字段后的第一个位置)
func decodeDomainName(msg []byte, offset int) (string, int) {
	var labels []string
	pos := offset
	for {
		if pos >= len(msg) {
			panic("invalid dns response")
		}
		length := int(msg[pos])
		// 达到 0 表示结束
		if length == 0 {
			pos++
			break
		}

		// 指针压缩以 11 双位开头
		if length&0xc0 == 0xc0 {
			if pos+2 > len(msg) {
				panic("invalid dns response")
			}
			newPos := int(binary.BigEndian.Uint16(msg[pos:pos+2]) & 0x3fff)
			label, _ := decodeDomainName(msg, newPos)
			labels = append(labels, label)
			pos += 2
			break
		} else {
			// 指向第一个有效的 label
			pos++
			if pos+length > len(msg) {
				panic("invalid dns response")
			}
			labels = append(labels, string(msg[pos:pos+length]))
			pos += length
		}
	}
	return strings.Join(labels, "."), pos
}

func main() {
	// 构造DNS请求报文
	query := buildQuery("www.baidu.com")

	// 指定DNS服务器的IP和端口
	socket, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(192, 168, 4, 1),
		Port: 53,
	})
	if err != nil {
		panic(err)
	}
	defer socket.Close()

	// 发送DNS请求报文
	_, err = socket.Write(query)
	if err != nil {
		panic(err)
	}

	// 读取响应数据(最多 512 字节)
	buf := make([]byte, 512)
	n, err := socket.Read(buf)
	if err != nil {
		panic(err)
	}
	buf = buf[:n]

	// 解析DNS响应报文
	for _, ip := range parseResponse(buf) {
		fmt.Println(ip)
	}
}
