package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
)

const (
	HttpStatusBadRequest = "HTTP/1.1 400 Bad Request\r\nContent-Length: 0\r\n\r\n"
)

const (
	HttpMethodGet  = "GET"
	HttpMethodPost = "POST"
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
	defer conn.Close()
	reader := bufio.NewReader(conn)

	// 读取请求行
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		_, _ = conn.Write([]byte(HttpStatusBadRequest))
		return
	}
	fmt.Printf("请求行: %s", requestLine)

	items := strings.Split(requestLine, " ")
	if len(items) != 3 {
		_, _ = conn.Write([]byte(HttpStatusBadRequest))
		return
	}

	switch items[0] {
	case HttpMethodGet:
		for {
			requestHeader, err := reader.ReadString('\n')
			if err != nil {
				return
			}
			if requestHeader == "\r\n" {
				// 请求头结束
				body := "Hello World\n"
				response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(body), body)
				_, _ = conn.Write([]byte(response))
				return
			}
			fmt.Printf("请求头: %s", requestHeader)
		}
	case HttpMethodPost:
		headers := make(map[string]string)
		for {
			requestHeader, err := reader.ReadString('\n')
			if err != nil {
				return
			}
			if requestHeader == "\r\n" {
				// 请求头结束
				break
			}
			parts := strings.SplitN(requestHeader, ": ", 2)
			if len(parts) == 2 {
				headers[strings.ToLower(parts[0])] = strings.TrimSpace(parts[1])
			}
		}

		lengthStr, ok := headers["content-length"]
		if !ok {
			_, _ = conn.Write([]byte(HttpStatusBadRequest))
			return
		}
		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			_, _ = conn.Write([]byte(HttpStatusBadRequest))
			return
		}
		if length < 0 {
			_, _ = conn.Write([]byte(HttpStatusBadRequest))
			return
		}
		// 假设只接收 application/json 类型
		if headers["content-type"] != "application/json" {
			_, _ = conn.Write([]byte(HttpStatusBadRequest))
			return
		}

		// 读取请求体
		body := make([]byte, length)
		_, err = reader.Read(body)
		if err != nil {
			_, _ = conn.Write([]byte(HttpStatusBadRequest))
			return
		}
		type Data struct {
			Option string `json:"option"`
		}
		var data Data
		if err = json.Unmarshal(body, &data); err != nil {
			_, _ = conn.Write([]byte(HttpStatusBadRequest))
			return
		}
		response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/json\r\nContent-Length: %d\r\n\r\n%s", length, body)
		if data.Option == "ping" {
			response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/json\r\nContent-Length: 17\r\n\r\n{\"option\":\"pong\"}")
		}
		_, _ = conn.Write([]byte(response))
	default:
		_, _ = conn.Write([]byte(HttpStatusBadRequest))
	}
}
