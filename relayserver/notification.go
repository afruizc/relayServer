package relayserver

import (
	"net"
	"fmt"
)

// New client connection.
// first int corresponds to the port.
// second int to a unique ID.
const newClientMsgFmt = "[NEW]%d:%d\n"

// Closed connection.
// int represents id of the client
const closedConnMsgFmt = "[CLOSE]%d\n"

func notifyNewClient(conn net.Conn, port, id int) error {
	s := fmt.Sprintf(newClientMsgFmt, port, id)
	return writeTo(conn, []byte(s))
}

func notifyClosedConnection(conn net.Conn, id int) error {
	if conn == nil {
		return nil
	}
	s := fmt.Sprintf(closedConnMsgFmt, id)
	return writeTo(conn, []byte(s))
}

// Write the given slice to the server
func writeTo(conn net.Conn, buf []byte) error {
	_, err := conn.Write(buf)

	if err != nil {
		return err
	}

	return nil
}
