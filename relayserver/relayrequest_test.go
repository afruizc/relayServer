package relayserver

import (
	"testing"
	"net"
	"fmt"
	"bufio"
)

type mockSync struct {
	syncCalled bool
}

func NewMockSync() (*mockSync) {
	return &mockSync{false}
}

func (s *mockSync) SynchronizeIO(conn1, conn2 net.Conn) {
	s.syncCalled = true
}

func TestRelayRequestHandler_RunWithClients_Success(t *testing.T) {
	// Arrange
	ds := NewMockSync()
	c, s := StartServerConnectClients(1)
	rr, err := NewRelayRequest(c[0], ds)

	if err != nil {
		t.Error(err)
	}

	cAddr := fmt.Sprintf("localhost:%d", rr.GetClientPort())
	sAddr := fmt.Sprintf("localhost:%d", rr.GetServerPort())

	// Act
	go rr.Run()
	connectClient(cAddr, t)
	connectClient(sAddr, t)

	// Assert
	reader := bufio.NewReader(s[0])
	msg, err := reader.ReadString('\n')

	expectedMsg := fmt.Sprintf("[NEW]localhost:%d\n", rr.GetServerPort())
	if msg != expectedMsg {
		t.Errorf("Expected %s got %s from relayServer", expectedMsg, msg)
	}

	if !ds.syncCalled {
		t.Errorf("Sync method wasn't called")
	}
}

func connectClient(addr string, t *testing.T) {
	_, err := net.Dial("tcp", addr)
	if err != nil {
		t.Error(err)
	}
}


