package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/jroimartin/gocui"
)

func Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (chatClient *Client) Update(g *gocui.Gui, v *gocui.View) error {

	chatClient.Name = v.Buffer()

	message := fmt.Sprintln("/connect>" + chatClient.Name)
	_, err := chatClient.Conn.Write([]byte(message))
	if err != nil {
		log.Fatalf("Error connecting to chat room, Error %v", err)
	}

	// Some UI changes
	g.SetViewOnTop("messages")
	g.SetViewOnTop("users")
	g.SetViewOnTop("input")
	g.SetCurrentView("input")
	go func(conn net.Conn, g *gocui.Gui) {
		messageView, err := g.View("messages")
		if err != nil {
			log.Fatalf("Error retrieving message view, Error: %v", err)
		}

		reader := bufio.NewReader(conn)

		for {
			data, _ := reader.ReadString('\n')
			msg := strings.TrimSpace(data)
			switch {
			default:
				g.Update(func(g *gocui.Gui) error {
					fmt.Fprintln(messageView, msg)
					return nil
				})
			}
		}
	}(chatClient.Conn, g)

	return nil

}

func (chatClient *Client) Send(g *gocui.Gui, v *gocui.View) error {
	message := fmt.Sprintln("/message>" + strings.TrimSpace(chatClient.Name) + ":" + v.Buffer())
	_, err := chatClient.Conn.Write([]byte(message))
	if err != nil {
		log.Fatalf("Error connecting to chat room, Error %v", err)
	}

	g.Update(func(g *gocui.Gui) error {
		v.Clear()
		v.SetCursor(0, 0)
		v.SetOrigin(0, 0)
		return nil
	})

	return nil

}
