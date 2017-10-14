package relayserver

import (
	"net"
	"fmt"
)

// When a client connects to the relay server, we send a
// string of bytes to the server notifying them that a new
// connection has been made.
// We then expect them to setup a connection to our server.

// This means, we need 2 listeners per request:	one for requests
// from clients and one for requests from the server.

const addrFormat = "[NEW]%d\n"

type RelayRequestHandler interface {
	Run()
	GetClientPort() int
	GetServerPort() int
}

type RelayRequest struct {
	clientL                net.Listener // Listener for clients
	serverL                net.Listener // Listener for servers
	c                      net.Conn     // Connection to the original server
	clientPort, serverPort int
	ds                     DataSynchronizer // Data sync
}

func NewRelayRequest(client net.Conn, ds DataSynchronizer) (RelayRequestHandler, error) {
	cl, sl, err := startServers()
	if err != nil {
		return nil, err
	}

	cPort := getPort(cl.Addr())
	sPort := getPort(sl.Addr())

	return &RelayRequest{cl, sl, client, cPort, sPort, ds}, nil
}

func (rr *RelayRequest) Run() {
	for {
		clientConn, err := rr.clientL.Accept()
		if err != nil {
			panic(err)
		}

		notifyServer(rr.c, rr.serverPort)
		serverConn, err := rr.serverL.Accept()
		if err != nil {
			panic(err)
		}

		rr.ds.SynchronizeIO(clientConn, serverConn)
	}
}

func (rr *RelayRequest) GetClientPort() int {
	return rr.clientPort
}

func (rr *RelayRequest) GetServerPort() int {
	return rr.serverPort
}

// Notifies the server that a new connection has been
// established. After this we block waiting for a client
// to connect from the server
func notifyServer(conn net.Conn, port int) error {
	msg := fmt.Sprintf(addrFormat, port)
	_, err := conn.Write([]byte(msg))

	if err != nil {
		return err
	}

	return nil
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
