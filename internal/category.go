package internal

import (
	"encoding/json"
	"fmt"
	"strings"
)

type CategoryType int

const (
	INCOME_CATEGORY CategoryType = iota
	SET_ASIDE_CATEGORY
	TOP_UP_CATEGORY
	FOOD
	TRANSPORT
	SUBSCRIPTIONS
	RENT
	GROCERIES
	UTILITIES
	SHOPPING
	ENTERTAINMENT
	MEDICAL
	FITNESS
	FAMILY
	GIFT
	PERSONAL_CARE
	OTHER
)

var CategoryName = map[CategoryType]string{
	INCOME_CATEGORY:    "Income",
	SET_ASIDE_CATEGORY: "Aside",
	TOP_UP_CATEGORY:    "Topup",
	FOOD:               "Food",
	TRANSPORT:          "Transport",
	SUBSCRIPTIONS:      "Subscriptions",
	RENT:               "Rent",
	GROCERIES:          "Groceries",
	UTILITIES:          "Utilities",
	SHOPPING:           "Shopping",
	ENTERTAINMENT:      "Entertainment",
	MEDICAL:            "Medical",
	FITNESS:            "Fitness",
	FAMILY:             "Family",
	GIFT:               "Gift",
	PERSONAL_CARE:      "Personal",
	OTHER:              "Other",
}

var NameToCategory = buildReverseCategoryMap()

func buildReverseCategoryMap() map[string]CategoryType {
	m := make(map[string]CategoryType, len(CategoryName))
	for k, v := range CategoryName {
		m[strings.ToLower(v)] = k
	}
	return m
}

func (c CategoryType) String() string {
	name := CategoryName[c]
	if name == "" {
		return "(name for category type not found)"
	}

	return name
}

func (c CategoryType) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func (c *CategoryType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	ct, ok := NameToCategory[strings.ToLower(s)]
	if !ok {

		return fmt.Errorf("unknown category type: %q", s)

	}
	*c = ct
	return nil
}
