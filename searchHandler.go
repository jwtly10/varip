package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

type SearchHandler struct {
}

func NewSearchHandler() *SearchHandler {
	return &SearchHandler{}
}

func (handler *SearchHandler) Search(path string, pattern *regexp.Regexp, showHidden bool) error {
	if !exists(path) {
		return fmt.Errorf("file or directory %s does not exist", path)
	}

	err := filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			handleError(err, path)
			verbose("Error walking directory %s: %s", path, err)
		}

		// Check if the entry is a symlink and skip it
		if d.Type().IsRegular() && d.Type()&os.ModeSymlink != 0 {
			verbose("Skipping symlink %s", path)
			return nil
		}

		if !showHidden {
			// If the entry is a hidden file or directory, skip it
			// Other than .env files, which are supported
			if strings.HasPrefix(filepath.Base(path), ".") && !strings.Contains(d.Name(), ".env") {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}

			// If the entry is a directory, check if it's in the ignore list
			for _, ignoredDir := range ignoredDirectories {
				if strings.Contains(path, ignoredDir) {
					if d.IsDir() {
						return filepath.SkipDir
					}
					return nil
				}
			}
		}

		if !d.IsDir() && isSupportedFileType(path) {
			err := parseFile(path, pattern)
			if err != nil {
				handleError(err, path)
				verbose("Error searching in file %s: %s", path, err)
			}
		}

		return nil
	})
	if err != nil {
		handleError(err, path)
		verbose("Error walking directory %s: %s", path, err)
	}

	return nil
}

// parseFile parses the file at the given path based on its file type and searches for the pattern.
func parseFile(path string, pattern *regexp.Regexp) error {
	var results []Match
	var err error

	switch {
	case strings.HasPrefix(filepath.Base(path), ".env") || strings.HasSuffix(path, ".properties"):
		results, err = ParseEnvFile(path, pattern)
	case strings.HasSuffix(path, ".json"):
		results, err = ParseJSONFile(path, pattern)
	case strings.HasSuffix(path, ".yml") || strings.HasSuffix(path, ".yaml"):
		results, err = ParseYAMLFile(path, pattern)

	default:
		err = fmt.Errorf("unsupported file type %s", path)
	}

	if err != nil {
		return err
	}

	printMatches(results, pattern)

	return nil
}

// printMatches pretty prints the matches found in the given list of Match objects.
// Highlights the matched pattern in the key.
func printMatches(m []Match, re *regexp.Regexp) {
	if len(m) == 0 {
		return
	}

	coloredPrintf(blue, "%s\n", m[0].Path)
	highlight := color.New(color.FgHiRed).SprintFunc()

	for _, match := range m {
		y := color.New(color.Faint).SprintFunc()
		highlightedLineNum := y(match.LineNum)

		highlightedKey := re.ReplaceAllStringFunc(match.Key, func(s string) string {
			return highlight(s)
		})

		f := color.New(color.Faint).SprintFunc()
		highlightedValue := f(match.Value)

		if showColor {
			// Line number is not supported for json parsing
			if match.LineNum == 0 {
				color.White("%s => %s", highlightedKey, highlightedValue)
				continue
			}
			color.White("%s: %s => %s", highlightedLineNum, highlightedKey, highlightedValue)
		} else {
			if match.LineNum == 0 {
				fmt.Printf("%s => %s\n", match.Key, match.Value)
				continue
			}
			fmt.Printf("%d: %s => %s\n", match.LineNum, match.Key, match.Value)
		}
	}

	fmt.Println()
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}
