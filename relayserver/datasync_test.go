package relayserver

import (
	"testing"
	"net"
	"io"
	"strings"
	"sync"
	"bytes"
	"time"
)

type MockWriter struct {
	written string
}

func (w *MockWriter) Write(bytes []byte) (n int, err error) {
	s := string(bytes)
	w.written = w.written + s
	return len(s), nil
}

func TestSynchronizeIO_OneSends_Success(t *testing.T) {
	// Arrange
	const bytesToSend = 1
	client0Sent := genString(bytesToSend)
	serv, s, c := StartServerConnectClients(2)
	defer serv.Close()
	defer s[0].Close()
	defer s[1].Close()

	wr := &MockWriter{}
	syncIO := NewClientServerSynchronizer(wr)

	var wg sync.WaitGroup
	wg.Add(2)

	// Act
	// By synchronizing c0 and c1 we should be able to
	// read on s1 everything that is written to s0 and vice versa
	go func() {
		defer wg.Done()
		syncIO.SynchronizeIO(s[0], s[1])
	}()

	// Assert
	c[0].Write([]byte(client0Sent))

	// Close the client after 1 millisecond
	go func() {
		defer wg.Done()
		time.Sleep(time.Millisecond)
		c[0].Close()
		c[1].Close()
	}()

	assertConnReceived(c[1], client0Sent, t)
	wg.Wait()

	assertClientDisconnectedNotification(wr, 1, t)
}


func TestSynchronizeIO_SendAndReceive_Success(t *testing.T) {

	// Arrange
	const bytesToSend0 = 1
	const bytesToSend1 = 10000
	client0Sent := genString(bytesToSend0)
	client1Sent := genString(bytesToSend1)
	serv, s, c := StartServerConnectClients(2)
	defer serv.Close()
	defer s[0].Close()
	defer s[1].Close()

	wr := &MockWriter{}
	syncIO := NewClientServerSynchronizer(wr)

	var wg sync.WaitGroup
	wg.Add(2)

	// Act
	// By synchronizing c0 and c1 we should be able to
	// read on s1 everything that is written to s0 and vice versa
	go func() {
		defer wg.Done()
		syncIO.SynchronizeIO(s[0], s[1])
	}()

	// Assert
	c[0].Write([]byte(client0Sent))
	c[1].Write([]byte(client1Sent))

	// Close the client after 1 millisecond
	go func() {
		defer wg.Done()
		time.Sleep(time.Millisecond)
		c[0].Close()
		c[1].Close()
	}()

	assertConnReceived(c[1], client0Sent, t)
	assertConnReceived(c[0], client1Sent, t)
	wg.Wait()

	assertClientDisconnectedNotification(wr, 2, t)
}

func assertClientDisconnectedNotification(w *MockWriter, times int, t *testing.T) {
	lineNo := 0
	lines := strings.Split(w.written, "\n")
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		if !strings.HasPrefix(l, "[CLOSE]") {
			t.Errorf("Expected to find prefix [CLOSE], found %s on line %d",
				w.written, lineNo)
		}
		lineNo++
	}
}

func assertConnReceived(conn net.Conn, expected string, t *testing.T) {
	in, s, err := assertStrIn(conn, expected)

	if err != nil && err != io.EOF {
		t.Error(err)
	}

	if !in {
		t.Errorf("Expected %s got %s", expected, s)
	}
}

func genString(length int) string {
	var buf bytes.Buffer
	for i := 0 ; i < length ; i++ {
		buf.WriteByte('a')
	}

	return buf.String()
}

