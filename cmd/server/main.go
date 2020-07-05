package main

import (
	"log"
	"net"

	"github.com/arunvm/chat_app/chat"
)

func main() {
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("Error when starting a tcp listener, Error: %v", err)
	}

	log.Println("Server starting")

	c := chat.New()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalf("Error when accepting a connection, Error: %v", err)
		}

		go c.Handle(&conn)
	}
}
