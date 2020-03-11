package main

import (
	"io"
	"log"
	"net"
	"sync"
)

type ClientServerSynchronizer interface {
	// Forwards all IO between conn1 and conn2
	// Blocks until both connections have been closed
	SynchronizeIO(clientConn, serverConn *net.TCPConn)
}

// These channels are used as flags for when
// a connection has read and is currently not
// done writing data. We dubbed them as `pending`
type ClientServerSyncImpl struct {
	c io.Writer // Original connection to the server. Used for notifications
}

func NewClientServerSynchronizer(serverConn io.Writer) ClientServerSynchronizer {
	return &ClientServerSyncImpl{serverConn}
}

func (df *ClientServerSyncImpl) SynchronizeIO(clientConn, serverConn *net.TCPConn) {
	df.sync(clientConn, serverConn)
}

// Signals the server that a client has disconnected
// expects both client and server to end the connections
// themselves
func (df *ClientServerSyncImpl) sync(client, server *net.TCPConn) {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		_, err := io.Copy(server, client)
		log.Println("Client disconnected", err)
		server.CloseWrite()
	}()

	go func() {
		defer wg.Done()
		_, err := io.Copy(client, server)
		log.Println("server disconnected", err)
		client.CloseWrite()
	}()

	wg.Wait()
}
