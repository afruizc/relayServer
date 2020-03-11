package main

import (
	"flag"
	"fmt"
)

func main() {
	port := flag.Int("port", 8080, "port to bind to")
	flag.Parse()

	relayServer, err := NewServer(*port)
	if err != nil {
		println(err)
	}
	fmt.Println("Waiting for clients...")
	relayServer.Listen()
}
