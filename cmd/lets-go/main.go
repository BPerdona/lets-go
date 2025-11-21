package main

import (
	"fmt"
	"lets-go/config"
	"lets-go/internal/server"
	"log"
)

func main() {
	config := config.NewServerConfig()
	srv := server.NewServer(config.ListenAddr)

	go func() {
		for msg := range srv.MessageChannel() {
			fmt.Printf("Received message from connection: (%s): %s", msg.From, string(msg.Payload))
		}
	}()

	log.Println("Server starting on", config.ListenAddr)
	log.Fatal(srv.Start())
}
