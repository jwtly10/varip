package main

import (
	"testing"
)

type testCaseParseFile struct {
	filePath        string
	searchPattern   string
	expectedMatches []Match
}

func TestParseJsonFile(t *testing.T) {

	testCases := []testCaseParseFile{
		{
			filePath:      "./testdata/config.json",
			searchPattern: "database",
			expectedMatches: []Match{
				{Path: "./testdata/config.json", LineNum: 0, Key: "database", Value: "{\n \"host\": \"localhost\",\n \"password\": \"secret\",\n \"port\": 5432,\n \"user\": \"admin\"\n}"},
			},
		},
	}

	for _, testCase := range testCases {
		re, err := generateRegex(testCase.searchPattern)
		if err != nil {
			t.Fatalf("Failed to compile regex pattern: %v", err)
		}

		matches, err := ParseJSONFile(testCase.filePath, re)
		if err != nil {
			t.Fatalf("ParseJSONFile returned an error: %v", err)
		}

		if len(matches) != len(testCase.expectedMatches) {
			t.Errorf("Expected %d matches, got %d", len(testCase.expectedMatches), len(matches))
		}

		for i, match := range matches {
			if match != testCase.expectedMatches[i] {
				t.Errorf("Expected match %#v, got %#v at index %d", testCase.expectedMatches[i], match, i)
			}
		}

	}

}

func TestParsePropertiesFile(t *testing.T) {
	testCase := testCaseParseFile{
		filePath:      "./testdata/application.properties",
		searchPattern: "sprIng",
		expectedMatches: []Match{
			{Path: "./testdata/application.properties", LineNum: 22, Key: "spring.jpa.hibernate.ddl-auto", Value: "update"},
			{Path: "./testdata/application.properties", LineNum: 23, Key: "spring.datasource.url", Value: "jdbc:mysql://testdb.jfkasdfie3eu-west-2.rds.amazonaws.com/dev"},
			{Path: "./testdata/application.properties", LineNum: 24, Key: "spring.datasource.username", Value: "${SPRING_DATASOURCE_USERNAME}"},
			{Path: "./testdata/application.properties", LineNum: 25, Key: "spring.datasource.password", Value: "${SPRING_DATASOURCE_PASSWORD}"},
			{Path: "./testdata/application.properties", LineNum: 26, Key: "spring.datasource.driver-class-name", Value: "com.mysql.cj.jdbc.Driver"},
			{Path: "./testdata/application.properties", LineNum: 27, Key: "spring.datasource.hikari.maximum-pool-size", Value: "5"},
		},
	}

	re, err := generateRegex(testCase.searchPattern)
	if err != nil {
		t.Fatalf("Failed to compile regex pattern: %v", err)
	}

	matches, err := ParseEnvFile(testCase.filePath, re)
	if err != nil {
		t.Fatalf("ParseEnvFile returned an error: %v", err)
	}

	if len(matches) != len(testCase.expectedMatches) {
		t.Errorf("Expected %d matches, got %d", len(testCase.expectedMatches), len(matches))
	}

	for i, match := range matches {
		if match != testCase.expectedMatches[i] {
			t.Errorf("Expected match %#v, got %#v at index %d", testCase.expectedMatches[i], match, i)
		}
	}
}

func TestParseEnvFile(t *testing.T) {
	testCase := testCaseParseFile{
		filePath:      "./testdata/.env.local",
		searchPattern: "db",
		expectedMatches: []Match{
			{Path: "./testdata/.env.local", LineNum: 1, Key: "DB_HOST", Value: "localhost"},
			{Path: "./testdata/.env.local", LineNum: 2, Key: "DB_PORT", Value: "5432"},
		},
	}

	re, err := generateRegex(testCase.searchPattern)
	if err != nil {
		t.Fatalf("Failed to compile regex pattern: %v", err)
	}

	matches, err := ParseEnvFile(testCase.filePath, re)
	if err != nil {
		t.Fatalf("ParseEnvFile returned an error: %v", err)
	}

	if len(matches) != len(testCase.expectedMatches) {
		t.Errorf("Expected %d matches, got %d", len(testCase.expectedMatches), len(matches))
	}

	for i, match := range matches {
		if match != testCase.expectedMatches[i] {
			t.Errorf("Expected match %#v, got %#v at index %d", testCase.expectedMatches[i], match, i)
		}
	}
}
