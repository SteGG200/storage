package logger

import (
	"fmt"
	"log"
	"os"
)

/*
SetFile sets the log file to a specified path.
If the path is "std", it will log to standard output.
If the file cannot be opened, it will log to the console and ask the user if they want to continue.
If the user chooses not to, the program will exit.
Otherwise, the program will continue and it will log to standard output

Params:

	filePath string // The path to the log file.
*/
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
