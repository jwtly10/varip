package main

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

type Match struct {
	Path    string
	LineNum int
	Key     string
	Value   string
}

// ParseEnvFile parses an env file and returns all matches
// parses .env, .properties
func ParseEnvFile(filePath string, re *regexp.Regexp) ([]Match, error) {

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	var matches []Match

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++
		split := strings.SplitN(line, "=", 2)

		if len(split) == 2 {
			if re.MatchString(split[0]) {
				matches = append(matches, Match{Path: filePath, LineNum: lineNum, Key: split[0], Value: split[1]})
			}
		}

	}

	return matches, nil
}

func ParseJSONFile(filePath string, re *regexp.Regexp) ([]Match, error) {
	return nil, nil
}

func ParseYAMLFile(filePath string, re *regexp.Regexp) ([]Match, error) {
	return nil, nil
}

func ParseTOMLFile(filePath string, re *regexp.Regexp) ([]Match, error) {
	return nil, nil
}
