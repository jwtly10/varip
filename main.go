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

var verboseEnabled bool = false
var showColor bool = true
var showErrors bool = false

var yellow = color.New(color.FgYellow)
var blue = color.New(color.FgBlue)
var red = color.New(color.FgRed)

func main() {
	app := &cli.App{
		Name:      "varip",
		Usage:     "Searches for environment variables in files. Searches in the current directory by default.",
		UsageText: "varip [options] [path] [pattern]",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "verbose",
				Usage: "Enable verbose debug logging",
			},
			&cli.BoolFlag{
				Name:  "errors",
				Usage: "Display errors in output, by default errors are hidden, so only matches are shown",
			},
			&cli.BoolFlag{
				Name:  "no-color",
				Usage: "Disable colorized output, useful if performance is slow or colors not supported by your terminal",
			},
			&cli.BoolFlag{
				Name:  "show-hidden",
				Usage: "Show hidden files and directories",
			},
		},
		Action: func(c *cli.Context) error {
			verboseEnabled = c.Bool("debug")
			showErrors = c.Bool("errors")
			showHidden := c.Bool("show-hidden")
			showColor = !c.Bool("no-color")

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

			coloredPrintf(yellow, "Searching for pattern '%s' in %s\n\n", pattern, fullPath)

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
			handleError(err, rootPath)
			verbose("Error walking directory %s: %s", rootPath, err)
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
		handleError(err, rootPath)
		verbose("Error walking directory %s: %s", rootPath, err)
	}

	return nil
}

// handleError prints the given error message
func handleError(err error, path string) {
	if showErrors {
		// TODO: Wrap errors for more context
		coloredPrintf(blue, "%s\n", path)
		coloredPrintf(red, "Error: %s\n\n", err)
	}
}

// isSupportedFileType checks if the file type of the given path is supported for searching.
func isSupportedFileType(path string) bool {

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

	coloredPrintf(blue, "%s\n", m[0].Path)
	highlight := color.New(color.FgHiRed).SprintFunc()

	for _, match := range m {
		y := color.New(color.FgHiYellow).SprintFunc()
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

// coloredPrintf prints the formatted string with or without color.
// The first argument is a color attribute (if color is enabled), followed by the format string and any additional arguments.
func coloredPrintf(colorAttribute *color.Color, format string, a ...interface{}) {
	formattedMessage := fmt.Sprintf(format, a...)

	if showColor && colorAttribute != nil {
		// If color is enabled, apply the color to the already formatted message
		fmt.Print(colorAttribute.Sprint(formattedMessage))
	} else {
		// If color is disabled, print the formatted message directly
		fmt.Print(formattedMessage)
	}
}
