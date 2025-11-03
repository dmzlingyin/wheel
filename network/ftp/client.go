package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

func main() {
	// 建立控制连接
	controlConn, err := net.Dial("tcp", "192.168.7.251:21")
	if err != nil {
		log.Fatalln(err)
	}
	defer controlConn.Close()

	reader := bufio.NewReader(controlConn)

	// 读取欢迎信息
	if err = welcome(reader); err != nil {
		log.Fatalln(err)
	}
	// 登录
	if err = login(controlConn, reader); err != nil {
		log.Fatalln(err)
	}
	// 指定被动模式
	addr, err := setPassiveMode(controlConn, reader)
	if err != nil {
		log.Fatalln(err)
	}
	// 获取文件列表
	if err = list(controlConn, reader, addr); err != nil {
		log.Fatalln(err)
	}
}

func welcome(r *bufio.Reader) error {
	// 读取欢迎信息
	code, msg, err := read(r)
	if err != nil {
		return err
	}
	if code != 220 {
		return errors.New(fmt.Sprintf("invalid response code: %d, message: %s", code, msg))
	}
	log.Printf("[step 1] welcome resp: (code: %d msg: %s)", code, msg)
	return nil
}

func login(conn net.Conn, r *bufio.Reader) error {
	if err := sendCommand(conn, "USER ubuntu"); err != nil {
		return err
	}
	code, msg, err := read(r)
	if err != nil {
		return err
	}
	if code != 331 && code != 230 {
		return errors.New(fmt.Sprintf("invalid response code: %d, message: %s", code, msg))
	}
	if code != 331 {
		return nil
	}

	if err = sendCommand(conn, "PASS xxx"); err != nil {
		return err
	}
	code, msg, err = read(r)
	if err != nil {
		return err
	}
	if code != 230 {
		return errors.New(fmt.Sprintf("invalid response code: %d, message: %s", code, msg))
	}
	log.Printf("[step 2] login resp: (code: %d msg: %s)", code, msg)
	return nil
}

func setPassiveMode(conn net.Conn, r *bufio.Reader) (string, error) {
	if err := sendCommand(conn, "PASV"); err != nil {
		return "", err
	}
	code, msg, err := read(r)
	if err != nil {
		return "", err
	}
	if code != 227 {
		return "", errors.New(fmt.Sprintf("invalid response code: %d, message: %s", code, msg))
	}
	log.Printf("[step 3] pasv resp: (code: %d msg: %s)", code, msg)

	// 解析227响应
	// 227 Entering Passive Mode (192,168,7,251,73,167).
	start := strings.Index(msg, "(")
	end := strings.Index(msg, ")")
	if start == -1 || end == -1 {
		return "", errors.New("invalid response")
	}
	items := strings.Split(msg[start+1:end], ",")
	if len(items) != 6 {
		return "", errors.New("invalid response")
	}
	ip := strings.Join(items[0:4], ".")

	p1, err := strconv.Atoi(items[4])
	if err != nil {
		return "", err
	}
	p2, err := strconv.Atoi(items[5])
	if err != nil {
		return "", err
	}
	port := p1*256 + p2
	return fmt.Sprintf("%s:%d", ip, port), nil
}

func list(conn net.Conn, r *bufio.Reader, addr string) error {
	// 建立数据连接
	dataConn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer dataConn.Close()

	if err = sendCommand(conn, "LIST"); err != nil {
		return err
	}
	code, msg, err := read(r)
	if err != nil {
		return err
	}
	if code != 150 && code != 125 {
		return errors.New(fmt.Sprintf("invalid response code: %d, message: %s", code, msg))
	}

	reader := bufio.NewReader(dataConn)
	var lines []string
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		lines = append(lines, strings.TrimRight(line, "\r\n"))
	}

	code, msg, err = read(r)
	if err != nil {
		return err
	}
	if code != 226 && code != 250 {
		return errors.New(fmt.Sprintf("invalid response code: %d, message: %s", code, msg))
	}

	log.Println("[step 4] list resp")
	log.Println("---------------------------------------------------")
	for _, line := range lines {
		fmt.Println(line)
	}
	log.Println("---------------------------------------------------")

	return nil
}

func sendCommand(conn net.Conn, cmd string) error {
	cmd = cmd + "\r\n"
	_, err := conn.Write([]byte(cmd))
	return err
}

func read(r *bufio.Reader) (int, string, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return 0, "", err
	}
	line = strings.TrimSpace(line)
	if len(line) < 3 {
		return 0, "", errors.New("invalid response")
	}
	code, err := strconv.Atoi(line[:3])
	return code, line, err
}
