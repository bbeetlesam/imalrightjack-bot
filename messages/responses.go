package messages

import (
	"fmt"
	"strconv"

	"github.com/bbeetlesam/imalrightjack-bot/models"
	"github.com/bbeetlesam/imalrightjack-bot/utils"
)

const (
	RespGreet   string = "Jack's here, mate\\. May I help you with those green stacks?"
	RespDefault string = "Jack doesn't know that command, sadly\\.\nMaybe you can try some /help\\."
	RespHelp    string = "*Need help with Jack? Here are some info:*\n" +
		"/start \\- Wake Jack's mind\\. You'll receive a greeting from him\\.\n" +
		"/help \\- Show this messages, obviously\\. Nothing more, nothing less\\.\n" +
		"/earn \\- Log your income\\. Probably the best command here\\.\n" +
		"/spend \\- Log your expense\\. Jack hates subtraction, just saying\\.\n\n" +
		"*Got it?* Now time to fully utilize Jack's simple mind\\."
	RespAbout string = "Who's Jack?\n\n*imalrightjack*, usually called Jack, is a bot that helps you track " +
		"your financial flow, like your income and outcome\\. His mind is very simple and intuitive though, like " +
		"basic income and expense logs, daily report, et cetera\\.\n\nMind you, Jack loves anything that smells money, " +
		"so I think you can trust him\\."
	RespTransactionFailed string = "ERROR: Failed to save the transaction.\nPlease try again\\."
	RespErrAmount         string = "Please specify the amount, and optionally the notes\\."
	RespErrInvalidAmount  string = "Invalid amount! Use positive numbers only (e.g. 67000)\nNo commas, dots, and letters allowed\\."
	RespFallbackMsg 			string = "⚠️ Sorry, there's a problem on Jack's voice currently.\n\n"
)

// will be later moved to somewhere proper, like configs
const currencySign = "Rp\\."

func RespTransactionSuccess(act string, amount int64, note string) string {
	noteText := "\\-"
	action := "spent"

	if note != "" {
		noteText = "\n" + utils.EscapeMarkdownV2(note)
	}
	if act == "earn" {
		action = "earned"
	}

	return fmt.Sprintf("Jack noted that you _%s_ *%s %d* with note: %s", action, currencySign, amount, noteText)
}

func RespTodayTransactions(transactions []models.Transaction, totalAmount int64) string {
	message := ""
	prefixEmoji := func(str string) string {
		if str == "spend" {
			return "➖ "
		}
		return "➕ "
	}

	if len(transactions) == 0 {
		message = "You have no transactions today, at least according to Jack's records."
	} else {
		message = "Your transactions today, recorded:\n\n"
		displayAmount := ""

		for _, transaction := range transactions {
			message += prefixEmoji(transaction.Type) + "\\[" + transaction.Time + "\\] "
			message += currencySign + " " + strconv.FormatInt(transaction.Amount, 10) + "\n"
		}

		if totalAmount >= 0 {
			displayAmount = currencySign + " " + strconv.FormatInt(totalAmount, 10)
		} else {
			displayAmount = "\\-" + currencySign + " " + strconv.FormatInt(-totalAmount, 10)
		}

		message += "\n💰 Total: " + displayAmount
	}

	return message
}
