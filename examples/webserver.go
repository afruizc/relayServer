package main

import (
	"net/http"
	"flag"
	"fmt"
	"sync"
	"net"
	"gitlab.com/afruizc/relayServer/clientutils"
	"os"
)
const addr = "localhost:12345"

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	var wg sync.WaitGroup

	serverAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		plnError("Error parsing address")
	}

	http.HandleFunc("/", handler)
	wg.Add(1)
	go func() {
		defer wg.Done()
		http.ListenAndServe(addr, nil)
	}()

	relayHost := flag.String("host", "localhost", "Relay server host")
	relayPort := flag.Int("port", 8080, "Relay server port")
	flag.Parse()

	relayServerAddr, err := net.ResolveTCPAddr("tcp",
		fmt.Sprintf("%s:%d", *relayHost, *relayPort))

	if err != nil {
		plnError("error resolving tcp address", err)
		return
	}

	relayConn, err := net.DialTCP("tcp", nil, relayServerAddr)
	if err != nil {
		plnError("error connecting to relay server", err)
	}

	clientutils.ProcessMessages(relayConn, *relayHost, serverAddr)

	wg.Wait()
}

func plnError(s ...interface{}) {
	fmt.Fprintln(os.Stderr, s...)
}
