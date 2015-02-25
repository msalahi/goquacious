package main

import (
	"bufio"
	"log"
	"net"
	"strings"
)

func (chatApp ChatApp) StatusMessage(status string) {
	chatApp.QueueMessageForDisplay(status)
}

func sendMessage(conn net.Conn, message string) {
	conn.Write([]byte(message))
	conn.Write([]byte("\n"))
}

func (chatApp ChatApp) broadcastMessage(message string) {
	for _, conn := range chatApp.Connections {
		messageFromMe := formatUserMessage(_USERNAME, message)
		sendMessage(conn, messageFromMe)
	}
}

func (chatApp ChatApp) SendMessages() {
	var message string
	for {
		message = <-chatApp.SendQueue
		chatApp.broadcastMessage(message)
	}
}

func exchangeUserNames(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)
	user, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	sendMessage(conn, _USERNAME)
	return strings.TrimSpace(user), nil
}

func (chatApp ChatApp) handleConversation(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	user, err := exchangeUserNames(conn)
	if err != nil {
		log.Panic(err)
	}
	chatApp.Connections[user] = conn
	chatApp.Gui.Flush()
	for {
		message, err := reader.ReadString('\n')
		chatApp.QueueUserMessageForDisplay(user, message)
		if err != nil || message == "" {
			break
		}
	}
	delete(chatApp.Connections, user)
	chatApp.Gui.Flush()
}

func (chatApp ChatApp) QueueMessageForSend(message string) {
	chatApp.SendQueue <- message
}

func (chatApp ChatApp) Listen() {
	listener, err := net.Listen("tcp", chatApp.ListenerAddress)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go chatApp.handleConversation(conn)
	}
}
