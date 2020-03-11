package main

import (
	"log"
	"net"
	"sync"
)

// When a client connects to the relay server, we send a
// string of bytes to the server notifying them that a new
// connection has been made.
// We then expect them to setup a connection to our server.

// This means, we need 2 listeners per request:	one for requests
// from clients and one for requests from the server.
type RelayRequestHandler interface {
	// Accepts clients, notifies the server and accepts a client
	// from the server.
	AcceptClients()

	// Return client port
	GetClientPort() int

	// Return server port
	GetServerPort() int
}

type RelayRequest struct {
	clientL                *net.TCPListener // Listener for clients
	serverL                *net.TCPListener // Listener for servers
	c                      *net.TCPConn     // Connection from the server
	clientPort, serverPort int
}

func NewRelayRequest(client *net.TCPConn) (RelayRequestHandler, error) {
	cl, sl, err := startServers()
	if err != nil {
		return nil, err
	}

	cPort := getPort(cl.Addr())
	sPort := getPort(sl.Addr())

	return &RelayRequest{cl, sl, client, cPort, sPort}, nil
}

func (rr *RelayRequest) AcceptClients() {
	var wg sync.WaitGroup

	for {
		clientSocket, err := rr.clientL.AcceptTCP()
		if err != nil {
			panic(err)
		}

		log.Println("Client connected:", clientSocket.RemoteAddr().String())

		notifyNewClient(rr.c, rr.serverPort)
		serverSocket, err := rr.serverL.AcceptTCP()
		if err != nil {
			panic(err)
		}
		log.Println("Server connected:", clientSocket.RemoteAddr().String())

		wg.Add(1)
		go func() {
			defer wg.Done()
			ds := NewClientServerSynchronizer(rr.c)
			ds.SynchronizeIO(clientSocket, serverSocket)
		}()
	}

	wg.Wait()
}

func (rr *RelayRequest) GetClientPort() int {
	return rr.clientPort
}

func (rr *RelayRequest) GetServerPort() int {
	return rr.serverPort
}

func startServers() (*net.TCPListener, *net.TCPListener, error) {
	clientL, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		return nil, nil, err
	}

	serverL, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		return nil, nil, err
	}

	return clientL.(*net.TCPListener), serverL.(*net.TCPListener), nil
}

func getPort(addr net.Addr) int {
	return addr.(*net.TCPAddr).Port
}
