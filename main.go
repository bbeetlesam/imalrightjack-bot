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
	initSchemaDB(db)

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
			responseMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			responseMsg.ParseMode = "Markdown"

			switch update.Message.Command() {
			case "start":
				responseMsg.Text = GREET_MSG
			case "help":
				responseMsg.Text = HELP_MSG
			case "about":
				responseMsg.Text = ABOUT_MSG
			case "earn":
				responseMsg.Text = "earn what"
			case "spend":
				responseMsg.Text = "spend what"
			default:
				responseMsg.Text = DEFAULT_MSG
			}

			bot.Send(responseMsg)
		}
	}
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
