package main

import (
	"net"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"gitlab.com/afruizc/relayServer/clientutils"
)

func main() {
	var wg sync.WaitGroup
	relayHost := flag.String("host", "localhost", "Relay server host")
	relayPort := flag.Int("port", 8080, "Relay server port")

	flag.Parse()

	laddr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		plnError("Error resolving address", err)
	}

	tcpServ, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		plnError("Cant start server", err)
		return
	}
	defer tcpServ.Close()

	addr := tcpServ.Addr().(*net.TCPAddr)

	wg.Add(1)
	go listen(tcpServ, &wg)
	fmt.Printf("Listening on: %s\n", addr)

	rsAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", *relayHost, *relayPort))
	if err != nil {
		plnError("Error resolving relayServer address", err)
	}

	relayConn, err := net.DialTCP("tcp", nil, rsAddr)
	if err != nil {
		plnError("Couldn't connect to relayServer", err)
		return
	}
	defer relayConn.Close()

	fmt.Println("Connected to relayServer", relayConn.RemoteAddr().String())
	clientutils.ProcessMessages(relayConn, *relayHost, addr)

	wg.Wait()
}

func listen(server *net.TCPListener, wg *sync.WaitGroup) {
	defer wg.Done()

	var localWg sync.WaitGroup

	for {
		c, err := server.Accept()
		if err != nil {
			panic(err)
		}

		tcpConn, ok := c.(*net.TCPConn)
		if !ok {
			plnError("Not a TCPConnection. Ignoring")
			c.Close()
			continue
		}

		fmt.Println("Client connected", tcpConn.RemoteAddr().String())
		localWg.Add(1)
		go handleConnection(tcpConn, &localWg)
	}

	localWg.Wait()
}

func handleConnection(tcpConn *net.TCPConn, wg *sync.WaitGroup) {
	defer wg.Done()
	defer tcpConn.Close()

	buf := make([]byte, 1024*32)
	for {
		nb, err := tcpConn.Read(buf)
		if err != nil {
			plnError("Error reading", err)
			break
		}
		log.Println("Server Read:", string(buf[:nb]))
		_, err = tcpConn.Write(buf[:nb])
		if err != nil {
			plnError("Error writing", err)
			break
		}
		log.Println("Server Wrote:", string(buf[:nb]))
	}
}

func plnError(s ...interface{}) {
	fmt.Fprintln(os.Stderr, s...)
}
