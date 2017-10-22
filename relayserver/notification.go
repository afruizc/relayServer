package relayserver

import (
	"fmt"
	"io"
)

// New client connection.
// first int corresponds to the port.
// second int to a unique ID.
const newClientMsgFmt = "[NEW]%d:%d\n"

// Closed connection.
// int represents id of the client
const closedConnMsgFmt = "[CLOSE]%d\n"

func notifyNewClient(socket io.Writer, port, id int) error {
	s := fmt.Sprintf(newClientMsgFmt, port, id)
	return writeTo(socket, []byte(s))
}

func notifyClosedConnection(socket io.Writer, id int) error {
	s := fmt.Sprintf(closedConnMsgFmt, id)
	return writeTo(socket, []byte(s))
}

// Write the given slice to the server
func writeTo(socket io.Writer, buf []byte) error {
	_, err := socket.Write(buf)

	if err != nil {
		return err
	}

	return nil
}
