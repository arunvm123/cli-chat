package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"github.com/arunvm/chat_app/chat/client"
	"github.com/jroimartin/gocui"
)

var conn net.Conn

func main() {

	var err error
	conn, err = net.Dial("tcp", ":8888")
	if err != nil {
		panic(err)
	}

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatal(err)
	}
	defer g.Close()

	g.SetManagerFunc(client.Layout)

	g.SetKeybinding("name", gocui.KeyEnter, gocui.ModNone, update)
	g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit)

	g.MainLoop()
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func update(g *gocui.Gui, v *gocui.View) error {
	b, err := ioutil.ReadAll(v)
	if err != nil {
		log.Fatalf("Error reading name, Error %v", err)
	}

	message := fmt.Sprintln("/connect>" + string(b))
	_, err = conn.Write([]byte(message))
	if err != nil {
		log.Fatalf("Error connecting to chat room, Error %v", err)
	}

	// Some UI changes
	g.SetViewOnTop("messages")
	g.SetViewOnTop("users")
	g.SetViewOnTop("input")
	g.SetCurrentView("input")

	return nil
}

// func main() {

// 	conn, err := net.Dial("tcp", ":8888")
// 	if err != nil {
// 		panic(err)
// 	}

// 	message := fmt.Sprintln("/connect>John")

// 	conn.Write([]byte(message))
// 	message = fmt.Sprintln("/message>Hello world")
// 	conn.Write([]byte(message))

// 	b := make([]byte, len("Hello world"))
// 	_, err = conn.Read(b)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(string(b))

// 	message = fmt.Sprintln("/disconnect>John")
// 	conn.Write([]byte(message))
// 	<-time.After(time.Second * 15)
// }
