package util

import (
	"log"
	"os"
)

func CheckAndHandleError(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
