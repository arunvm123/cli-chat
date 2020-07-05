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
	<-time.After(time.Second * 15)
}
