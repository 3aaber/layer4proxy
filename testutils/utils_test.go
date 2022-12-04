package testutils

import (
	"testing"
)

func TestEcho(t *testing.T) {

	addr := "127.0.0.1:9001"
	requestCh := make(chan []byte)
	responseCh := make(chan []byte)

	l, err := TcpEchoServer(addr, nil)
	if err != nil || l == nil {
		t.Error(err)
		return
	}
	defer l.Close()

	c, err := TcpClient(addr, requestCh, responseCh)
	if err != nil || c == nil {
		t.Error(err)
		return
	}
	defer c.Close()

	sampleText := "This Is Test\n"

	requestCh <- []byte(sampleText)
	responseTest := <-responseCh

	if string(responseTest) != sampleText {
		t.Error()
	}
}
