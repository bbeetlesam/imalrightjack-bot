package database

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/bbeetlesam/imalrightjack-bot/messages"
	"github.com/bbeetlesam/imalrightjack-bot/models"
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

	return &models.Transaction{Type: cmdType, Amount: amount, Note: note}, ""
}

func GetTodayTransactions(db *sql.DB, userID int64) ([]models.Transaction, int64, error) {
	jakartaLoc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		jakartaLoc = time.FixedZone("WIB", 7*60*60)
	}

	now := time.Now().In(jakartaLoc)

	// calculate start and end of today (00:00:00) - (24:00:00) in Jakarta
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, jakartaLoc)
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
