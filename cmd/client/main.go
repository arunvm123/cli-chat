package main

import (
	"log"

	"github.com/arunvm/chat_app/chat"
	"github.com/arunvm/chat_app/chat/client"
	"github.com/jroimartin/gocui"
	"google.golang.org/grpc"
)

func main() {
	var err error
	conn, err := grpc.Dial(":8888", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error connecting to server, Error=%v", err)
	}
	defer conn.Close()
	broadcastClient := chat.NewBroadcastClient(conn)

	chatClient := client.New(broadcastClient)

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatal(err)
	}
	defer g.Close()

	g.SetManagerFunc(client.Layout)

	g.SetKeybinding("name", gocui.KeyEnter, gocui.ModNone, chatClient.Update)
	g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, chatClient.Send)
	g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, chatClient.Quit)

	g.MainLoop()
}
