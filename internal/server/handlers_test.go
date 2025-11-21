package server

import (
	"net"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleHelp(t *testing.T) {
	srv := NewServer(":0")
	client, conn := net.Pipe()

	defer client.Close()
	defer conn.Close()
	defer srv.Stop()

	go srv.handleMessage(conn, []byte("/help\n"))

	buf := make([]byte, 1024)
	n, err := client.Read(buf)
	require.NoError(t, err)

	response := string(buf[:n])
	assert.Contains(t, response, "Available commands:")
}

func TestHandleQuit(t *testing.T) {
	srv := NewServer(":0")
	client, conn := net.Pipe()

	serverAddr := conn.RemoteAddr().String()
	srv.peers[serverAddr] = conn

	initialPeers := len(srv.peers)
	assert.Equal(t, 1, initialPeers)

	defer client.Close()
	defer conn.Close()
	defer srv.Stop()

	srv.handleMessage(conn, []byte("/quit\n"))

	assert.Equal(t, 0, len(srv.peers))

	// Verifica que a conexão foi fechada tentando escrever nela
	_, err := conn.Write([]byte("test"))
	assert.Error(t, err, "Connection should be closed")
}

func TestHandlePeers(t *testing.T) {
	srv := NewServer(":0")
	client1, conn1 := net.Pipe()
	client2, conn2 := net.Pipe()

	serverAddr1 := conn1.RemoteAddr().String()
	serverAddr2 := conn2.RemoteAddr().String()
	srv.peers[serverAddr1] = conn1
	srv.peers[serverAddr2] = conn2

	defer client1.Close()
	defer conn1.Close()
	defer client2.Close()
	defer conn2.Close()
	defer srv.Stop()

	go srv.handleMessage(conn1, []byte("/peers\n"))

	buf := make([]byte, 1024)
	n, err := client1.Read(buf)
	require.NoError(t, err)

	response := string(buf[:n])
	assert.Contains(t, response, "Total Peers: "+strconv.Itoa(len(srv.peers)))
}
