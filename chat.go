package main

import (
	"bufio"
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
	"net"
	"strings"
)

var (
	USERNAME              = "murad"
	MESSAGE_DISPLAY_QUEUE = make(chan UserMessage)
	MESSAGE_SEND_QUEUE    = make(chan UserMessage)
	CONNECTIONS           = make(map[string]net.Conn)
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
	contactList.Clear()
	for user, _ := range CONNECTIONS {
		fmt.Fprintf(contactList, "• %s\n", user)
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

func queueMessageForDisplay(user string, message string) {
	MESSAGE_DISPLAY_QUEUE <- UserMessage{user: user, text: message}
}

func queueMessageForSend(user string, message string) {
	MESSAGE_SEND_QUEUE <- UserMessage{user: user, text: message}
}

func sendMessage(message UserMessage) {
	for user, conn := range CONNECTIONS {
		print(user)
		messageText := fmt.Sprintf("%s\n%s\n", message.user, message.text)
		conn.Write([]byte(messageText))
	}
}

func handleInputMessage(gui *gocui.Gui, chatInput *gocui.View) error {
	line, _ := chatInput.Line(0)
	queueMessageForDisplay(USERNAME, line)
	queueMessageForSend(USERNAME, line)
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
	if err := gui.SetKeybinding("chatInput", gocui.KeyEnter, gocui.ModNone, handleInputMessage); err != nil {
		return err
	}
	return nil
}

func listenForMessagesAndSend() {
	var message UserMessage
	for {
		message = <-MESSAGE_SEND_QUEUE
		sendMessage(message)
	}
}

func listenForMessagesAndDisplay(gui *gocui.Gui, displayViewName string) {
	var message UserMessage
	var view *gocui.View
	for {
		if view != nil {
			message = <-MESSAGE_DISPLAY_QUEUE
			displayMessageFromUser(view, message)
			gui.Flush()
		} else {
			view, _ = gui.View(displayViewName)
		}
	}
}

func waitForUser(reader *bufio.Reader) (string, error) {
	user, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(user), nil
}

func handleConversation(gui *gocui.Gui, conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	user, err := waitForUser(reader)
	if err != nil {
		log.Panic(err)
	}
	CONNECTIONS[user] = conn
	gui.Flush()
	for {
		message, err := reader.ReadString('\n')
		queueMessageForDisplay(user, message)
		if err != nil || message == "" {
			break
		}
	}
	delete(CONNECTIONS, user)
	gui.Flush()
}

func listenForConnections(localAddress string, gui *gocui.Gui) {
	listener, err := net.Listen("tcp", localAddress)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConversation(gui, conn)
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
	go listenForMessagesAndSend()
	go listenForConnections("127.0.0.1:55555", gui)

	err := gui.MainLoop()
	if err != nil && err != gocui.Quit {
		log.Panicln(err)
	}
}
