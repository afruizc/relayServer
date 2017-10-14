package relayserver

import (
	"net"
	"fmt"
	"log"
)

type Server interface {
	Listen()
}

type ServerImpl struct {
	server     net.Listener
	newClients chan net.Conn
	ds         DataSynchronizer
}

func NewServer(port int, ds DataSynchronizer) (Server, error) {
	laddr := fmt.Sprintf("localhost:%d", port)
	server, err := net.Listen("tcp", laddr)

	if err != nil {
		return nil, err
	}

	return &ServerImpl{server, make(chan net.Conn), ds}, nil
}

func (s *ServerImpl) Listen() {
	for {
		client, err := s.server.Accept()
		log.Printf("Client request to relay: %s\n", client)
		if err != nil {
			panic(err)
		}

		go s.startRelay(client)
	}
}

func (s *ServerImpl) startRelay(conn net.Conn) {
	relayRequest, err := NewRelayRequest(conn, s.ds)
	log.Printf("Serving relay requests on port %d", relayRequest.GetClientPort())
	if err != nil {
		panic(err)
	}

	relayRequest.Run()
}

// TODO tests this better
