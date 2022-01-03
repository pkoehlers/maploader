package util

import (
	"log"
	"os"
)

func InitLogging() {
	err := RotateFile(7, "log/maploader", "log")
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.OpenFile("log/maploader.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	log.SetOutput(f)
	log.Println("Logging initialized")
}
