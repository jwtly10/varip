package main

import (
	"log"
	"testing"
)

type testCaseParseFile struct {
	filePath        string
	searchPattern   string
	expectedMatches []Match
}

// app:
//
//	name: "app"
//	deployment:
//	  image:
//	    tag: latest
//	    repository: "app"
//	    pullPolicy: Always
//	  replicaCount: 2
//	  JAVA_OPTS: "-Xms256m -Xmx512m"
//	  resources:
//	    cpuRequest: 200m
//	    memoryRequest: 512Mi
//	    memoryLimit: 768Mi
//	  environmentVariables:
//	    - name: API_URL
//	      value: https://api.gatewayurl.com/api/gateway/
//	    - name: SPRING_REDIS_HOST
//	      value: redis.host.com
//	    - name: SPRING_REDIS_PORT
//	      value: 6379
//	    - name: SPRING_REDIS_PASSWORD
//	      value: ${REDIS_PASSWORD}
//	    - name: SPRING_AUTOCONFIGURE_EXCLUDE
//	      value: "org.springframework.boot.autoconfigure.thymeleaf.ThymeleafAutoConfiguration"
//	    - name: API_PAYMENT_EXPIRATION_PERIOD_MINUTES
//	      value: 10
func TestParseYamlFile(t *testing.T) {
	testCases := []testCaseParseFile{
		{
			filePath:      "./testdata/prod.yaml",
			searchPattern: "resources",
			expectedMatches: []Match{
				{Path: "./testdata/prod.yaml", LineNum: 0, Key: "app.deployment.resources.cpuRequest", Value: "200m"},
				{Path: "./testdata/prod.yaml", LineNum: 0, Key: "app.deployment.resources.memoryLimit", Value: "768Mi"},
				{Path: "./testdata/prod.yaml", LineNum: 0, Key: "app.deployment.resources.memoryRequest", Value: "512Mi"},
			},
		},
	}

	for _, testCase := range testCases {
		re, err := generateRegex(testCase.searchPattern)
		if err != nil {
			t.Fatalf("Failed to compile regex pattern: %v", err)
		}

		matches, err := ParseYAMLFile(testCase.filePath, re)
		if err != nil {
			t.Fatalf("ParseJSONFile returned an error: %v", err)
		}

		log.Println(matches)

		if len(matches) != len(testCase.expectedMatches) {
			t.Errorf("Expected %d matches, got %d", len(testCase.expectedMatches), len(matches))
		}

		for _, match := range matches {
			if !contains(testCase.expectedMatches, match) {
				t.Errorf("Expected match %v in expected matches", match)
			}
		}

	}

}

func TestParseJsonFile(t *testing.T) {
	testCases := []testCaseParseFile{
		{
			filePath:      "./testdata/config.json",
			searchPattern: "database",
			expectedMatches: []Match{
				{Path: "./testdata/config.json", LineNum: 0, Key: "database.host", Value: "localhost"},
				{Path: "./testdata/config.json", LineNum: 0, Key: "database.password", Value: "secret"},
				{Path: "./testdata/config.json", LineNum: 0, Key: "database.port", Value: "5432"},
				{Path: "./testdata/config.json", LineNum: 0, Key: "database.user", Value: "admin"},
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

		log.Println(matches)

		if len(matches) != len(testCase.expectedMatches) {
			t.Errorf("Expected %d matches, got %d", len(testCase.expectedMatches), len(matches))
		}

		for _, match := range matches {
			if !contains(testCase.expectedMatches, match) {
				t.Errorf("Expected match %v in expected matches", match)
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

	log.Println(matches)

	if len(matches) != len(testCase.expectedMatches) {
		t.Errorf("Expected %d matches, got %d", len(testCase.expectedMatches), len(matches))
	}

	for i, match := range matches {
		if match != testCase.expectedMatches[i] {
			t.Errorf("Expected match %#v, got %#v at index %d", testCase.expectedMatches[i], match, i)
		}
	}
}

func contains(matches []Match, match Match) bool {
	for _, m := range matches {
		if m == match {
			return true
		}
	}
	return false
}
