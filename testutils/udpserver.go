package testutils

import (
	"fmt"
	"net"
)

const UDP_PACKET_SIZE = 65507

type MessageRecieved struct {
	Addr    *net.UDPAddr
	Message []byte
	Size    int
}

func UDPServer(ip string, port int, messageChanel chan MessageRecieved) (net.Listener, error) {

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
		message := MessageRecieved{
			Addr:    &net.UDPAddr{},
			Message: make([]byte, UDP_PACKET_SIZE),
		}
		message.Size, message.Addr, err = conn.ReadFromUDP(message.Message[:])

		if err != nil {
			fmt.Println(err.Error())
		}
		if messageChanel != nil {
			messageChanel <- message
		}

	}
}
