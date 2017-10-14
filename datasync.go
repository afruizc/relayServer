package relayserver

import (
	"net"
	"io"
)

type DataSynchronizer interface {
	// Forwards all IO between conn1 and conn2
	// That is, everything that is read from conn1 is written to conn2
	// and vice versa
	SynchronizeIO(conn1, conn2 net.Conn)
}

type DataSynchronizerImpl struct {}

func NewDataSynchronizer() (*DataSynchronizerImpl) {
	return &DataSynchronizerImpl{}
}

func (df *DataSynchronizerImpl) SynchronizeIO(conn1, conn2 net.Conn) {
	go func() {
		_, err := io.Copy(conn2, conn1)

		if err != nil {
			panic(err)
		}
	}()

	go func() {
		_, err := io.Copy(conn1, conn2)

		if err != nil {
			panic(err)
		}
	}()
}

