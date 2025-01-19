package logger

import (
	"fmt"
	"log"
	"os"
)

func SetFile(filePath string) {
	if filePath == "std" {
		return
	}

	logFile, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.SetPrefix("WARNING: ")
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		log.Printf("Cannot open log file: %s. Do you want to continue? (yes | no)", filePath)
		var response string
		fmt.Scanln(&response)
		if response != "yes" {
			os.Exit(1)
		}
	} else {
		InfoLogger.SetOutput(logFile)
		WarningLogger.SetOutput(logFile)
		ErrorLogger.SetOutput(logFile)
	}
}
