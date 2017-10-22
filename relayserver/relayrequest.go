package relayserver

import (
	"net"
)



// When a client connects to the relay server, we send a
// string of bytes to the server notifying them that a new
// connection has been made.
// We then expect them to setup a connection to our server.

// This means, we need 2 listeners per request:	one for requests
// from clients and one for requests from the server.
type RelayRequestHandler interface {
	Run()
	GetClientPort() int
	GetServerPort() int
}


type RelayRequest struct {
	clientL                net.Listener // Listener for clients
	serverL                net.Listener // Listener for servers
	c                      net.Conn     // Connection from the server
	clientPort, serverPort int
	clientId               int // Incremental ID for clients
}

func NewRelayRequest(client net.Conn) (RelayRequestHandler, error) {
	cl, sl, err := startServers()
	if err != nil {
		return nil, err
	}

	cPort := getPort(cl.Addr())
	sPort := getPort(sl.Addr())

	return &RelayRequest{cl, sl, client, cPort, sPort, 0},nil
}

func (rr *RelayRequest) Run() {
	for {
		rr.clientId++
		clientConn, err := rr.clientL.Accept()
		if err != nil {
			panic(err)
		}

		notifyNewClient(rr.c, rr.serverPort, rr.clientId)
		serverConn, err := rr.serverL.Accept()
		if err != nil {
			panic(err)
		}
		ds := NewClientServerSynchronizer(rr.c, rr.clientId)

		go func() {
			defer func() {
				clientConn.Close()
				serverConn.Close()
			}()
			ds.SynchronizeIO(clientConn, serverConn)
		}()
	}
}

func (rr *RelayRequest) GetClientPort() int {
	return rr.clientPort
}

func (rr *RelayRequest) GetServerPort() int {
	return rr.serverPort
}


func startServers() (net.Listener, net.Listener, error) {
	clientL, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		return nil, nil, err
	}

	serverL, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		return nil, nil, err
	}

	return clientL, serverL, nil
}

func getPort(addr net.Addr) int {
	return addr.(*net.TCPAddr).Port
}
