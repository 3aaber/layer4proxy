package tcp

import (
	"fmt"
	"layer4proxy/testutils"
	"net"
	"testing"
	"time"
)

func getConnections(addr string) error {

	return nil
}

func TestEcho(t *testing.T) {

}
func TestProxy(t *testing.T) {

	// client --> proxy --> upstream
	// client <-- proxy <-- upstream

	upatreamAddress := "127.0.0.1:9000"
	proxyAddress := "127.0.0.1:9001"

	requestC := make(chan []byte)
	responseC := make(chan []byte)

	requestP := make(chan []byte)
	responseP := make(chan []byte)

	// Upstream Server
	// upstreamServer := make(chan net.Conn)
	u, err := testutils.TcpEchoServer(upatreamAddress, nil)
	if err != nil || u == nil {
		t.Error(err)
		return
	}

	// Proxy Server
	proxyServer := make(chan net.Conn)
	p, err := testutils.TcpServer(proxyAddress, proxyServer)
	if err != nil || p == nil {
		t.Error(err)
		return
	}

	// Proxy Client that connect to upstream
	proxyClient, err := testutils.TcpClient(upatreamAddress, requestC, responseC)
	if err != nil || proxyClient == nil {
		t.Error(err)
		return
	}

	// User Client that connect to proxy
	userClient, err := testutils.TcpClient(proxyAddress, requestP, responseP)
	if err != nil || userClient == nil {
		t.Error(err)
		return
	}

	from := <-proxyServer

	proxy(proxyClient, from, time.Second*100)
	proxy(from, proxyClient, time.Second*100)

	requestP <- []byte("This is Test \n")
	resp := <-responseP
	fmt.Println(string(resp))

}

func TestCopy(t *testing.T) {

}
