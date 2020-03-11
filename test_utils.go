package main

import (
	"bytes"
	"fmt"
	"net"
)

// Starts a server and connects clientNo clients to the
// server. It returns all endpoints for the established
// connections
func StartServerConnectClients(clientNo int) (*net.TCPListener, []*net.TCPConn, []*net.TCPConn) {
	const serverAddr = "localhost:9001"
	clientEndpoints := make([]*net.TCPConn, 0)
	serverEndpoints := make([]*net.TCPConn, 0)

	addr, err := net.ResolveTCPAddr("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error parsing tcp address", addr)
	}

	server := startServer(addr)
	serverEndpointsChannel := make(chan *net.TCPConn)
	go func() {
		for i := 0; i < clientNo; i++ {
			c, _ := server.AcceptTCP()
			fmt.Println("Client connected", )
			serverEndpointsChannel <- c
		}
	}()

	for i := 0; i < clientNo; i++ {
		c, _ := net.DialTCP("tcp", nil, addr)
		clientEndpoints = append(clientEndpoints, c)
		serverEndpoints = append(serverEndpoints, <-serverEndpointsChannel)
	}

	return server, serverEndpoints, clientEndpoints
}

func startServer(addr *net.TCPAddr) *net.TCPListener {

	server, err := net.ListenTCP("tcp", addr)
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
