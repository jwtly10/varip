package main

import (
	"bufio"
	"encoding/json"
	"fmt"
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
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	jsonBlob, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var jsonData map[string]interface{}
	err = json.Unmarshal(jsonBlob, &jsonData)
	if err != nil {
		return nil, err
	}

	var matches []Match
	for k, v := range jsonData {
		if re.MatchString(k) {
			value, err := json.MarshalIndent(v, "", " ")
			if err != nil {
				return nil, err
			}

			matches = append(matches, Match{Path: filePath, Key: k, Value: string(value)})
		}
	}

	return matches, nil
}

func flattenJSON(prefix string, value interface{}, flattened map[string]interface{}) {
	switch child := value.(type) {
	case map[string]interface{}:
		for k, v := range child {
			flattenJSON(prefix+k+".", v, flattened)
		}
	case []interface{}:
		for i, v := range child {
			flattenJSON(fmt.Sprintf("%s[%d].", prefix, i), v, flattened)
		}
	default:
		flattened[prefix[:len(prefix)-1]] = value
	}
}

func ParseYAMLFile(filePath string, re *regexp.Regexp) ([]Match, error) {
	return nil, nil
}

func ParseTOMLFile(filePath string, re *regexp.Regexp) ([]Match, error) {
	return nil, nil
}
