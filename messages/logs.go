package messages

import "fmt"

const (
	LogStart       string = "Waiting to get updates (from her)..."
	LogDBConnected string = "Successfully connected to the database. A present from Nancy!"
)

func LogMessageReceived(username string, userID int64, text string) string {
	return fmt.Sprintf("Message from %s (%d): %s", username, userID, text)
}

func LogTransactionSaved(act string, amount int64, userID int64) string {
	return fmt.Sprintf("Transaction saved: %s Rp. %d by user %d", act, amount, userID)
}

func LogBotAuthorised(botName string) string {
	return fmt.Sprintf("Bot authorised as: %s.", botName)
}

func LogDBError(err error) string {
	return fmt.Sprintf("Database error: %v", err)
}
