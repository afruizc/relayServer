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
}

func NewServer(port int) (Server, error) {
	laddr := fmt.Sprintf(":%d", port)
	server, err := net.Listen("tcp", laddr)

	if err != nil {
		return nil, err
	}

	return &ServerImpl{server, make(chan net.Conn)}, nil
}

func (s *ServerImpl) Listen() {
	for {
		client, err := s.server.Accept()
		if err != nil {
			panic(err)
		}
		log.Printf("Client request to relay: %s\n", client)

		go startRelay(client)
	}
}

func startRelay(conn net.Conn) {
	defer conn.Close()
	relayRequest, err := NewRelayRequest(conn)
	if err != nil {
		panic(err)
	}
	log.Printf("Serving relay requests on port %d", relayRequest.GetClientPort())

	relayRequest.Run()
}

// TODO tests this
