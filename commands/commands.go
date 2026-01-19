// Package commands handles Telegram bot command processing and routing.
package commands

import (
	"database/sql"
	"log"

	"github.com/bbeetlesam/imalrightjack-bot/database"
	"github.com/bbeetlesam/imalrightjack-bot/messages"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleMessage(update tgbotapi.Update, db *sql.DB, done <-chan struct{}) *tgbotapi.MessageConfig {
	if update.Message == nil {
		return nil
	}

	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	log.Println(messages.LogMessageReceived(update.Message.From.UserName, userID, update.Message.Text))

	if !update.Message.IsCommand() {
		return nil
	}

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
		responseMsg.Text = handleTransaction(update, db, userID, done)
	case "today":
		responseMsg.Text = handleTodayReport(db, userID, done)
	default:
		responseMsg.Text = messages.RespDefault
	}

	return &responseMsg
}

func handleTransaction(update tgbotapi.Update, db *sql.DB, userID int64, done <-chan struct{}) string {
	tx, userErrMsg := database.ParseTransactionMsg(update.Message.Text)
	if userErrMsg != "" {
		return userErrMsg
	}

	// check shutdown before db write (prevents duplicate transactions on restart)
	if shouldShutdown(done) {
		log.Println("Shutdown signal received, skipping transaction")
		return messages.RespTransactionFailed
	}

	if err := database.AddTransaction(db, userID, tx); err != nil {
		log.Println(messages.LogDBError(err))
		return messages.RespTransactionFailed
	}

	log.Println(messages.LogTransactionSaved(tx.Type, tx.Amount, userID))
	return messages.RespTransactionSuccess(tx.Type, tx.Amount, tx.Note)
}

func handleTodayReport(db *sql.DB, userID int64, done <-chan struct{}) string {
	// check shutdown before db read
	if shouldShutdown(done) {
		log.Println("Shutdown signal received, skipping report")
		return messages.RespDefault
	}

	transactions, totalAmount, err := database.GetTodayTransactions(db, userID)
	if err != nil {
		log.Printf("Failed to get today's transactions: %v", err)
		return "Failed to retrieve today's transactions. Please try again."
	}

	return messages.RespTodayTransactions(transactions, totalAmount)
}

func shouldShutdown(done <-chan struct{}) bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}
