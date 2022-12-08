package testutils

import (
	"fmt"
	"net"
)

func UDPEchoServer(ip string, port int, messageChanel chan MessageRecieved) (net.Listener, error) {

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

		n, err := conn.Write(message.Message[:])
		if err != nil {
			fmt.Println(err.Error())
		}

		if n != message.Size {
			fmt.Println("byte sized in send and recieved are different")
		}

	}
}
