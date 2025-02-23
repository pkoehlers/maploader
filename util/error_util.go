package util

import (
	"log"
)

func CheckAndHandleError(context string, err error) {
	if err != nil {
		log.Fatalf("[ERROR] %s: %v", context, err)
	}
}
