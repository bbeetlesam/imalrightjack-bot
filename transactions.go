package main

import (
	"fmt"
	"strconv"
	"strings"
)

type TransactionInput struct {
	Type string
	Amount int64
	Note string
}

func parseTransactionMsg(msgText string) (*TransactionInput, error) {
	args := strings.SplitN(msgText, " ", 3)
	maxNoteLength := 75
	note := ""

	if len(args) < 2 {
		return nil, fmt.Errorf("bad usage.")
	}

	// parse command type [spend | earn]
	cmdType := strings.TrimPrefix(args[0], "/")

	// parse amount (int, not float)
	amount, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil || amount <= 0 {
		return nil, fmt.Errorf("invalid amount (numbers only)")
	}

	// parse note (truncated if length > 75)
	if len(args) >= 3 {
		note = args[2]
		if len(note) > maxNoteLength {
			note = note[:maxNoteLength]
		}
	}

	return &TransactionInput{ Type: cmdType, Amount: amount, Note: note }, nil
}
