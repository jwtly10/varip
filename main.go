package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

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
	h := NewSearchHandler()
	app := setupApp(h)

	err := app.Run(os.Args)
	if err != nil {
		coloredPrintf(red, "Error running varip: %s\n", err)
		os.Exit(1)
	}
}

type FileSearcher interface {
	Search(path string, pattern *regexp.Regexp, showHidden bool) error
}

func setupApp(searchHandler FileSearcher) *cli.App {
	app := &cli.App{
		Name:      "varip",
		Usage:     "Searches for environment variables in files. Searches in the current directory by default.",
		UsageText: "varip [options] [pattern] [path]",
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
				// If there are two or more arguments, first is pattern, second is path
				pattern = c.Args().Get(0)
				path = c.Args().Get(1)
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

			return searchHandler.Search(fullPath, regPattern, showHidden)
		},
	}
	return app
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

// generateRegex generates a case-insensitive regular expression from the given pattern.
func generateRegex(pattern string) (*regexp.Regexp, error) {
	pattern = "(?i)" + regexp.QuoteMeta(pattern)
	re, err := regexp.Compile(pattern)
	return re, err
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
