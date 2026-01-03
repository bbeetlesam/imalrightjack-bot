package main

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOTKEN")) // use env var for early testing
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	fmt.Println("Hands off your stack next time.")
}
