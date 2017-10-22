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
// acts based on messages received there. Right now, only 2 types of messages
// are supported: [NEW] and [CLOSE]. This blocks on Reading from the
// TCP Connection.
func ProcessMessages(tcpConn *net.TCPConn, relayHost string, serverAddr *net.TCPAddr) {
	var wg sync.WaitGroup

	clients := &sync.Map{}
	for {
		bufReader := bufio.NewReader(tcpConn)
		msg, err := bufReader.ReadString('\n')
		if err != nil {
			fmt.Println("[RS] Error reading from connection", err)
			return
		}

		msg = strings.TrimSpace(msg)

		switch {
		case strings.HasPrefix(msg, "[NEW]"):
			// Adds to the wg
			processNew(msg, relayHost, clients, serverAddr, &wg)
		case strings.HasPrefix(msg, "[CLOSE]"):
			processClose(msg, clients)
		}
	}

	wg.Wait()
}

// Creates a new relayClient and starts synchronization between
// the relayServer and the server. This blocks until the
// synchronization is done.
func processNew(msg string, relayHost string, clients *sync.Map,
		serveraddr *net.TCPAddr, wg *sync.WaitGroup) {
	id, relayaddr := resolveNewMessage(msg, relayHost)
	client, err := newClient(relayaddr, serveraddr)
	if err != nil {
		fmt.Println("[RS] Error creating client. Skipping")
		return
	}
	fmt.Println("[RS] New connection from client", id, client.relay.RemoteAddr())
	clients.Store(id, client)

	wg.Add(1)
	go func() {
		defer wg.Done()
		client.sync()
	}()
}

func processClose(msg string, clients *sync.Map) {
	id := getId(msg)
	clientInt, ok := clients.Load(id)
	if !ok {
		fmt.Println("[RS] Can't find client with id", id)
	}

	client := clientInt.(*relayClient)
	client.close()
	clients.Delete(id)
	fmt.Println("[RS] Client with id", id, "disconnected")
}

// Parses the message and
func resolveNewMessage(msg, rlHost string) (int, *net.TCPAddr) {
	var id, port int
	fmt.Sscanf(msg, "[NEW]%d:%d", &port, &id)
	raddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", rlHost, port))

	if err != nil {
		log.Println("[RS] Error parsing tcp address", err)
	}
	return id, raddr
}

func getId(s string) (id int) {
	opCode := "[CLOSE]"
	idx := strings.Index(s, opCode) + len(opCode)
	fmt.Sscanf(s[idx:], "%d", &id)
	return
}

// Relay client has the server and
// relay server connections.
type relayClient struct {
	relay  *net.TCPConn
	server *net.TCPConn
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

	return &relayClient{tcpRelayConn, tcpServerConn}, nil
}

func (c *relayClient) sync() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		copyStream(c.server, c.relay)
	}()

	go func() {
		defer wg.Done()
		copyStream(c.relay, c.server)
	}()

	wg.Wait()
}

// Closes the connection making sure
// everything is flushed first
func (c *relayClient) close() {
	//time.Sleep(time.Millisecond)
	c.server.Close()
	c.relay.Close()
}

func copyStream(dst, src *net.TCPConn) (written int64, err error)  {
	buf := make([]byte, 32 * 1024)
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			log.Println("Read Relay:", string(buf[:nr]))
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
				log.Println("Wrote Relay:", string(buf[:nw]))
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return written, err

}