package main

import (
	"fmt"
	"os"

	"github.com/thenanox/ficsit-discord-bot/cmd"
)

func main() {
	token := os.Getenv("DISCORD_BOT_TOKEN")

	if err := cmd.Execute(token); err != nil {
		fmt.Printf("Error executing cmd %v", err)
		os.Exit(1)
	}
}
