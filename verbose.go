package main

import "log"

func verbose(message string, args ...interface{}) {
	if verboseEnabled {
		log.Printf(message, args...)
	}
}
