package testutils

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"sync"
)

func echo(conn net.Conn) {
	r := bufio.NewReader(conn)
	for {
		line, err := r.ReadBytes(byte('\n'))
		if err != nil {
			if err != io.EOF {
				fmt.Println("failed to read data, err:", err)
			}
			return
		}

		conn.Write(line)
	}
}

func TcpEchoServer(addr string, connChan chan net.Conn) (net.Listener, error) {
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
			go echo(conn)
		}
	}(wg)
	wg.Wait()
	return l, nil
}
