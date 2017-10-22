package relayserver

import (
	"net"
	"fmt"
	"log"
	"sync"
)

type Server interface {
	Listen()
}

type ServerImpl struct {
	server     *net.TCPListener
	newClients chan net.Conn
}

func NewServer(port int) (Server, error) {
	laddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	server, err := net.ListenTCP("tcp", laddr)

	if err != nil {
		return nil, err
	}

	return &ServerImpl{server, make(chan net.Conn)}, nil
}

func (s *ServerImpl) Listen() {
	var wg sync.WaitGroup

	for {
		client, err := s.server.AcceptTCP()
		if err != nil {
			panic(err)
		}

		log.Println("Request to relay from:", client.RemoteAddr())

		wg.Add(1)
		go startRelay(client, &wg)
	}

	wg.Wait()
}

func startRelay(conn *net.TCPConn, wg *sync.WaitGroup) {
	defer wg.Done()
	defer conn.Close()
	relayRequest, err := NewRelayRequest(conn)
	if err != nil {
		panic(err)
	}
	log.Printf("Serving relay requests on port %d", relayRequest.GetClientPort())

	relayRequest.AcceptClients()
}

// TODO tests this
