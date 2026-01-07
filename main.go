package main

import (
	"fmt"
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

	log.Printf(startLog)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		userID := update.Message.From.ID
		chatID := update.Message.Chat.ID

		log.Printf(
			"Message from %s (%d): %s", update.Message.From.UserName, userID, update.Message.Text,
		)

		if update.Message.IsCommand() {
			responseMsg := tgbotapi.NewMessage(chatID, "")
			responseMsg.ParseMode = "Markdown"

			switch update.Message.Command() {
			case "start":
				responseMsg.Text = greetMsg
			case "help":
				responseMsg.Text = helpMsg
			case "about":
				responseMsg.Text = aboutMsg
			case "earn", "spend":
				tx, err := parseTransactionMsg(update.Message.Text)
				if err != nil {
					responseMsg.Text = "ERROR: " + err.Error()
					break
				}

				if err := addTransactionToDB(db, userID, tx); err != nil {
					responseMsg.Text = "ERROR: " + "Failed to save the transaction.\nPlease try again."
					log.Printf("Database error: %v", err)
				} else {
					note := "-"
					action := "earned"

					if tx.Note != "" {
						note = "\n" + tx.Note
					}
					if tx.Type == "spend" {
						action = "spent"
					}

					responseMsg.Text = fmt.Sprintf(
						"Jack noted that you _%s_ *Rp. %d* with note: %s", action, tx.Amount, note,
					)
					log.Printf("Transaction saved: %s Rp. %d by user %d", tx.Type, tx.Amount, userID)
				}
			default:
				responseMsg.Text = defaultMsg
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
