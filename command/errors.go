package command

import "errors"

var (
	ErrIllegalCombo       = errors.New("Illegal combination of arguments")
	ErrIllegalIncomeCombo = errors.New("Illegal combination of arguments for positive pay type")
	ErrUnknownCategory    = errors.New("Unknown Category")
)
