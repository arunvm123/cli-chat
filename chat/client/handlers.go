package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/arunvm/chat_app/chat"
	"github.com/jroimartin/gocui"
)

func (chatClient *Client) Quit(g *gocui.Gui, v *gocui.View) error {
	message := fmt.Sprintf(chat.Disconnect + chatClient.Name)
	chatClient.Conn.Write([]byte(message))
	return gocui.ErrQuit
}

func (chatClient *Client) Update(g *gocui.Gui, v *gocui.View) error {

	chatClient.Name = v.Buffer()

	message := fmt.Sprintln(chat.Connect + chatClient.Name)
	_, err := chatClient.Conn.Write([]byte(message))
	if err != nil {
		log.Fatalf("Error connecting to chat room, Error %v", err)
	}

	// Some UI changes
	g.SetViewOnTop(chat.MessageView)
	g.SetViewOnTop(chat.UsersView)
	g.SetViewOnTop(chat.InputView)
	g.SetCurrentView(chat.InputView)
	go func(conn net.Conn, g *gocui.Gui) {
		messageView, err := g.View("messages")
		if err != nil {
			log.Fatalf("Error retrieving message view, Error: %v", err)
		}
		usersView, err := g.View("users")
		if err != nil {
			log.Fatalf("Error retrieving users view, Error: %v", err)
		}

		reader := bufio.NewReader(conn)

		for {
			data, _ := reader.ReadString('\n')
			msg := strings.TrimSpace(data)
			switch {
			case strings.HasPrefix(msg, chat.Users):
				splitUsers := strings.Split(strings.SplitAfter(msg, ">")[1], " ")
				users := strings.Join(splitUsers, "\n")
				g.Update(func(g *gocui.Gui) error {
					usersView.Title = fmt.Sprintf(" %d users: ", len(splitUsers))
					usersView.Clear()
					fmt.Fprintln(usersView, users)
					return nil
				})
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
	message := fmt.Sprintln(chat.Message + strings.TrimSpace(chatClient.Name) + ":" + v.Buffer())
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
