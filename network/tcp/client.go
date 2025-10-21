package main

import (
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:3000")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte("hello, world\n"))
	if err != nil {
		panic(err)
	}
}
