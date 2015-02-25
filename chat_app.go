package main

import (
	"github.com/jroimartin/gocui"
	"log"
	"net"
)

type ChatApp struct {
	Gui                     *gocui.Gui
	DisplayQueue, SendQueue chan string
	Connections             map[string]net.Conn
	ListenerAddress         string
}

func CreateChatApp() ChatApp {
	gui := gocui.NewGui()
	if err := gui.Init(); err != nil {
		log.Fatal(err)
	}

	gui.ShowCursor = true
	connections := make(map[string]net.Conn)
	displayQueue := make(chan string)
	sendQueue := make(chan string)
	listenerAddress := net.JoinHostPort(_HOST, _PORT)
	chatApp := ChatApp{gui, displayQueue, sendQueue, connections, listenerAddress}

	gui.SetLayout(chatApp.RenderLayout)
	chatApp.setKeyBindings()

	return chatApp
}

func (chatApp ChatApp) MainLoop() {
	defer chatApp.Close()

	go chatApp.BootstrapPeers()
	go chatApp.Listen()
	go chatApp.DisplayMessages()
	go chatApp.SendMessages()

	err := chatApp.Gui.MainLoop()
	if err != nil && err != gocui.Quit {
		log.Panicln(err)
	}
}
func (chatApp ChatApp) Close() {
	chatApp.Gui.Close()
}
