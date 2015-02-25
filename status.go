package main

import (
	"fmt"
)

func (chatApp ChatApp) UserJoinedStatusMessage(user string) {
	status := fmt.Sprintf("<%s> has joined!")
	chatApp.QueueMessageForDisplay(status)
}

func (chatApp ChatApp) UserLeftStatusMessage(user string) {
	status := fmt.Sprintf("<%s> has left :( :( :(")
	chatApp.QueueMessageForDisplay(status)
}
