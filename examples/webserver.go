package main

import (
	"net/http"
	"flag"
	"fmt"
//	"net"
//	"bufio"
//	"strings"
//	"io"
)
//const addr = "localhost:12345"

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main2() {
	http.HandleFunc("/", handler)
	go http.ListenAndServe(addr, nil)

	relayHost := flag.String("host", "localhost", "Relay server host")
	relayPort := flag.Int("port", 8080, "Relay server port")
	flag.Parse()

	relayServerAddr := fmt.Sprintf("%s:%d", *relayHost, *relayPort)
	relayConn := connectToRelay(relayServerAddr)

	for {
		processMessages(relayConn, *relayHost)
	}
}
