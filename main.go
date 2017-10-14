package main

import (
	"github.com/afruizc/relayServer/relayserver"
	"fmt"
)

func main() {

	ds := relayserver.NewDataSynchronizer()

	relayServer, err := relayserver.NewServer(8080, ds)
	if err != nil {
		println(err)
	}
	fmt.Println("Waiting for clients...")
	relayServer.Listen()
}
