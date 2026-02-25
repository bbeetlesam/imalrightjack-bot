package domain

import "fmt"

type TransactionType string

const (
	TransactionTypeEarn  TransactionType = "earn"
	TransactionTypeSpend TransactionType = "spend"
)

type Transaction struct {
	ID     int64
	UserID int64
	Type   TransactionType
	Amount int64
	Note   string
	Time   int64 // in UTC
}

func (t *Transaction) IsValid() bool {
	if t.Amount <= 0 {
		return false
	}
	if t.Type != TransactionTypeEarn && t.Type != TransactionTypeSpend {
		return false
	}
	if t.UserID <= 0 {
		return false
	}
	return true
}

func (t *Transaction) FormatAmount() string {
	return formatCurrency(t.Amount)
}

func formatCurrency(amount int64) string {
	return fmt.Sprintf("%.2f", float64(amount)/100.0)
}
