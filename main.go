package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.mau.fi/whatsmeow/types/events"

	_ "github.com/mattn/go-sqlite3"
)

func eventHandler(e any) {
	switch v := e.(type) {
	case *events.Message:
		fmt.Println("Received message", v.Message.GetConversation())
	}
}

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()
	whatsAppClient := LaunchWhatsAppClient()

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	whatsAppClient.Disconnect()
}
