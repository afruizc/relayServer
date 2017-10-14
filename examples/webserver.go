package main

import (
	"net/http"
	"fmt"
	"flag"
	"net"
	"bufio"
	"strings"
	"io"
)

const addr = "localhost:12345"

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", handler)
	go http.ListenAndServe(addr, nil)

	relayHost := flag.String("host", "localhost", "Relay server host")
	relayPort := flag.Int("port", 8080, "Relay server port")
	flag.Parse()

	relayServerAddr := fmt.Sprintf("%s:%d", *relayHost, *relayPort)
	relayConn := connectToRelay(relayServerAddr)

	for {
		waitForClient(relayConn, *relayHost)
	}
}

func connectToRelay(relayServerAddr string) net.Conn {
	relayConn, err := net.Dial("tcp", relayServerAddr)

	if err != nil {
		panic(err)
	}

	return relayConn
}

func waitForClient(conn net.Conn, host string) {
	bufReader := bufio.NewReader(conn)
	msg, err := bufReader.ReadString('\n')

	if err != nil {
		fmt.Println(err)
	}

	laddr := strings.TrimSpace(getAddress(msg, host))

	connToRelay, err := net.Dial("tcp", laddr)
	connToServer, err := net.Dial("tcp", addr)

	syncConns(connToRelay, connToServer)
}

func syncConns(conn1 net.Conn, conn2 net.Conn) {
	go io.Copy(conn1, conn2)
	go io.Copy(conn2, conn1)
}

func getAddress(s string, host string) string {
	var port int
	opCode := "[NEW]"
	idx := strings.Index(s, opCode) + len(opCode)
	fmt.Sscanf(s[idx:], "%d", &port)

	return fmt.Sprintf("%s:%d", host, port)
}
