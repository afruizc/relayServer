package main

import (
	"net"
	"flag"
	"fmt"
	"bufio"
	"strings"
	"io"
	"log"
	"time"
)

var addr string

type relayClient struct {
	relay net.Conn
	server net.Conn
}

func main() {
	relayHost := flag.String("host", "localhost", "Relay server host")
	relayPort := flag.Int("port", 8080, "Relay server port")
	flag.Parse()

	server, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		fmt.Println(server)
	}

	port := server.Addr().(*net.TCPAddr).Port
	addr = fmt.Sprintf("localhost:%d", port)

	go func() {
		for {
			c, err := server.Accept()
			if err != nil {
				panic(err)
			}

			go handleConnection(c)
		}
	}()

	fmt.Printf("Listening on: %s\n", addr)

	relayServerAddr := fmt.Sprintf("%s:%d", *relayHost, *relayPort)
	relayConn := connectToRelay(relayServerAddr)
	clients := make(map[int]*relayClient)

	for {
		processMessages(relayConn, *relayHost, clients)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024*32)
	for {
		nb, err := conn.Read(buf)
		if err != nil {
			log.Println("ERROR:", err)
			break
		}
		_, err = conn.Write(buf[:nb])
		if err != nil {
			log.Println("ERROR:", err)
			break
		}
		log.Printf("Client sent: % X \n", string(buf[:nb]))
	}
}

func connectToRelay(relayServerAddr string) net.Conn {
	relayConn, err := net.Dial("tcp", relayServerAddr)

	if err != nil {
		panic(err)
	}

	return relayConn
}

func processMessages(conn net.Conn, relayHost string, clients map[int]*relayClient) {
	bufReader := bufio.NewReader(conn)
	msg, err := bufReader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
	}

	msg = strings.TrimSpace(msg)

	switch {
	case strings.HasPrefix(msg, "[NEW]"):
		id, port := getIdAndPort(msg)
		rlAddr := fmt.Sprintf("%s:%d", relayHost, port)
		client := newClient(rlAddr)
		clients[id] = client
		go client.sync()
	case strings.HasPrefix(msg, "[CLOSE]"):
		time.Sleep(time.Millisecond)
		id := getId(msg)
		clients[id].close()
	}
}


func newClient(rlAddr string) (*relayClient) {
	connToRelay, err := net.Dial("tcp", rlAddr)
	if err != nil {
		panic(err)
	}

	connToServer, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}

	return &relayClient{connToRelay, connToServer}
}

func (c *relayClient) sync() {
	go func() {
		io.Copy(c.server, c.relay)
	}()

	io.Copy(c.relay, c.server)
}

func (c *relayClient) close() {
	c.relay.Close()
	c.server.Close()
}

func getId(s string) (id int) {
	opCode := "[CLOSE]"
	idx := strings.Index(s, opCode) + len(opCode)
	fmt.Sscanf(s[idx:], "%d", &id)
	return
}

func getIdAndPort(s string) (id, port int) {
	opCode := "[NEW]"
	idx := strings.Index(s, opCode) + len(opCode)
	fmt.Sscanf(s[idx:], "%d:%d", &port, &id)
	return
}
