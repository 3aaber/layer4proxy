package testutils

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
)

func UDPClient(addrStr string, request chan []byte, serverResponse chan []byte) (net.Conn, error) {

	addr, err := net.ResolveUDPAddr("udp", addrStr)
	if err != nil {
		return nil, fmt.Errorf("could not resolve udp address %s: %v", addrStr, err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, fmt.Errorf("could not dial UDP addr %v: %v", addr, err)
	}

	response := bufio.NewReader(conn)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		wg.Done()
		for {
			userLine, ok := <-request
			if !ok {
				break
			}
			conn.Write(userLine)

			serverLine, err := response.ReadBytes(byte('\n'))
			switch err {
			case nil:
				serverResponse <- serverLine
			case io.EOF:
				os.Exit(0)
			default:
				fmt.Println("ERROR", err)
			}
		}
	}(wg)
	wg.Wait()

	return conn, nil
}
