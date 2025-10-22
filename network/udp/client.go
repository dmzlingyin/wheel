package main

import "net"

func main() {
	socket, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	if err != nil {
		panic(err)
	}
	defer socket.Close()

	_, err = socket.Write([]byte("hello world\n"))
	if err != nil {
		panic(err)
	}
}
