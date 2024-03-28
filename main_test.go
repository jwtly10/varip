package main

import (
	"fmt"
	"regexp"
	"testing"
)

type MockSearchHandler struct {
	Path       string
	Pattern    *regexp.Regexp
	ShowHidden bool
	CallCount  int
}

func (m *MockSearchHandler) Search(path string, pattern *regexp.Regexp, showHidden bool) error {
	m.Path = path
	m.Pattern = pattern
	m.ShowHidden = showHidden
	m.CallCount++
	return nil
}

func TestCLIParsingOrder(t *testing.T) {
	mock := &MockSearchHandler{}
	app := setupApp(mock)

	err := app.Run([]string{"varip", "pattern", "./path/testPath"})
	if err != nil {
		fmt.Print("Expected error")
	}

	expectedPath := abs("./path/testPath")
	expectedPattern := "(?i)" + "pattern"

	if mock.Path != expectedPath {
		t.Errorf("Expected path to be %s, got %s", expectedPath, mock.Path)
	}

	if mock.Pattern.String() != expectedPattern {
		t.Errorf("Expected pattern to be %s, got %s", expectedPattern, mock.Pattern.String())
	}
}
