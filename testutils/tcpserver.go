package testutils

import (
	"fmt"
	"net"
	"sync"
)

func TcpServer(addr string, connChan chan net.Conn) (net.Listener, error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("ERROR", err)
		return nil, err
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		wg.Done()
		for {
			conn, err := l.Accept()
			if err != nil {
				fmt.Println("ERROR", err)
				return
			}
			if connChan != nil {
				connChan <- conn
			}
		}
	}(wg)
	wg.Wait()
	return l, nil
}
