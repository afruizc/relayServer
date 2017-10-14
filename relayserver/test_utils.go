package relayserver

import (
	"net"
	"fmt"
	"testing"
)

// Starts a server and connects clientNo clients to the
// server. It returns all endpoints for the established
// connections
func StartServerConnectClients(clientNo int) ([]net.Conn, []net.Conn) {
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

func ReadString(conn net.Conn, t *testing.T) string {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		t.Error(err)
	}

	return string(buf[:n])
}