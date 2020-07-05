package main

import (
	"fmt"
	"net"
	"time"
)

func main() {

	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		panic(err)
	}

	message := fmt.Sprintln("/connect>John")

	conn.Write([]byte(message))
	message = fmt.Sprintln("/message>Hello world")
	conn.Write([]byte(message))

	b := make([]byte, len("Hello world"))
	_, err = conn.Read(b)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))

	message = fmt.Sprintln("/disconnect>John")
	conn.Write([]byte(message))
	<-time.After(time.Second * 15)
}
