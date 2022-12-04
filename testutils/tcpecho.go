package testutils

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
)

func echo(conn net.Conn) {
	r := bufio.NewReader(conn)
	for {
		line, err := r.ReadBytes(byte('\n'))
		switch err {
		case nil:
			break
		case io.EOF:
		default:
			fmt.Println("ERROR", err)
		}
		conn.Write(line)
	}
}

func TcpEchoServer(addr string, connChan chan net.Conn) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		wg.Done()
		for {
			conn, err := l.Accept()
			if err != nil {
				fmt.Println("ERROR", err)
				continue
			}
			connChan <- conn
			go echo(conn)
		}
	}(wg)
	wg.Wait()
}
