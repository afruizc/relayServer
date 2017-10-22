package relayserver

import (
	"net"
	"fmt"
	"bytes"
)

// Starts a server and connects clientNo clients to the
// server. It returns all endpoints for the established
// connections
func StartServerConnectClients(clientNo int) (net.Listener, []net.Conn, []net.Conn) {
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

	return server, serverEndpoints, clientEndpoints
}

func startServer(serverAddr string) net.Listener {
	server, err := net.Listen("tcp", serverAddr)
	if err != nil {
		fmt.Println("ERROR:", err)
	}

	return server
}

func readAllFromConn(conn net.Conn) (string, error) {
	var ans bytes.Buffer
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			return ans.String(), err
		}

		ans.WriteString(string(buf[:n]))
	}

	return ans.String(), nil
}

func assertStrIn(conn net.Conn, str string) (bool, string, error) {
	var ans bytes.Buffer
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			return ans.String() == str, ans.String(), err
		}

		ans.WriteString(string(buf[:n]))
		if ans.String() == str {
			break
		}
	}

	return ans.String() == str, ans.String(), nil
}
