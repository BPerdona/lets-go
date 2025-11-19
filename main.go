package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type Message struct {
	from    string
	payload []byte
}

type Server struct {
	listenAddr string
	ln         net.Listener
	quitch     chan struct{}
	msgch      chan Message
	peers      map[string]net.Conn
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
		msgch:      make(chan Message, 1024),
		peers:      make(map[string]net.Conn),
	}
}

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

func (s *Server) broadcastMessage(msg string, from string) int {
	if len(s.peers) == 0 {
		return 0
	}
	for _, peer := range s.peers {
		if peer.RemoteAddr().String() == from {
			continue
		}
		peer.Write([]byte("[" + from + "]: " + msg + "\n"))
	}
	return len(s.peers) - 1
}

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

func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 2048)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			delete(s.peers, conn.RemoteAddr().String())
			fmt.Println("Peer", conn.RemoteAddr(), "disconnected")
			conn.Close()
			return
		}

		msg := strings.TrimSpace(string(buf[:n]))

		switch msg {
		case "/help":
			conn.Write([]byte("Available commands: /quit, /peers, /help, /broadcast\n"))
		case "/quit":
			delete(s.peers, conn.RemoteAddr().String())
			fmt.Println("Peer", conn.RemoteAddr(), "disconnected")
			conn.Close()
			return
		case "/peers":
			fmt.Fprintf(conn, "Total Peers: %d\n", len(s.peers))
		case "/broadcast":
			conn.Write([]byte("Enter the message to broadcast: "))
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Println("Error reading from connection:", err)
				continue
			}
			msg = strings.TrimSpace(string(buf[:n]))
			totalPeers := s.broadcastMessage(msg, conn.RemoteAddr().String())
			fmt.Fprintf(conn, "Message Broadcasted to %d peers\n", totalPeers)
		default:
			s.msgch <- Message{
				from:    conn.RemoteAddr().String(),
				payload: buf[:n],
			}
			conn.Write([]byte("Message Sent\n"))
		}
	}
}

func main() {
	server := NewServer(":8080")

	go func() {
		for msg := range server.msgch {
			fmt.Printf("Received message from connection: (%s): %s", msg.from, string(msg.payload))
		}
	}()

	log.Fatal(server.Start())
}
