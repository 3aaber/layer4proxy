package tcp

import (
	"fmt"
	"layer4proxy/testutils"
	"net"
	"testing"
	"time"
)

func TestEcho(t *testing.T) {

}
func TestProxy(t *testing.T) {

	testMessage := []byte("This is Test \n")
	upatreamAddress := "127.0.0.1:9000"
	proxyAddress := "127.0.0.1:9001"

	// USER CLIENT <<===>> PROXY SERVER/CLIENT(127.0.0.1:9001) <<===>> UPSTREAM SERVER : (127.0.0.1:9000)

	requestC := make(chan []byte)
	responseC := make(chan []byte)

	requestP := make(chan []byte)
	responseP := make(chan []byte)

	// Upstream Server
	u, err := testutils.TcpEchoServer(upatreamAddress, nil)
	if err != nil || u == nil {
		t.Error(err)
		return
	}

	// Proxy Server
	proxyServerConnectionChannel := make(chan net.Conn)
	p, err := testutils.TcpServer(proxyAddress, proxyServerConnectionChannel)
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

	proxyServer := <-proxyServerConnectionChannel

	// Proxy TO <--> FROM
	proxy(proxyClient, proxyServer, time.Second*100)
	proxy(proxyServer, proxyClient, time.Second*100)

	// Send Message from client to Proxy
	requestP <- testMessage

	// Recieved Message Sent From Echo Server, Proxied from Proxy server
	resp := <-responseP

	if string(testMessage) != string(resp) {
		t.Errorf("Sent Recieved Message Missmatch , %s , %s", string(testMessage), string(resp))
	}

	fmt.Println(string(resp))

}

func TestCopy(t *testing.T) {

}
