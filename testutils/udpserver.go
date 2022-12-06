package testutils

import (
	"fmt"
	"net"
)

func UDPServer(ip string, port int, connChan chan net.Conn) (net.Listener, error) {

	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: port,
		IP:   []byte(ip),
	})
	if err != nil {
		panic(err)
	}

	defer conn.Close()
	fmt.Printf("server listening %s\n", conn.LocalAddr().String())

	for {
		message := make([]byte, 20)
		_, _, err := conn.ReadFromUDP(message[:])
		if err != nil {
			fmt.Println(err.Error())
		}

	}
}
