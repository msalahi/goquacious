package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"net"
)

func initTitleBar(titleBar *gocui.View) {
	fmt.Fprintln(titleBar, " > "+_TITLE)
	titleBar.FgColor = gocui.ColorCyan
}

func drawTitleBar(gui *gocui.Gui) error {
	width, _ := gui.Size()
	titleBar, err := gui.SetView(_TITLE_BAR_VIEW, 0, 0, width-1, _TITLE_HEIGHT)
	if err == gocui.ErrorUnkView {
		initTitleBar(titleBar)
		return nil
	}
	return err
}

func initContactList(contactList *gocui.View) {
	contactList.FgColor = gocui.ColorGreen
}

func drawContactList(gui *gocui.Gui, connections map[string]net.Conn) error {
	width, height := gui.Size()
	x0, y0 := width-_CONTACT_LIST_WIDTH, _TITLE_HEIGHT+1
	x1, y1 := width-1, height-_CHAT_INPUT_HEIGHT
	contactList, err := gui.SetView(_CONTACT_LIST_VIEW, x0, y0, x1, y1)
	if err == gocui.ErrorUnkView {
		initContactList(contactList)
		return nil
	}
	contactList.Clear()
	for user, _ := range connections {
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
	x0, y0 := 0, height-_CHAT_INPUT_HEIGHT+1
	x1, y1 := width-1, height-1
	chatInput, err := gui.SetView(_CHAT_INPUT_VIEW, x0, y0, x1, y1)
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
	x0, y0 := 0, _TITLE_HEIGHT+1
	x1, y1 := width-_CONTACT_LIST_WIDTH-2, height-_CHAT_INPUT_HEIGHT
	chatLog, err := gui.SetView(_CHAT_LOG_VIEW, x0, y0, x1, y1)
	if err == gocui.ErrorUnkView {
		initChatLog(chatLog)
		return nil
	}
	return err
}

func (chatApp ChatApp) RenderLayout(gui *gocui.Gui) error {
	if err := drawTitleBar(gui); err != nil {
		return err
	}
	if err := drawContactList(gui, chatApp.Connections); err != nil {
		return err
	}
	if err := drawChatInput(gui); err != nil {
		return err
	}
	if err := drawChatLog(gui); err != nil {
		return err
	}
	if err := gui.SetCurrentView(_CHAT_INPUT_VIEW); err != nil {
		return err
	}
	return nil
}
