package main

import (
	"github.com/jroimartin/gocui"
	"strings"
)

func (chatApp ChatApp) handleInputMessage(gui *gocui.Gui, chatInput *gocui.View) error {
	line, _ := chatInput.Line(0)
	line = strings.TrimSpace(line)
	if line != "" {
		chatApp.QueueUserMessageForDisplay(_USERNAME, line)
		chatApp.QueueMessageForSend(line)
	}
	chatInput.Clear()
	chatInput.SetCursor(0, 0)
	return nil
}

func sendQuit(gui *gocui.Gui, view *gocui.View) error {
	return gocui.Quit
}

func (chatApp ChatApp) setKeyBindings() error {
	if err := chatApp.Gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, sendQuit); err != nil {
		return err
	}
	if err := chatApp.Gui.SetKeybinding("chatInput", gocui.KeyEnter, gocui.ModNone, chatApp.handleInputMessage); err != nil {
		return err
	}
	return nil
}
