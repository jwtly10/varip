package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

var debugEnabled bool = false

func main() {
	app := &cli.App{
		Name:      "varip",
		Usage:     "Searches for environment variables in files. Searches in the current directory by default.",
		UsageText: "varip [options] [path] [pattern]",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Enable verbose debug logging",
			},
			&cli.BoolFlag{
				Name:  "no-color",
				Usage: "Disable colorized output, useful if performance is slow",
			},
			&cli.BoolFlag{
				Name:  "show-hidden",
				Usage: "Show hidden files and directories",
			},
		},
		Action: func(c *cli.Context) error {
			debugEnabled = c.Bool("debug")
			showHidden := c.Bool("show-hidden")

			path := "."
			pattern := ""
			if c.NArg() > 1 {
				// If there are two or more arguments, first is path, second is pattern
				path = c.Args().Get(0)
				pattern = c.Args().Get(1)
			} else if c.NArg() == 1 {
				// If there is only one argument, it's the pattern
				pattern = c.Args().Get(0)
			} else {
				// No arguments provided
				cli.ShowAppHelpAndExit(c, 1)
			}

			fullPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}

			regPattern, err := generateRegex(pattern)
			if err != nil {
				log.Fatalf("Error generating regex: %s", err)
			}

			color.Yellow("Searching for pattern '%s' in %s\n", pattern, fullPath)

			return search(fullPath, regPattern, showHidden)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running varip: %s\n", err)
		os.Exit(1)
	}
}

// search recursively searches for the given pattern in files under the rootPath directory.
func search(rootPath string, pattern *regexp.Regexp, showHidden bool) error {
	err := filepath.WalkDir(rootPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			handleError(err)
			verbose("Error walking directory %s: %s", rootPath, err)
		}

		// Skip hidden files and directories, unless showHidden is true or it's a .env file
		if !showHidden && (strings.HasPrefix(d.Name(), ".") && !strings.Contains(d.Name(), ".env")) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if !d.IsDir() {
			if isSupportedFileType(path) {
				err := parseFile(path, pattern)
				if err != nil {
					handleError(err)
					verbose("Error searching in file %s: %s", path, err)
				}
			}
		}

		return nil
	})
	if err != nil {
		handleError(err)
		verbose("Error walking directory %s: %s", rootPath, err)
	}

	return nil
}

// handleError prints the given error message
func handleError(err error) {
	// TODO: Wrap errors for more context
	color.Red("Error: %s", err)
}

// isSupportedFileType checks if the file type of the given path is supported for searching.
func isSupportedFileType(path string) bool {
	supportedFileTypes := []string{".env*", "*.json", "*.properties", "*.yml", "*.yaml"}

	for _, fileType := range supportedFileTypes {
		fileName := filepath.Base(path)
		if match, _ := filepath.Match(fileType, fileName); match {
			return true
		}
	}

	verbose("Unsupported file type %s", path)
	return false
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

// generateRegex generates a case-insensitive regular expression from the given pattern.
func generateRegex(pattern string) (*regexp.Regexp, error) {
	pattern = "(?i)" + regexp.QuoteMeta(pattern)
	re, err := regexp.Compile(pattern)
	return re, err
}

// printMatches pretty prints the matches found in the given list of Match objects.
// Highlights the matched pattern in the key.
func printMatches(m []Match, re *regexp.Regexp) {
	if len(m) == 0 {
		return
	}

	color.Blue("%s", m[0].Path)

	highlight := color.New(color.FgHiRed).SprintFunc()

	for _, match := range m {
		highlightedKey := re.ReplaceAllStringFunc(match.Key, func(s string) string {
			return highlight(s)
		})
		// highlightedValue := re.ReplaceAllStringFunc(match.Value, func(s string) string {
		// 	return highlight(s)
		// })

		color.White("%d: %s => %s", match.LineNum, highlightedKey, match.Value)
	}
}
