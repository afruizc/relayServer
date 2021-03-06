package main

import (
	"bufio"
	"fmt"
	"net"
	"testing"
	"time"
)

func TestRelayRequestHandler_RunWithClients_Success(t *testing.T) {
	// Arrange
	serv, c, s := StartServerConnectClients(1)
	rr, err := NewRelayRequest(c[0])

	if err != nil {
		t.Error(err)
	}

	cAddr := fmt.Sprintf("localhost:%d", rr.GetClientPort())
	sAddr := fmt.Sprintf("localhost:%d", rr.GetServerPort())

	// Act
	go rr.AcceptClients()
	connectClient(cAddr, t)
	connectClient(sAddr, t)

	// Assert
	reader := bufio.NewReader(s[0])
	msg, err := reader.ReadString('\n')

	expectedMsg := fmt.Sprintf("[NEW]%d\n", rr.GetServerPort())
	if msg != expectedMsg {
		t.Errorf("Expected %s got %s from relayServer", expectedMsg, msg)
	}

	time.Sleep(time.Millisecond)

	serv.Close()
}

func connectClient(addr string, t *testing.T) {
	_, err := net.Dial("tcp", addr)
	if err != nil {
		t.Error(err)
	}
}
