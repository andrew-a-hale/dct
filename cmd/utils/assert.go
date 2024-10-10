package utils

import "log"

func Assert(b bool, msg string) {
	if !b {
		log.Fatalf("Assertion Error: %s\n", msg)
	}
}
