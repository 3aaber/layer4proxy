package session

import (
	"layer4proxy/core"
	"layer4proxy/testutils"
	"net"
	"testing"
)

func TestNewSession(t *testing.T) {
	messageChanel := make(chan testutils.MessageRecieved)
	testutils.UDPEchoServer("127.0.0.1", 8050, messageChanel)

	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8050")
	if err != nil {
		t.Errorf("could not resolve udp address %s: %v", "127.0.0.1:8050", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		t.Errorf("problem in dial udp, %s", err.Error())
	}

	ss := NewSession(addr, conn, core.Upstream{}, Config{})

	ss.Write([]byte("test"))

	<-messageChanel
}
