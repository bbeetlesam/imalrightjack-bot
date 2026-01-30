// Package commands handles Telegram bot command processing and routing.
package commands

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/bbeetlesam/imalrightjack-bot/database"
	"github.com/bbeetlesam/imalrightjack-bot/messages"
	"github.com/bbeetlesam/imalrightjack-bot/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleMessage(ctx context.Context, update tgbotapi.Update, db *sql.DB) *tgbotapi.MessageConfig {
	if update.Message == nil {
		return nil
	}

	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	utils.LogColor("updt", messages.LogMessageReceived(
		update.Message.From.UserName,
		userID,
		update.Message.Text,
		update.Message.Date,
	))

	if !update.Message.IsCommand() {
		return nil
	}

	responseMsg := tgbotapi.NewMessage(chatID, "")
	responseMsg.ParseMode = "MarkdownV2"

	switch update.Message.Command() {
	case "start":
		responseMsg.Text = messages.RespGreet
	case "help":
		responseMsg.Text = messages.RespHelp
	case "about":
		responseMsg.Text = messages.RespAbout
	case "earn", "spend":
		responseMsg.Text = handleTransaction(ctx, update, db, userID)
	case "today":
		responseMsg.Text = handleTodayReport(ctx, db, userID)
	case "getlog":
		responseMsg.Text = handleGetLog(ctx, update, db, userID)
	default:
		responseMsg.Text = messages.RespDefault
	}

	return &responseMsg
}

func handleTransaction(ctx context.Context, update tgbotapi.Update, db *sql.DB, userID int64) string {
	tx, userErrMsg := database.ParseTransactionMsg(update.Message.Text)
	if userErrMsg != "" {
		return userErrMsg
	}

	// check shutdown before db write (prevents duplicate transactions on restart)
	if ctx.Err() != nil {
		utils.LogColor("warn", "Shutdown signal received, skipping transaction")
		return messages.RespTransactionFailed
	}

	var err error
	if tx.ID, err = database.AddTransaction(ctx, db, userID, tx); err != nil {
		utils.LogColor("errs", messages.LogDBError(err))
		return messages.RespTransactionFailed
	}

	utils.LogColor("dbwr", messages.LogTransactionSaved(tx.Type, tx.Amount, userID))
	return messages.RespTransactionSuccess(tx.Type, tx.ID, tx.Amount, tx.Note)
}

func handleTodayReport(ctx context.Context, db *sql.DB, userID int64) string {
	// check shutdown before db read
	if ctx.Err() != nil {
		utils.LogColor("warn", "Shutdown signal received, skipping report")
		return messages.RespDefault
	}

	transactions, totalAmount, err := database.GetTodayTransactions(ctx, db, userID)
	if err != nil {
		utils.LogColorf("errs", "Failed to get today's transactions: %v", err)
		return "Failed to retrieve today's transactions. Please try again."
	}

	return messages.RespTodayTransactions(transactions, totalAmount)
}

// currently applies for by ID
// TODO: use this for global getter (today, this month, week, date, etc)
func handleGetLog(ctx context.Context, update tgbotapi.Update, db *sql.DB, userID int64) string {
	if ctx.Err() != nil {
		utils.LogColor("warn", "idontkow")
		return messages.RespDefault
	}

	msg := utils.ParseCommandMsg(update.Message.Text)
	if len(msg) < 2 {
		return "Please specify the ID\\."
	} else if len(msg) > 2 {
		return "Provide only the ID number, not anything else\\."
	}

	txID, err := strconv.Atoi(msg[1])
	if err != nil {
		return messages.RespInvalidParse
	}

	tx, err := database.GetTransactionByID(ctx, db, userID, int64(txID))
	if err != nil {
		utils.LogColorf("errs", "failed to get transaction #%s: %v", msg[1], err)

		if err == sql.ErrNoRows {
			return messages.RespTransactionNotExist
		} else {
			return messages.RespInvalidParse
		}
	}

	return messages.RespDetailedTransaction(tx)
}
