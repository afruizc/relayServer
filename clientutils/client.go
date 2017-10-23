package clientutils

import (
	"net"
	"bufio"
	"fmt"
	"strings"
	"io"
	"sync"
	"log"
)

// Process messages reads from the connection to the relay server and
// acts based on messages received there. Only 1 type of message
// is supported right now:
// [NEW] message: "[NEW]<port>" where port number is an int
func ProcessMessages(tcpConn *net.TCPConn, relayHost string, serverAddr *net.TCPAddr) {
	var wg sync.WaitGroup

	for {
		bufReader := bufio.NewReader(tcpConn)
		msg, err := bufReader.ReadString('\n')
		if err != nil {
			fmt.Println("[RS] Error reading from connection", err)
			return
		}

		msg = strings.TrimSpace(msg)

		if !strings.HasPrefix(msg, "[NEW]") {
			fmt.Println("[RS] Error. Msg not understood", msg)
			return
		}
		processNew(msg, relayHost, serverAddr, &wg)
	}

	wg.Wait()
}

// Creates a new relayClient and starts synchronization between
// the relayServer and the server.
func processNew(msg string, relayHost string, serveraddr *net.TCPAddr,
		wg *sync.WaitGroup) {
	relayaddr := parseNewMsg(msg, relayHost)
	client, err := newClient(relayaddr, serveraddr)
	if err != nil {
		fmt.Println("[RS] Error creating client. Skipping")
		return
	}
	fmt.Println("[RS] New connection from client", client.relay.RemoteAddr())

	wg.Add(1)
	go func() {
		defer wg.Done()
		client.sync()
	}()
}

// Returns the id and address to connect to from the message
func parseNewMsg(msg, rlHost string) (*net.TCPAddr) {
	var port int
	fmt.Sscanf(msg, "[NEW]%d", &port)
	raddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", rlHost, port))

	if err != nil {
		log.Println("[RS] Error parsing tcp address", err)
	}
	return raddr
}

// Relay client has the server and
// relay server connections.
type relayClient struct {
	relay  *net.TCPConn
	server *net.TCPConn
	openClient chan int
}

func newClient(relayAddr, serverAddr *net.TCPAddr) (*relayClient, error) {
	tcpRelayConn, err := net.DialTCP("tcp", nil, relayAddr)
	if err != nil {
		return nil, err
	}

	tcpServerConn, err := net.DialTCP("tcp", nil, serverAddr)
	if err != nil {
		return nil, err
	}

	return &relayClient{tcpRelayConn, tcpServerConn,
		make(chan int, 1)}, nil
}

func (c *relayClient) sync() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		io.Copy(c.server, c.relay)
		c.server.CloseWrite()
	}()

	go func() {
		defer wg.Done()
		io.Copy(c.relay, c.server)
		c.relay.CloseWrite()
	}()

	wg.Wait()
}
