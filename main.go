package main

import (
	"github.com/afruizc/relayServer/relayserver"
	"fmt"
	"flag"
)

func main() {
	port := flag.Int("port", 8080, "port to bind to")
	flag.Parse()

	ds := relayserver.NewDataSynchronizer()

	relayServer, err := relayserver.NewServer(*port, ds)
	if err != nil {
		println(err)
	}
	fmt.Println("Waiting for clients...")
	relayServer.Listen()
}
