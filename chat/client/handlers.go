package client

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/arunvm/chat_app/chat"
	"github.com/jroimartin/gocui"
)

func (chatClient *Client) Quit(g *gocui.Gui, v *gocui.View) error {
	chatClient.Conn.Disconnect(context.Background(), &chat.User{Name: chatClient.Name})
	return gocui.ErrQuit
}

func (chatClient *Client) Update(g *gocui.Gui, v *gocui.View) error {

	chatClient.Name = strings.TrimSpace(v.Buffer())

	stream, err := chatClient.Conn.Connect(context.Background(), &chat.User{Name: chatClient.Name})
	if err != nil {
		log.Fatalf("Error connecting to chat room, Error %v", err)
	}

	// Some UI changes
	g.SetViewOnTop(chat.MessageView)
	g.SetViewOnTop(chat.UsersView)
	g.SetViewOnTop(chat.InputView)
	g.SetCurrentView(chat.InputView)
	go func(s chat.Broadcast_ConnectClient, conn chat.BroadcastClient, g *gocui.Gui) {
		messageView, err := g.View(chat.MessageView)
		if err != nil {
			log.Fatalf("Error retrieving message view, Error: %v", err)
		}
		usersView, err := g.View(chat.UsersView)
		if err != nil {
			log.Fatalf("Error retrieving users view, Error: %v", err)
		}

		for {
			msg, err := s.Recv()
			if err != nil {
				log.Printf("Error while reading message,Error=%v", err)
				return
			}
			switch msg.Type {
			case chat.UserList:
				splitUsers := strings.Split(msg.GetMessage(), " ")
				users := strings.Join(splitUsers, "\n")
				g.Update(func(g *gocui.Gui) error {
					usersView.Title = fmt.Sprintf(" %d users: ", len(splitUsers))
					usersView.Clear()
					fmt.Fprintln(usersView, users)
					return nil
				})
			case chat.BroadcastMessage:
				g.Update(func(g *gocui.Gui) error {
					fmt.Fprint(messageView, msg.Message)
					return nil
				})
			}
		}
	}(stream, chatClient.Conn, g)

	return nil

}

func (chatClient *Client) Send(g *gocui.Gui, v *gocui.View) error {
	message := fmt.Sprint(strings.TrimSpace(chatClient.Name) + ":" + v.Buffer())
	_, err := chatClient.Conn.BroadcastMessage(context.Background(), &chat.Message{Message: message, Type: chat.BroadcastMessage})
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
