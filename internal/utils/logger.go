package utils

import (
	"log"
	"os"
)

var logger = log.New(os.Stdout, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile)

func LogInfo(msg string) {
	logger.Println("INFO:", msg)
}

func LogError(msg string) {
	logger.Println("ERROR:", msg)
}
