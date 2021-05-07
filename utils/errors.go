package utils

import "log"

func Must(err error, message string) {
	if err != nil {
		log.Fatal(message)
	}
}
