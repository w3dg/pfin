package db

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/w3dg/pfin/internal"
)

type JsonlReader struct {
	file    *os.File
	scanner *bufio.Scanner
}

type JsonlWriter struct {
	file   *os.File
	writer *json.Encoder
}

func NewReader(filename string) (*JsonlReader, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}

	return &JsonlReader{
		file:    f,
		scanner: bufio.NewScanner(f),
	}, nil
}

func (jr *JsonlReader) Close() error {
	return jr.file.Close()
}

func (jr *JsonlReader) ReadAll() ([]internal.Entry, error) {
	entries := make([]internal.Entry, 0, 128)

	for jr.scanner.Scan() {
		line := jr.scanner.Bytes()
		if len(line) == 0 {
			continue // skip blank lines
		}

		var entry internal.Entry
		if err := json.Unmarshal(line, &entry); err != nil {
			return nil, fmt.Errorf("parsing line: %w", err)
		}
		entries = append(entries, entry)
	}

	if err := jr.scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning file: %w", err)
	}

	return entries, nil
}

func NewWriter(filename string) (*JsonlWriter, error) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("could not open file for writing: %w", err)
	}

	return &JsonlWriter{
		file:   f,
		writer: json.NewEncoder(f),
	}, nil
}

func (jw *JsonlWriter) Write(entry internal.Entry) error {
	return jw.writer.Encode(entry)
}

func (jw *JsonlWriter) Close() error {
	return jw.file.Close()
}
