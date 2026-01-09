package messages

import "fmt"

const (
	RespGreet   string = "Jack's here, mate. May I help you with those green stacks?"
	RespDefault string = "Jack doesn't know that command, sadly.\nMaybe you can try to type /help."
	RespHelp    string = "*Need help with Jack? Here are some info:*\n" +
		"/start - Wake Jack's mind. You'll receive a greeting from him.\n" +
		"/help - Show this messages, obviously. What are you expecting?\n" +
		"/earn - Log your income. Probably the best command here.\n" +
		"/spend - Log your expense. Jack hates subtraction, just saying.\n\n" +
		"*Got it?* Now time to fully utilize Jack's simple mind."
	RespAbout string = "Who's Jack?\n\n*imalrightjack*, usually called Jack, is a bot that helps you track " +
		"your financial flow, like your income and outcome. His mind is very simple and intuitive though, like " +
		"basic income and expense logs, daily report, et cetera.\nMind you, Jack loves anything that smells money, " +
		"so I think you can trust him."
	RespTransactionFailed string = "ERROR: Failed to save the transaction.\nPlease try again."
	RespErrAmount         string = "Please specify the amount, and optionally the notes."
	RespErrInvAmount      string = "Invalid amount! Use positive numbers only (e.g. 67000)\nNo commas, dots, and letters allowed."
)

const (
	LogStart            string = "Waiting to get updates (from her)..."
	LogDBConnected      string = "Successfully connected to the database. A present from Nancy!"
	LogTeletokenMissing string = "$TELETOKEN env variable not set."
	LogDBTokenMissing   string = "$TURSOTOKEN env variable not set."
	LogDBUrlMissing     string = "$TURSOURL env variable not set."
)

func RespTransactionSuccess(act string, amount int64, note string) string {
	noteText := "-"
	action := "spent"

	if note != "" {
		noteText = "\n" + note
	}
	if act == "earn" {
		action = "earned"
	}

	return fmt.Sprintf("Jack noted that you _%s_ *Rp. %d* with note: %s", action, amount, noteText)
}

func LogMessageReceived(username string, userID int64, text string) string {
	return fmt.Sprintf("Message from %s (%d): %s", username, userID, text)
}

func LogTransactionSuccess(act string, amount int64, userID int64) string {
	return fmt.Sprintf("Transaction saved: %s Rp. %d by user %d", act, amount, userID)
}

func LogBotAuthorised(botName string) string {
	return fmt.Sprintf("Bot authorised as: %s.", botName)
}

func LogDBError(err error) string {
	return fmt.Sprintf("Database error: %v", err)
}
