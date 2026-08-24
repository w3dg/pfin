package command

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/w3dg/pfin/internal"
	"github.com/w3dg/pfin/tracker"
)

type Command struct {
	Desc     string
	Dispatch func(t *tracker.Tracker, args []string)
}

var Commands = map[string]Command{
	"show": {
		Desc: "Show all entries",
		Dispatch: func(t *tracker.Tracker, args []string) {
			t.PrettyPrint()
		},
	},
	"add": {
		Desc:     "Add a new entry",
		Dispatch: runAdd,
	},
	// more commands go here as you build them
}

func runAdd(t *tracker.Tracker, args []string) {
	fs := flag.NewFlagSet("add", flag.ExitOnError)
	etype := fs.String("etype", "", "income, aside, expense, topup")
	amount := fs.String("amount", "", "amount, e.g. 45.50")
	date := fs.String("date", time.Now().Format(time.DateOnly), "date, eg. 2026-08-24, if omitted, current date is used")
	category := fs.String("category", "", "category name")
	note := fs.String("note", "", "note")
	fs.Parse(args)

	switch {
	case *etype == "":
		pfatal("pfin add: -etype is required")
	case *amount == "":
		pfatal("pfin add: -amount is required")
	case *category == "":
		pfatal("pfin add: -category is required")
	}

	_, err := time.Parse(time.DateOnly, *date)
	if err != nil {
		pfatal("Could not parse date: ", *date, err.Error())
	}

	*category = strings.ToLower(*category)

	err = IsValidAdd(*etype, *category)

	if err != nil && *etype != "expense" {
		if errors.Is(err, ErrUnknownCategory) {
			pfatal("pfin add: unknown category for positive payment. Valid values: income, aside, topup")
		}

		if errors.Is(err, ErrIllegalCombo) {
			pfatal("pfin add: illegal combination, a non expense record cannot have the category:", *category)
		}

		if errors.Is(err, ErrIllegalIncomeCombo) {
			fmt.Fprintf(os.Stderr, "pfin add: illegal combination, a %v record cannot have category %v\n", *etype, *category)
			os.Exit(2)
		}
	}

	if err != nil && *etype == "expense" {
		if errors.Is(err, ErrUnknownCategory) {
			fmt.Fprintln(os.Stderr, "Available categories:")
			for k, v := range internal.NameToCategory {
				fmt.Fprintln(os.Stderr, k, v)
			}
			pfatal("pfin add: unknown category", *category)
		}

		if errors.Is(err, ErrIllegalCombo) {
			pfatal("pfin add: illegal combination with expense type entry")
		}
	}

	if *etype == "expense" {
		if err := t.AddExpense(*amount, *date, *category, *note); err != nil {
			log.Fatal("Error adding expense: ", err)
		}

	} else {
		// t.AddPositivePayment(*etype, *amount, *category, *note) — wire up once tracker supports it
		fmt.Printf("would add positive pay: etype=%s amount=%s category=%s note=%s\n", *etype, *amount, *category, *note)
	}
}

func IsValidAdd(etype, category string) error {
	categoryEnum, ok := internal.NameToCategory[category]
	if !ok {
		return ErrUnknownCategory
	}

	allowedCategoriesForPositivePayment := []internal.CategoryType{internal.INCOME_CATEGORY, internal.SET_ASIDE_CATEGORY, internal.TOP_UP_CATEGORY}
	if etype != "expense" {
		if !ok {
			return ErrUnknownCategory
		}

		if !slices.Contains(allowedCategoriesForPositivePayment, categoryEnum) {
			return ErrIllegalCombo
		}

		if (categoryEnum == internal.INCOME_CATEGORY && etype != "income") ||
			(categoryEnum == internal.SET_ASIDE_CATEGORY && etype != "aside") ||
			(categoryEnum == internal.TOP_UP_CATEGORY && etype != "topup") {
			return ErrIllegalIncomeCombo
		}
	} else if slices.Contains(allowedCategoriesForPositivePayment, categoryEnum) {
		return ErrIllegalCombo
	}

	return nil
}

func pfatal(s ...string) {
	fmt.Fprintln(os.Stderr, s)
	os.Exit(2)
}
