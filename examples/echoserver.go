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
				fmt.Println("Damn")
				panic(err)
			}

			go handleConnection(c)
		}
	}()

	fmt.Printf("Listening on: %s\n", addr)

	relayServerAddr := fmt.Sprintf("%s:%d", *relayHost, *relayPort)
	relayConn := connectToRelay(relayServerAddr)

	for {
		processConnection(relayConn, *relayHost)
	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	clientReader := bufio.NewReader(conn)

	for {
		rawBytes, err := clientReader.ReadBytes('\n')
		log.Printf("Client sent: % X \n", rawBytes)
		if err != nil {
			break
		}
		conn.Write(rawBytes)
	}
}

func connectToRelay(relayServerAddr string) net.Conn {
	relayConn, err := net.Dial("tcp", relayServerAddr)

	if err != nil {
		panic(err)
	}

	return relayConn
}

func processConnection(conn net.Conn, relayHost string) {
	bufReader := bufio.NewReader(conn)
	msg, err := bufReader.ReadString('\n')

	if err != nil {
		fmt.Println(err)
	}

	laddr := strings.TrimSpace(getAddress(msg, relayHost))

	connToRelay, err := net.Dial("tcp", laddr)
	if err != nil {
		panic(err)
	}

	connToServer, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}

	go syncConns(connToRelay, connToServer)
}

func syncConns(conn1 net.Conn, conn2 net.Conn) {
	go io.Copy(conn1, conn2)
	io.Copy(conn2, conn1)

	time.Sleep(time.Second)

	conn1.Close()
	conn2.Close()
}

func getAddress(s string, host string) string {
	var port int
	opCode := "[NEW]"
	idx := strings.Index(s, opCode) + len(opCode)
	fmt.Sscanf(s[idx:], "%d", &port)

	return fmt.Sprintf("%s:%d", host, port)
}
