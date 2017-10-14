package relayserver

import (
	"net"
	"io"
	"log"
)

type DataSynchronizer interface {
	// Forwards all IO between conn1 and conn2
	// That is, everything that is read from conn1 is written to conn2
	// and vice versa
	SynchronizeIO(conn1, conn2 net.Conn)
}

type DataSynchronizerImpl struct {}

func NewDataSynchronizer() (DataSynchronizer) {
	return &DataSynchronizerImpl{}
}

func (df *DataSynchronizerImpl) SynchronizeIO(conn1, conn2 net.Conn) {
	go sync(conn1, conn2)
	go sync(conn2, conn1)
}

func sync(conn1, conn2 net.Conn) {
	defer conn1.Close()
	defer conn2.Close()
	_, err := io.Copy(conn1, conn2)

	if err != nil {
		log.Println("Connection closed")
	}
}

