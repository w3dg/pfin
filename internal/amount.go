package internal

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Amount stores value in the smallest currency unit (e.g. paise/cents)
// to avoid floating point precision issues.
type Amount struct {
	Value int // e.g. 4550 represents 45.50
}

// NewAmount builds an Amount from whole and fractional (0-99) parts.
func NewAmount(whole, fraction int) Amount {
	return Amount{Value: whole*100 + fraction}
}

// NewAmountFromValue builds an Amount directly from the smallest unit.
func NewAmountFromValue(value int) Amount {
	return Amount{Value: value}
}

func (a Amount) Whole() int {
	return a.Value / 100
}

func (a Amount) Fraction() int {
	v := a.Value % 100
	if v < 0 {
		v = -v
	}
	return v
}

func (a Amount) String() string {
	if a.Value < 0 && a.Whole() == 0 {
		// handle -0.xx correctly since Whole() would print "0" without a sign
		return fmt.Sprintf("-0.%02d", a.Fraction())
	}
	return fmt.Sprintf("%d.%02d", a.Whole(), a.Fraction())
}

func (a Amount) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

// no negative amounts are allowed, the entry type handles that
func (a *Amount) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	parts := strings.SplitN(s, ".", 2)
	whole, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("invalid amount %q: %w", s, err)
	}

	fraction := 0
	if len(parts) == 2 {
		fracStr := parts[1]
		if len(fracStr) > 2 {
			fracStr = fracStr[:2] // truncate beyond cents
		} else if len(fracStr) == 1 {
			fracStr += "0" // pad e.g. ".5" -> ".50"
		}
		fraction, err = strconv.Atoi(fracStr)
		if err != nil {
			return fmt.Errorf("invalid amount %q: %w", s, err)
		}
	}

	value := whole*100 + fraction
	a.Value = value
	return nil
}
