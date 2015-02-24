package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
	"strings"
	"time"
)

var (
	USERNAME = "murad"
	MESSAGES = make(chan UserMessage)
)

const (
	TITLE              = "ScrintedChat"
	TITLE_HEIGHT       = 2
	CONTACT_LIST_WIDTH = 30
	CHAT_INPUT_HEIGHT  = 10
)

type UserMessage struct {
	text string
	user string
}

func initTitleBar(titleBar *gocui.View) {
	fmt.Fprintln(titleBar, " > "+TITLE)
	titleBar.FgColor = gocui.ColorCyan
}

func drawTitleBar(gui *gocui.Gui) error {
	width, _ := gui.Size()
	titleBar, err := gui.SetView("titleBar", 0, 0, width-1, TITLE_HEIGHT)
	if err == gocui.ErrorUnkView {
		initTitleBar(titleBar)
		return nil
	}
	return err
}

func initContactList(contactList *gocui.View) {
	contactList.FgColor = gocui.ColorGreen
	fmt.Fprintln(contactList, "• SadComputer")
}

func drawContactList(gui *gocui.Gui) error {
	width, height := gui.Size()
	x0, y0 := width-CONTACT_LIST_WIDTH, TITLE_HEIGHT+1
	x1, y1 := width-1, height-CHAT_INPUT_HEIGHT
	contactList, err := gui.SetView("contactList", x0, y0, x1, y1)
	if err == gocui.ErrorUnkView {
		initContactList(contactList)
		return nil
	}
	return err
}

func initChatInput(chatInput *gocui.View) {
	chatInput.Editable = true
	chatInput.Highlight = true
	chatInput.Wrap = true
}

func drawChatInput(gui *gocui.Gui) error {
	width, height := gui.Size()
	x0, y0 := 0, height-CHAT_INPUT_HEIGHT+1
	x1, y1 := width-1, height-1
	chatInput, err := gui.SetView("chatInput", x0, y0, x1, y1)
	if err == gocui.ErrorUnkView {
		initChatInput(chatInput)
		return nil
	}
	return err
}

func initChatLog(chatLog *gocui.View) {
	chatLog.Autoscroll = true
	chatLog.Wrap = true
}

func drawChatLog(gui *gocui.Gui) error {
	width, height := gui.Size()
	x0, y0 := 0, TITLE_HEIGHT+1
	x1, y1 := width-CONTACT_LIST_WIDTH-2, height-CHAT_INPUT_HEIGHT
	chatLog, err := gui.SetView("chatLog", x0, y0, x1, y1)
	if err == gocui.ErrorUnkView {
		initChatLog(chatLog)
		return nil
	}
	return err
}

func layout(gui *gocui.Gui) error {
	if err := drawTitleBar(gui); err != nil {
		return err
	}
	if err := drawContactList(gui); err != nil {
		return err
	}
	if err := drawChatInput(gui); err != nil {
		return err
	}
	if err := drawChatLog(gui); err != nil {
		return err
	}
	if err := gui.SetCurrentView("chatInput"); err != nil {
		return err
	}
	return nil
}

func displayMessageFromUser(chatLog *gocui.View, message UserMessage) {
	text, user := strings.TrimSpace(message.text), message.user
	if text != "" {
		fmt.Fprintf(chatLog, "<%s>: %s\n", user, text)
	}
}

func queueMessage(user string, message string) {
	MESSAGES <- UserMessage{user: user, text: message}
}

func sendMessage(gui *gocui.Gui, chatInput *gocui.View) error {
	line, _ := chatInput.Line(0)
	queueMessage(USERNAME, line)
	chatInput.Clear()
	chatInput.SetCursor(0, 0)
	return nil
}

func sendQuit(gui *gocui.Gui, view *gocui.View) error {
	return gocui.Quit
}

func keyBindings(gui *gocui.Gui) error {
	if err := gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, sendQuit); err != nil {
		return err
	}
	if err := gui.SetKeybinding("chatInput", gocui.KeyEnter, gocui.ModNone, sendMessage); err != nil {
		return err
	}
	return nil
}

func listenForMessagesAndDisplay(gui *gocui.Gui, displayViewName string) {
	var message UserMessage
	var view *gocui.View
	for {
		if view != nil {
			message = <-MESSAGES
			displayMessageFromUser(view, message)
			gui.Flush()
		} else {
			view, _ = gui.View(displayViewName)
		}
	}
}

func talkToSadComputer() {
	for {
		queueMessage("SadComputer", "i am a computer")
		time.Sleep(2 * time.Second)
	}
}

func main() {
	gui := gocui.NewGui()
	if err := gui.Init(); err != nil {
		log.Fatal(err)
	}
	defer gui.Close()

	gui.SetLayout(layout)
	keyBindings(gui)
	gui.ShowCursor = true

	go listenForMessagesAndDisplay(gui, "chatLog")
	go talkToSadComputer()

	err := gui.MainLoop()
	if err != nil && err != gocui.Quit {
		log.Panicln(err)
	}
}
