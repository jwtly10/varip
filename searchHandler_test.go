package main

import (
	"fmt"
	"testing"
)

func TestSearchHandlerInvalidFile(t *testing.T) {
	h := NewSearchHandler()

	pat, err := generateRegex("pattern")
	if err != nil {
		t.Fatalf("Failed to compile regex pattern: %v", err)
	}

	err = h.Search("/path/to/invalidFile", pat, false)

	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	fmt.Print(err)
}
