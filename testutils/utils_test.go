package testutils

import (
	"net"
	"testing"
)

func TestEcho(t *testing.T) {

	addr := "127.0.0.1:9001"
	requestCh := make(chan []byte)
	responseCh := make(chan []byte)
	connChan := make(chan net.Conn)
	TcpEchoServer(addr, connChan)
	TcpClient(addr, requestCh, responseCh)
}
