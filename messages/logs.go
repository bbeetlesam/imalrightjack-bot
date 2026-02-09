package messages

import (
	"fmt"
	"os"
	"time"

	"github.com/bbeetlesam/imalrightjack-bot/models"
)

const (
	LogStart       string = "Waiting to get updates (from her)..."
	LogDBConnected string = "Successfully connected to the database. A present from Nancy!"
	LogExitProgram string = "Shutdown complete. Fare thee well!"
)

func LogMessageReceived(username string, userID int64, text string, msgDate int) string {
	msgTime := time.Unix(int64(msgDate), 0)
	timestamp := msgTime.Format("2006-01-02 15:04:05 MST")

	return fmt.Sprintf("Message from %s (%d) [%s]: %s", username, userID, timestamp, text)
}

func LogTransactionSaved(act models.TransactionType, amount int64, userID int64) string {
	return fmt.Sprintf("Transaction saved: %s Rp. %d by user %d", act, amount, userID)
}

func LogBotAuthorised(botName string) string {
	return fmt.Sprintf("Bot authorised as: %s.", botName)
}

func LogDBError(err error) string {
	return fmt.Sprintf("Database error: %v", err)
}

func LogSignalOSReceived(sig os.Signal) string {
	return fmt.Sprintf("Received OS signal: %v", sig)
}
