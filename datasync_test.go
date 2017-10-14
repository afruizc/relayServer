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
	s, c := startServerConnectClients(2)

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
	received := readString(conn, t)
	fmt.Println("received", received)
	if received != expected {
		t.Error("Expected", expected, "got", received)
	}
}

func readString(conn net.Conn, t *testing.T) string {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		t.Error(err)
	}

	return string(buf[:n])
}

// This is a dummy server that sends the vowels
// to the first `sendToClients`
func startServerConnectClients(clientNo int) ([]net.Conn, []net.Conn) {
	const serverAddr = "localhost:9001"
	clientEndpoints := make([]net.Conn, 0)
	serverEndpoints := make([]net.Conn, 0)

	server := startServer(serverAddr)

	serverEndpointsChannel := make(chan net.Conn)
	go func() {
		for i := 0; i < clientNo; i++ {
			c, _ := server.Accept()
			fmt.Println("Client connected", )
			serverEndpointsChannel <- c
		}
	}()

	for i := 0; i < clientNo; i++ {
		c, _ := net.Dial("tcp", serverAddr)
		clientEndpoints = append(clientEndpoints, c)
		serverEndpoints = append(serverEndpoints, <-serverEndpointsChannel)
	}

	return serverEndpoints, clientEndpoints
}

func startServer(serverAddr string) net.Listener {
	server, err := net.Listen("tcp", serverAddr)
	if err != nil {
		fmt.Println("ERROR:", err)
	}

	return server
}
