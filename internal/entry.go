package internal

import (
	"encoding/json"
	"fmt"
	"strings"
)

type EntryType int

const (
	INCOME EntryType = iota
	SET_ASIDE
	TOP_UP
	EXPENSE
)

type Entry struct {
	EntryType EntryType    `json:"type"`
	Date      string       `json:"date"` // or time.Time
	Amount    Amount       `json:"amount"`
	Category  CategoryType `json:"category"`
	Notes     string       `json:"note"`
}

var EntryTypeName = map[EntryType]string{
	INCOME:    "Income",
	SET_ASIDE: "Aside",
	TOP_UP:    "Topup",
	EXPENSE:   "Expense",
}

func (et EntryType) String() string {
	etv, ok := EntryTypeName[et]
	if !ok {
		return "(entry type not found)"
	}
	return etv
}

var nameToEntryType = buildReverseEntryTypeMap()

func buildReverseEntryTypeMap() map[string]EntryType {
	m := make(map[string]EntryType, len(EntryTypeName))
	for k, v := range EntryTypeName {
		m[strings.ToLower(v)] = k
	}
	return m
}

func (e EntryType) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

func (e *EntryType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	et, ok := nameToEntryType[strings.ToLower(s)]
	if !ok {
		return fmt.Errorf("unknown entry type: %q", s)
	}
	*e = et
	return nil
}
