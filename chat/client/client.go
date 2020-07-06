package client

import "net"

type Client struct {
	Conn net.Conn
	Name string
}

func New(conn net.Conn) *Client {
	return &Client{
		Conn: conn,
	}
}
