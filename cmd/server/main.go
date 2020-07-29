package main

import (
	"log"
	"net"

	"github.com/arunvm/chat_app/chat"
	"github.com/arunvm/chat_app/chat/server"
	"google.golang.org/grpc"
)

func main() {
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("Error when starting a tcp listener, Error: %v", err)
	}

	log.Println("Server starting")

	s := grpc.NewServer()
	chat.RegisterBroadcastServer(s, server.New())
	log.Fatal(s.Serve(l))
}
