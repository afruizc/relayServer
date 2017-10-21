package main

import (
	"fmt"
	"flag"
	"gitlab.com/afruizc/relayServer/relayserver"
)

func main() {
	port := flag.Int("port", 8080, "port to bind to")
	flag.Parse()

	relayServer, err := relayserver.NewServer(*port)
	if err != nil {
		println(err)
	}
	fmt.Println("Waiting for clients...")
	relayServer.Listen()
}
