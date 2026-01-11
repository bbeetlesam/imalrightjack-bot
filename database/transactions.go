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

	timestamp := time.Now().Format(time.RFC3339)
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

func GetTodayTransactions(db *sql.DB, userID int64) ([]Transaction, int64, error) {
	today := time.Now().Format("2006-01-02")
	query := `SELECT type, timestamp, amount, note
		FROM transactions
		WHERE user_id = ? AND timestamp LIKE ?
		ORDER BY timestamp DESC
	;`

	rows, err := db.Query(query, userID, today+"%")
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var (
		transactions []Transaction
		totalAmount  int64
	)
	for rows.Next() {
		var transaction Transaction
		var timestamp string

		if err := rows.Scan(&transaction.Type, &timestamp, &transaction.Amount, &transaction.Note); err != nil {
			return nil, 0, err
		}

		switch transaction.Type {
		case "spend":
			totalAmount -= transaction.Amount
		case "earn":
			totalAmount += transaction.Amount
		}

		transactions = append(transactions, transaction)
	}

	return transactions, totalAmount, nil
}
