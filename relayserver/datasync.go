package relayserver

import (
	"net"
	"io"
	"log"
	"time"
)

type DataSynchronizer interface {
	// Forwards all IO between conn1 and conn2
	// That is, everything that is read from conn1 is written to conn2
	// and vice versa
	SynchronizeIO(clientConn, serverConn net.Conn)
}

type DataSynchronizerImpl struct {
	clientServerCopy chan bool
	serverClientCopy chan bool
}

func NewDataSynchronizer() (DataSynchronizer) {
	return &DataSynchronizerImpl{
		make(chan bool),
		make(chan bool)}
}

func (df *DataSynchronizerImpl) SynchronizeIO(clientConn, serverConn net.Conn) {
	go sync(clientConn, serverConn) // Copy from the client
	go sync(serverConn, clientConn) // Copy from the server
}

func sync(conn1, conn2 net.Conn) {
	defer func() {
		conn1.Close()
		conn2.Close()
	}()
	n, err := io.Copy(conn1, conn2)
	if err == nil {
		// Client/Server disconnected successfully after writing
		// give 1 sec so all IO finishes.
		time.Sleep(time.Second)
	}
	log.Println("Connection closed.", "err:", err, "n:", n)
}

