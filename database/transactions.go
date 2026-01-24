package database

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/bbeetlesam/imalrightjack-bot/messages"
	"github.com/bbeetlesam/imalrightjack-bot/models"
	"github.com/bbeetlesam/imalrightjack-bot/utils"
)

func AddTransaction(db *sql.DB, userID int64, tx *models.Transaction) error {
	query := `INSERT INTO transactions (user_id, type, timestamp, amount, note)
		VALUES (?, ?, ?, ?, ?)
	;`

	timestamp := time.Now().UTC().Format(time.RFC3339)
	_, err := db.Exec(query, userID, tx.Type, timestamp, tx.Amount, tx.Note)
	return err
}

func ParseTransactionMsg(msgText string) (*models.Transaction, string) {
	args := strings.SplitN(msgText, " ", 3)
	maxNoteLength := 75
	maxAmount := 999999999999
	note := ""

	if len(args) < 2 {
		return nil, messages.RespErrAmount
	}

	// parse command type, and its bot name (after @) if any
	command := utils.ParseCommand(args[0])

	// parse amount (positive int, not float)
	amount, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil || amount <= 0 || amount >= int64(maxAmount) {
		return nil, messages.RespErrInvalidAmount
	}

	// parse note (truncated if length > 75)
	if len(args) >= 3 {
		note = args[2]
		if len(note) > maxNoteLength {
			note = note[:maxNoteLength]
		}
	}

	return &models.Transaction{Type: command.Action, Amount: amount, Note: note}, ""
}

func GetTodayTransactions(db *sql.DB, userID int64) ([]models.Transaction, int64, error) {
	now := time.Now() // use system timezone

	// calculate start and end of today (00:00:00) - (24:00:00) in the used timezone
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	// convert to UTC for database query
	startUTC := startOfDay.UTC().Format(time.RFC3339)
	endUTC := endOfDay.UTC().Format(time.RFC3339)

	query := `SELECT type, timestamp, amount, note
		FROM transactions
		WHERE user_id = ? AND timestamp >= ? AND timestamp < ?
		ORDER BY timestamp DESC
	;`

	rows, err := db.Query(query, userID, startUTC, endUTC)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var (
		transactions []models.Transaction
		totalAmount  int64
	)
	for rows.Next() {
		var transaction models.Transaction
		var timestamp string

		if err := rows.Scan(&transaction.Type, &timestamp, &transaction.Amount, &transaction.Note); err != nil {
			return nil, 0, err
		}

		txTime, err := time.Parse(time.RFC3339, timestamp)
		if err != nil {
			transaction.Time = ""
		} else {
			localTime := txTime.In(now.Location())
			transaction.Time = localTime.Format("15:04")
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
