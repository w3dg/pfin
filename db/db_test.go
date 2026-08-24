package db

import (
	"testing"
)

func TestReadAll(t *testing.T) {
	reader, err := NewReader("../testdata/db_test_ledger.jsonl")
	if err != nil {
		t.Fatalf("NewReader failed: %v", err)
	}
	defer reader.Close()

	entries, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll failed: %v", err)
	}

	if len(entries) != 4 {
		t.Errorf("expected 4 entries, got %d", len(entries))
	}
}
