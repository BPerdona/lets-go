package server

import (
	"fmt"
	"lets-go/internal/message"
	"net"
)

// Server represents the TCP server that manages peer connections
type Server struct {
	listenAddr string
	ln         net.Listener
	quitch     chan struct{}
	msgch      chan message.Message
	peers      map[string]net.Conn
}

// NewServer creates a new instance of the server
func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
		msgch:      make(chan message.Message, 1024),
		peers:      make(map[string]net.Conn),
	}
}

// Start starts the server and starts accepting connections
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	s.ln = ln

	go s.acceptLoop()

	<-s.quitch
	close(s.msgch)

	return nil
}

// Stop stops the server
func (s *Server) Stop() {
	close(s.quitch)
}

// MessageChannel returns the message channel
func (s *Server) MessageChannel() <-chan message.Message {
	return s.msgch
}

// acceptLoop accepts new connections continuously
func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		fmt.Println("New connection to the server", conn.RemoteAddr())

		s.peers[conn.RemoteAddr().String()] = conn
		go s.readLoop(conn)
	}
}

// readLoop reads messages from a specific connection
func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 2048)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			s.removePeer(conn)
			return
		}

		s.handleMessage(conn, buf[:n])
	}
}

// removePeer removes a peer from the list of connections
func (s *Server) removePeer(conn net.Conn) {
	delete(s.peers, conn.RemoteAddr().String())
	fmt.Println("Peer", conn.RemoteAddr(), "disconnected")
	conn.Close()
}
