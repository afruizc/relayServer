package relayserver

import (
	"testing"
	"net"
	"io"
	"time"
)

func TestSynchronizeIO_SendOneHalt_Success(t *testing.T) {
	// Arrange
	const bytesToSend = 1024 * 128
	client0Sent := genBytes(bytesToSend)
	s, c := StartServerConnectClients(2)

	// Act
	// By synchronizing c0 and c1 we should be able to
	// read on s1 everything that is written to s0 and vice versa
	forwarder := NewDataSynchronizer()
	forwarder.SynchronizeIO(s[0], s[1])

	c[0].Write(client0Sent)
	c[0].Close()
	time.Sleep(time.Second)

	// Assert
	assertConnReceived(c[1], bytesToSend, t)
	c[1].Close()
}

func TestSynchronizeIO_SendAndReceive_Success(t *testing.T) {
	// Arrange
	const sentBytesFrom0 = 8
	const sentBytesFrom1 = 1024
	client0Sent := genBytes(sentBytesFrom0)
	client1Sent := genBytes(sentBytesFrom1)
	s, c := StartServerConnectClients(2)

	// Act
	// By synchronizing c0 and c1 we should be able to
	// read on s1 everything that is written to s0 and vice versa
	forwarder := NewDataSynchronizer()
	forwarder.SynchronizeIO(s[0], s[1])

	c[0].Write(client0Sent)
	c[0].Close()
	c[1].Write(client1Sent)

	// Assert
	assertConnReceived(c[1], sentBytesFrom0, t)
	c[1].Close()
}

func assertConnReceived(conn net.Conn, expectedCount int, t *testing.T) {
	byteCount, err := GetByteCount(conn)

	if err != nil && err != io.EOF {
		t.Error(err)
	}

	if expectedCount != byteCount {
		t.Errorf("Expected %d got %d bytes", expectedCount, byteCount)
	}
}

func genBytes(length int) []byte {
	r := make([]byte, length)
	for i := 0 ; i < length ; i++ {
		r[i] = 'a'
	}

	return r
}

