package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var TELEBOT_TOKEN string = os.Getenv("TELETOKEN")

func main() {
	// open/connect to the database (remote)
	db, err := openDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// connect to the bot with its token
	bot, err := tgbotapi.NewBotAPI(TELEBOT_TOKEN)
	if err != nil {
		log.Fatal(err)
	}

	setBotCommands(bot)

	log.Printf("Authorised as: %s", bot.Self.UserName) // if succesfully connected

	u := tgbotapi.NewUpdate(0) // creates an UpdateConfig obj
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	bot.Debug = false

	log.Printf(START_LOG)

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
				sendMessage(bot, update.Message.Chat.ID, GREET_MSG)
			case "help":
				sendMessage(bot, update.Message.Chat.ID, HELP_MSG)
			case "about":
				sendMessage(bot, update.Message.Chat.ID, ABOUT_MSG)
			case "earn":
				sendMessage(bot, update.Message.Chat.ID, "earn what?")
			case "spend":
				sendMessage(bot, update.Message.Chat.ID, "spend what?")
			default:
				sendMessage(bot, update.Message.Chat.ID, DEFAULT_MSG)
			}
		}
	}
}

func sendMessage(bot *tgbotapi.BotAPI, chatID int64, msg string) {
	message := tgbotapi.NewMessage(chatID, msg)
	message.ParseMode = "Markdown"
	bot.Send(message)
}

func setBotCommands(bot *tgbotapi.BotAPI) {
	commands := []tgbotapi.BotCommand{
		{Command: "start", Description: "Wake Jack up"},
		{Command: "help", Description: "How this thing works"},
		{Command: "about", Description: "About Jack, and who made him"},
		{Command: "earn", Description: "Log an income"},
		{Command: "spend", Description: "Log an expense"},
		{Command: "today", Description: "Today’s damage report"},
	}

	config := tgbotapi.NewSetMyCommands(commands...)
	_, err := bot.Request(config)
	if err != nil {
		log.Println("Failed to set commands:", err)
	}
}
