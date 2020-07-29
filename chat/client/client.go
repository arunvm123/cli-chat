package client

import (
	"github.com/arunvm/chat_app/chat"
)

type Client struct {
	Conn chat.BroadcastClient
	Name string
}

func New(conn chat.BroadcastClient) *Client {
	return &Client{
		Conn: conn,
	}
}
