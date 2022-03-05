package utils

import (
	"log"
	"os"
)

func ExitOnError(message string, err error) {
	if err != nil {
		log.Printf("[%s] - %s", message, err)
		os.Exit(-1)
	}
}

func PrintOnError(message string, err error) {
	if err != nil {
		log.Printf("[%s] - %s", message, err)
	}
}
