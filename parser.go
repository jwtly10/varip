package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type Match struct {
	Path    string
	LineNum int
	Key     string
	Value   string
}

// ParseEnvFile parses an env file and returns all matches.
// Parses .env, .properties
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

// ParseJSONFile parses a JSON file and returns all matches.
// Parses .json
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

	flattened := make(map[string]interface{})
	flattenJSON("", jsonData, flattened)

	var matches []Match
	for k, v := range flattened {
		if re.MatchString(k) {
			matches = append(matches, Match{Path: filePath, Key: k, Value: fmt.Sprintf("%v", v)})
		}

	}

	return matches, nil
}

// ParseYAMLFile parses the YAML file, flattens it, and finds matches based on the given regexp.
// Parses .yml, .yaml
func ParseYAMLFile(filePath string, re *regexp.Regexp) ([]Match, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var data interface{} // Use interface{} to handle arrays at the root
	if err = yaml.Unmarshal(file, &data); err != nil {
		return nil, err
	}

	flattened := make(map[string]string)
	flattenYAML("", data, flattened)

	var matches []Match
	for k, v := range flattened {
		if re.MatchString(k) {
			matches = append(matches, Match{Path: filePath, Key: k, Value: v})
		}
	}

	return matches, nil
}

// flattenYAML converts a nested YAML structure into a flat key-value map.

func flattenYAML(prefix string, value interface{}, flattened map[string]string) {
	switch child := value.(type) {
	case map[string]interface{}:
		for k, v := range child {
			flattenYAML(prefix+k+".", v, flattened)
		}
	case []interface{}:
		for i, v := range child {
			flattenYAML(fmt.Sprintf("%s[%d].", prefix, i), v, flattened)
		}
	case map[interface{}]interface{}: // YAML parsing might produce this
		for k, v := range child {
			flattenYAML(prefix+fmt.Sprintf("%v.", k), v, flattened)
		}
	default:
		flattened[prefix[:len(prefix)-1]] = fmt.Sprintf("%v", value) // Convert value to string
	}
}

// flattenJSON flattens a JSON object into a map.
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
