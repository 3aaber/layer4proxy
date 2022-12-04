package testutils

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
)

func TcpClient(addr string, request chan []byte, serverResponse chan []byte) (net.Conn, error) {

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("ERROR", err)
		return nil, err
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
				os.Exit(2)
			}
		}
	}(wg)
	wg.Wait()
	return conn, nil
}
