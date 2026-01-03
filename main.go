package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	GREET string = "Hello, I'm Jack. May I help you with those green stacks?"
	HELP  string = "Use /start to start, and /help to help.\nUse /spend or /earn to manage the stack."
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELETOKEN")) // connect to the bot with its token
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Authorised as: %s", bot.Self.UserName) // if succesfully connected

	u := tgbotapi.NewUpdate(0) // creates an UpdateConfig obj
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	bot.Debug = false

	log.Printf("Waiting to get updates (from her)...")

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf(
			"Message from %s (%d): %s", update.Message.From.UserName, update.Message.Chat.ID, update.Message.Text,
		)

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				sendMessage(bot, update.Message.Chat.ID, GREET)
			case "help":
				sendMessage(bot, update.Message.Chat.ID, HELP)
			}
		}
	}
}

func sendMessage(bot *tgbotapi.BotAPI, chatID int64, msg string) {
	message := tgbotapi.NewMessage(chatID, msg)
	bot.Send(message)
}
