package tracker

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/w3dg/pfin/config"
	"github.com/w3dg/pfin/db"
	"github.com/w3dg/pfin/internal"
)

type Tracker struct {
	Cfg     config.Config
	Entries []internal.Entry
}

func NewTracker(cfg config.Config) Tracker {
	t := Tracker{
		Cfg:     cfg,
		Entries: nil,
	}

	t.init()
	return t
}

func (t *Tracker) init() {
	if strings.Contains(t.Cfg.DataDir, "~") {
		homedir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal("pfin config contains ~ as part of the data dir setting, however it could not be resolved. Exiting.")
		}
		t.Cfg.DataDir = strings.Replace(t.Cfg.DataDir, "~", homedir, 1) // 1 to replace the first occurence
	}

	dataFile := path.Join(t.Cfg.DataDir, "pfin.jsonl")

	_, err := os.Stat(dataFile)
	if err != nil {
		fmt.Printf("No data file - pfin.jsonl found at %v, creating a new one.\n", t.Cfg.DataDir)
		err := os.MkdirAll(t.Cfg.DataDir, 0755)
		if err != nil {
			log.Fatal("Error creating data dir - ", t.Cfg.DataDir, err)
		}
		_, err = os.Create(dataFile)
		if err != nil {
			log.Fatal("Error creating the pfin.jsonl data file", err)
		}
	}

	reader, err := db.NewReader(dataFile)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	entries, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Error reading entries ", err)
	}

	t.Entries = entries
}

// t.AddExpense(*amount, *category, *note) — wire up once tracker supports it
func (t *Tracker) AddExpense(amount, date, category, note string) error {
	a := internal.Amount{}

	// force double quotes around the amount
	if err := a.UnmarshalJSON([]byte(fmt.Sprintf("\"%v\"", amount))); err != nil {
		log.Fatal("Could not parse amount")
	}

	categoryEnum, ok := internal.NameToCategory[category]
	if !ok {
		return ErrUnknownExpenseCategory
	}

	e := internal.Entry{
		EntryType: internal.EXPENSE,
		Date:      date,
		Amount:    a,
		Category:  categoryEnum,
		Notes:     note,
	}

	dbw, err := db.NewWriter(t.Cfg.GetDataFilePath())
	if err != nil {
		log.Fatal("Error initializing new db writer", err)
	}

	defer dbw.Close()

	if err := dbw.Write(e); err != nil {
		log.Fatalf("Error writing expense to db: %v, %v", t.Cfg.GetDataFilePath(), err)
	}

	fmt.Printf("Added expense: %v\n", e)

	return nil
}

func (t *Tracker) String() string {
	return fmt.Sprintf("[ cfg: %v, entries: %v ]", t.Cfg, t.Entries)
}

func (t *Tracker) PrettyPrint() {
	for _, v := range t.Entries {
		fmt.Printf("[ type: %v, date: %v, amount: %v, category: %v, notes: %v ]", v.EntryType, v.Date, v.Amount, v.Category, v.Notes)
		fmt.Println()
	}
}
