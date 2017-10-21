package relayserver

import (
	"net"
	"io"
	"log"
	"sync"
)

type ClientServerSynchronizer interface {
	// Forwards all IO between conn1 and conn2
	// That is, everything that is read from conn1 is written to conn2
	// and vice versa
	SynchronizeIO(clientConn, serverConn net.Conn)
}

// These channels are used as flags for when
// a connection has read and is currently not
// done writing data. We dubbed them as `pending`
type ClientServerSyncImpl struct {
	c  net.Conn // Original connection to the server. Used for notifications
	id int      // ID of the client
}

func NewClientServerSynchronizer(serverConn net.Conn, id int) (ClientServerSynchronizer) {
	return &ClientServerSyncImpl{serverConn, id}
}

func (df *ClientServerSyncImpl) SynchronizeIO(clientConn, serverConn net.Conn) {
	df.sync(clientConn, serverConn)
}

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
	server.Close()
	client.Close()
}
