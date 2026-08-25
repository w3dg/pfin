package tracker

import "errors"

var (
	ErrInvalidExpenseAmount       = errors.New("Invalid expense amount entered, could not unmarshal")
	ErrUnknownExpenseCategory     = errors.New("Unknown category passed to add expense, not found in name to category map")
	ErrUnknownPositivePayCategory = errors.New("Unknown category passed to add for positive payment, not found in name to category map")
)
