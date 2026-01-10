package database

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/bbeetlesam/imalrightjack-bot/messages"
)

type Transaction struct {
	Type   string
	Amount int64
	Note   string
}

func AddTransaction(db *sql.DB, userID int64, tx *Transaction) error {
	query := `INSERT INTO transactions (user_id, type, timestamp, amount, note)
		VALUES (?, ?, ?, ?, ?)
	;`

	timestamp := time.Now().Unix()
	_, err := db.Exec(query, userID, tx.Type, timestamp, tx.Amount, tx.Note)
	return err
}

func ParseTransactionMsg(msgText string) (*Transaction, string) {
	args := strings.SplitN(msgText, " ", 3)
	maxNoteLength := 75
	note := ""

	if len(args) < 2 {
		return nil, messages.RespErrAmount
	}

	// parse command type [spend | earn]
	cmdType := strings.TrimPrefix(args[0], "/")

	// parse amount (positive int, not float)
	amount, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil || amount <= 0 {
		return nil, messages.RespErrInvalidAmount
	}

	// parse note (truncated if length > 75)
	if len(args) >= 3 {
		note = args[2]
		if len(note) > maxNoteLength {
			note = note[:maxNoteLength]
		}
	}

	return &Transaction{Type: cmdType, Amount: amount, Note: note}, ""
}
