package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

var (
	users = map[string]struct{}{
		"user1@example.com": {},
		"user2@example.com": {},
		"user3@example.com": {},
	}
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
	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	var (
		from, to  string
		dataStep  bool
		dataLines []string
	)

	// 发送邮件欢迎信息
	send(w, "220 server@example.com Simple Mail Transfer Service Ready")

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return
		}
		line = strings.TrimRight(line, "\r\n")
		items := strings.Fields(line)
		if len(items) <= 0 {
			send(w, "500 Syntax error, command unrecognized")
			continue
		}

		if dataStep {
			if line != "." {
				dataLines = append(dataLines, line)
			} else {
				showDetail(from, to, dataLines)
				send(w, "250 OK")
				from, to = "", ""
				dataStep = false
				dataLines = nil
			}
			continue
		}

		switch strings.ToUpper(items[0]) {
		case "HELO", "EHLO":
			send(w, "250 server@example.com")
		case "MAIL":
			if len(items) < 2 {
				send(w, "501 Syntax error in parameters or arguments")
				continue
			}
			if !strings.HasPrefix(items[1], "FROM:") {
				send(w, "501 Syntax error in parameters or arguments")
				continue
			}
			from = strings.TrimLeft(items[1], "FROM:")
			if !strings.HasPrefix(from, "<") || !strings.HasSuffix(from, ">") {
				send(w, "501 Syntax error in parameters or arguments")
				continue
			}
			from = from[1 : len(from)-1]
			if from == "" {
				send(w, "501 Syntax error in parameters or arguments")
				continue
			}
			send(w, "250 OK")
		case "RCPT":
			if len(items) < 2 {
				send(w, "501 Syntax error in parameters or arguments")
				continue
			}
			if !strings.HasPrefix(items[1], "TO:") {
				send(w, "501 Syntax error in parameters or arguments")
				continue
			}
			to = strings.TrimLeft(items[1], "TO:")
			if !strings.HasPrefix(to, "<") || !strings.HasSuffix(to, ">") {
				send(w, "501 Syntax error in parameters or arguments")
				continue
			}
			to = to[1 : len(to)-1]
			if to == "" {
				send(w, "501 Syntax error in parameters or arguments")
				continue
			}
			if _, ok := users[to]; !ok {
				send(w, "550 Requested action not taken: No such user here")
				continue
			}
			send(w, "250 OK")
		case "DATA":
			if from == "" || to == "" {
				send(w, "503 Bad sequence of commands")
				continue
			}
			dataStep = true
			send(w, "354 Start mail input; end with <CRLF>.<CRLF>")
		case "QUIT":
			send(w, "221 service@example.com Service closing transmission channel")
		default:
			send(w, "502 Command not implemented")
		}
	}
}

func send(w *bufio.Writer, msg string) {
	msg += "\r\n"
	_, _ = w.Write([]byte(msg))
	_ = w.Flush()
}

func showDetail(from, to string, dataLines []string) {
	fmt.Println("--------------------------------")
	fmt.Println("Received mail from:", from)
	fmt.Println("To:", to)
	fmt.Println("Data:")
	for _, l := range dataLines {
		fmt.Println(l)
	}
	fmt.Println("--------------------------------")
}
