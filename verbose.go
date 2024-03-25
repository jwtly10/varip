package main

import "log"

func verbose(message string, args ...interface{}) {
	if debugEnabled {
		log.Printf(message, args...)
	}
}
