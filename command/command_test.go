package command

import (
	"errors"
	"fmt"
	"testing"
)

func TestIllegalCombinationForAdd(t *testing.T) {
	etype := "income"
	category := "food"
	if err := IsValidAdd(etype, category); !errors.Is(err, ErrIllegalCombo) {
		fmt.Printf("got %v want %v\n", err, ErrIllegalCombo)
		t.Fail()
	}

	etype = "aside"
	category = "food"
	if err := IsValidAdd(etype, category); !errors.Is(err, ErrIllegalCombo) {
		fmt.Printf("got %v want %v\n", err, ErrIllegalCombo)
		t.Fail()
	}

	etype = "topup"
	category = "food"
	if err := IsValidAdd(etype, category); !errors.Is(err, ErrIllegalCombo) {
		fmt.Printf("got %v want %v\n", err, ErrIllegalCombo)
		t.Fail()
	}
}

func TestIllegalIncomeCombinationForAdd(t *testing.T) {
	etype := "income"
	category := "aside"
	if err := IsValidAdd(etype, category); !errors.Is(err, ErrIllegalIncomeCombo) {
		fmt.Printf("got %v want %v\n", err, ErrIllegalIncomeCombo)
		t.Fail()
	}

	etype = "aside"
	category = "topup"
	if err := IsValidAdd(etype, category); !errors.Is(err, ErrIllegalIncomeCombo) {
		fmt.Printf("got %v want %v\n", err, ErrIllegalIncomeCombo)
		t.Fail()
	}

	etype = "topup"
	category = "income"
	if err := IsValidAdd(etype, category); !errors.Is(err, ErrIllegalIncomeCombo) {
		fmt.Printf("got %v want %v\n", err, ErrIllegalIncomeCombo)
		t.Fail()
	}

	etype = "expense"
	category = "income"
	if err := IsValidAdd(etype, category); !errors.Is(err, ErrIllegalCombo) {
		fmt.Printf("got %v want %v\n", err, ErrIllegalCombo)
		t.Fail()
	}
}

func TestCorrectAdd(t *testing.T) {
	etype := "income"
	category := "income"
	if err := IsValidAdd(etype, category); err != nil {
		fmt.Printf("got %v want %v\n", err, nil)
		t.Fail()
	}

	etype = "expense"
	category = "food"
	if err := IsValidAdd(etype, category); err != nil {
		fmt.Printf("got %v want %v\n", err, nil)
		t.Fail()
	}
}
