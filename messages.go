package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"strings"
)

func formatUserMessage(user string, message string) string {
	return fmt.Sprintf("<%s>: %s", user, message)
}

func displayMessage(chatLog *gocui.View, message string) {
	text := strings.TrimSpace(message)
	if text != "" {
		fmt.Fprintln(chatLog, message)
	}
}

func (chatApp ChatApp) QueueMessageForDisplay(message string) {
	chatApp.DisplayQueue <- message
}

func (chatApp ChatApp) QueueUserMessageForDisplay(user string, message string) {
	chatApp.DisplayQueue <- formatUserMessage(user, message)
}

func (chatApp ChatApp) DisplayMessages() {
	var message string
	var view *gocui.View
	for {
		if view != nil {
			message = <-chatApp.DisplayQueue
			displayMessage(view, message)
			chatApp.Gui.Flush()
		} else {
			view, _ = chatApp.Gui.View(_CHAT_LOG_VIEW)
		}
	}
}
