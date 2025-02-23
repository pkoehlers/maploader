package util

import (
	"log"
)

func CheckAndHandleError(err error) {
	if err != nil {
		log.Fatalf("[ERROR] %s: %v", err)
	}
}
