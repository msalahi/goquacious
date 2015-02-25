package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <username>\n", os.Args[0])
		os.Exit(1)
	}
	_USERNAME = os.Args[1]
	chatApp := CreateChatApp()
	chatApp.MainLoop()
}
