package tcp

import (
	"layer4proxy/testutils"
	"net"
	"testing"
)

func getConnections(addr string) error {

	return nil
}

func TestEcho(t *testing.T) {

}
func TestProxy(t *testing.T) {

	// client --> proxy --> upstream
	// client <-- proxy <-- upstream

	var addr string
	request := make(chan []byte)
	response := make(chan []byte)

	upstream := make(chan net.Conn)
	u, err := testutils.TcpEchoServer(addr, upstream)
	if err != nil || u == nil {
		t.Error(err)
		return
	}

	proxy := make(chan net.Conn)
	p, err := testutils.TcpServer(addr, proxy)
	if err != nil || p == nil {
		t.Error(err)
		return
	}

	client, err := testutils.TcpClient(addr, request, response)
	if err != nil || client == nil {
		t.Error(err)
		return
	}

}

func TestCopy(t *testing.T) {

}
