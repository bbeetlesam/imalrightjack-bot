package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/bbeetlesam/imalrightjack-bot/commands"
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

	bot.Debug = false
	updates := bot.GetUpdatesChan(u)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	done := make(chan struct{})

	var wg sync.WaitGroup

	log.Println(messages.LogStart)

	wg.Go(func() {
		for {
			select {
			case <-done:
				return
			case update := <-updates:
				responseMsg := commands.HandleMessage(update, db, done)
				if responseMsg != nil {
					if _, err := bot.Send(responseMsg); err != nil {
						log.Printf("Failed to send message: %v", err)
						fallbackMsg := tgbotapi.NewMessage(responseMsg.ChatID, messages.RespFallbackMsg)
						fallbackMsg.ParseMode = "Markdown"

						if _, err := bot.Send(fallbackMsg); err != nil {
							log.Printf("Failed to send fallback message: %v", err)
						}
					}
				}
			}
		}
	})

	sig := <-sigChan
	log.Println(messages.LogSignalOSReceived(sig))

	close(done)
	bot.StopReceivingUpdates()

	wg.Wait()
	log.Println(messages.LogExitProgram)
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
