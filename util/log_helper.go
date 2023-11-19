package util

import (
	"io"
	"log"
	"os"
)

func InitLogging(logFolder string) {
	os.MkdirAll(logFolder, os.ModePerm)
	err := RotateFile(7, logFolder+"/maploader", "log")
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.OpenFile(logFolder+"/maploader.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	log.SetOutput(io.MultiWriter(f, os.Stdout))
	log.Println("Logging initialized")
}
