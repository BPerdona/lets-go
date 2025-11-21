package server

import (
	"fmt"
	"lets-go/internal/message"
	"net"
	"strings"
)

// handleMessage processes incoming messages and executes commands
func (s *Server) handleMessage(conn net.Conn, data []byte) {
	msg := strings.TrimSpace(string(data))

	switch msg {
	case "/help":
		s.HandleHelp(conn)
	case "/quit":
		s.HandleQuit(conn)
	case "/peers":
		s.HandlePeers(conn)
	case "/broadcast":
		s.HandleBroadcast(conn)
	default:
		s.HandleDefaultMessage(conn, data)
	}
}

// handleHelp displays available commands
func (s *Server) HandleHelp(conn net.Conn) {
	conn.Write([]byte("Available commands: /quit, /peers, /help, /broadcast\n"))
}

// handleQuit disconnects the peer
func (s *Server) HandleQuit(conn net.Conn) {
	s.removePeer(conn)
}

// handlePeers displays the total number of connected peers
func (s *Server) HandlePeers(conn net.Conn) {
	fmt.Fprintf(conn, "Total Peers: %d\n", len(s.peers))
}

// handleBroadcast requests and sends a broadcast message
func (s *Server) HandleBroadcast(conn net.Conn) {
	conn.Write([]byte("Enter the message to broadcast: "))
	buf := make([]byte, 2048)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading from connection:", err)
		return
	}
	msg := strings.TrimSpace(string(buf[:n]))
	totalPeers := s.broadcastMessage(msg, conn.RemoteAddr().String())
	fmt.Fprintf(conn, "Message Broadcasted to %d peers\n", totalPeers)
}

// handleDefaultMessage processes normal messages (not commands)
func (s *Server) HandleDefaultMessage(conn net.Conn, data []byte) {
	s.msgch <- message.Message{
		From:    conn.RemoteAddr().String(),
		Payload: data,
	}
	conn.Write([]byte("Message Sent\n"))
}

// broadcastMessage sends a message to all peers except the sender
func (s *Server) broadcastMessage(msg string, from string) int {
	if len(s.peers) == 0 {
		return 0
	}
	count := 0
	for addr, peer := range s.peers {
		if addr == from {
			continue
		}
		peer.Write([]byte("[" + from + "]: " + msg + "\n"))
		count++
	}
	return count
}
