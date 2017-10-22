package relayserver

import (
	"net"
	"io"
	"log"
	"sync"
)

type ClientServerSynchronizer interface {
	// Forwards all IO between conn1 and conn2
	// Blocks until both connections have been finished
	SynchronizeIO(clientConn, serverConn net.Conn)
}

// These channels are used as flags for when
// a connection has read and is currently not
// done writing data. We dubbed them as `pending`
type ClientServerSyncImpl struct {
	c  io.Writer // Original connection to the server. Used for notifications
	id int      // ID of the client
}

func NewClientServerSynchronizer(serverConn io.Writer, id int) (ClientServerSynchronizer) {
	return &ClientServerSyncImpl{serverConn, id}
}

func (df *ClientServerSyncImpl) SynchronizeIO(clientConn, serverConn net.Conn) {
	df.sync(clientConn, serverConn)
}

// Signals the server that a client has disconnected
// expects both client and server to end the connections
// themselves
func (df *ClientServerSyncImpl) sync(client, server net.Conn) {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		_, err := io.Copy(server, client)
		notifyClosedConnection(df.c, df.id)
		log.Println("Client disconnected", err)
	}()

	go func() {
		defer wg.Done()
		_, err := io.Copy(client, server)
		log.Println("server disconnected", err)
	}()

	wg.Wait()
}
