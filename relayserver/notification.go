package relayserver

import (
	"fmt"
	"io"
)

// New client connection.
// first int corresponds to the port.
// second int to a unique ID.
const newClientMsgFmt = "[NEW]%d\n"

func notifyNewClient(socket io.Writer, port int) error {
	s := fmt.Sprintf(newClientMsgFmt, port)
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
