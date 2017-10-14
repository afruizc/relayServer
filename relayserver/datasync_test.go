package relayserver

import (
	"testing"
	"net"
	"fmt"
)

func TestSynchronizeIO_SendOnBothEnds_Success(t *testing.T) {
	// Arrange
	const server0Sent = "data"
	const server1Sent = "more data"
	s, c := StartServerConnectClients(2)

	// Act
	// By synchronizing c0 and c1 we should be able to
	// read on s1 everything that is written to s0 and vice versa
	forwarder := NewDataSynchronizer()
	forwarder.SynchronizeIO(c[0], c[1])

	s[0].Write([]byte(server0Sent))
	s[1].Write([]byte(server1Sent))


	// Assert
	assertConnReceived(s[0], server1Sent, t)
	assertConnReceived(s[1], server0Sent, t)
}

func assertConnReceived(conn net.Conn, expected string, t *testing.T) {
	received := ReadString(conn, t)
	fmt.Println("received", received)
	if received != expected {
		t.Error("Expected", expected, "got", received)
	}
}

