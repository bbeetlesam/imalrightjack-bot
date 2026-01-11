package main

import (
	"log"
	"strconv"

	"github.com/bbeetlesam/imalrightjack-bot/config"
	"github.com/bbeetlesam/imalrightjack-bot/database"
	"github.com/bbeetlesam/imalrightjack-bot/messages"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	botCfg, err := config.LoadBotConfig()
	if err != nil {
		log.Fatal(err)
	}

	// open/connect to the database (remote)
	db, err := database.Open(botCfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := database.InitSchema(db); err != nil {
		log.Fatal(err)
	}

	// connect to the bot with its token
	bot, err := tgbotapi.NewBotAPI(botCfg.TelebotToken)
	if err != nil {
		log.Fatal(err)
	}

	setBotCommands(bot)

	log.Println(messages.LogBotAuthorised(bot.Self.UserName)) // if successfully connected

	u := tgbotapi.NewUpdate(0) // creates an UpdateConfig obj
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	bot.Debug = false

	log.Println(messages.LogStart)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		userID := update.Message.From.ID
		chatID := update.Message.Chat.ID

		log.Println(messages.LogMessageReceived(update.Message.From.UserName, userID, update.Message.Text))

		if update.Message.IsCommand() {
			responseMsg := tgbotapi.NewMessage(chatID, "")
			responseMsg.ParseMode = "Markdown"

			switch update.Message.Command() {
			case "start":
				responseMsg.Text = messages.RespGreet
			case "help":
				responseMsg.Text = messages.RespHelp
			case "about":
				responseMsg.Text = messages.RespAbout
			case "earn", "spend":
				tx, userErrMsg := database.ParseTransactionMsg(update.Message.Text)
				if userErrMsg != "" {
					responseMsg.Text = userErrMsg
					break
				}

				if err := database.AddTransaction(db, userID, tx); err != nil {
					responseMsg.Text = messages.RespTransactionFailed
					log.Println(messages.LogDBError(err))
				} else {
					responseMsg.Text = messages.RespTransactionSuccess(tx.Type, tx.Amount, tx.Note)
					log.Println(messages.LogTransactionSaved(tx.Type, tx.Amount, userID))
				}
			case "today":
				transactions, totalAmount, err := database.GetTodayTransactions(db, userID)
				if err != nil {
					log.Println("cant read. placeholder")
					break
				}

				if len(transactions) == 0 {
					responseMsg.Text = "You have no transactions today, at least according to Jack's records."
				} else {
					responseMsg.Text = "Your transactions today, recorded:\n\n"
					for i, transaction := range transactions {
						responseMsg.Text += strconv.Itoa(i+1) + ". " + transaction.Type + " " + strconv.FormatInt(transaction.Amount, 10) + "\n"
					}
					responseMsg.Text += "\nTotal: Rp. " + strconv.FormatInt(totalAmount, 10)
				}
			default:
				responseMsg.Text = messages.RespDefault
			}

			if _, err := bot.Send(responseMsg); err != nil {
				log.Printf("Failed to send message: %v", err)
			}
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

	cmdConfig := tgbotapi.NewSetMyCommands(commands...)
	_, err := bot.Request(cmdConfig)
	if err != nil {
		log.Printf("Failed to set commands: %v", err)
	}
}
