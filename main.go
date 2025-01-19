package main

import (
	"flag"

	"github.com/SteGG200/storage/logger"
)

var (
	dbFilePath       = flag.String("database", "./data.db", "Path to your database file.")
	storageDirectory = flag.String("storage", "./storage", "Path to your storage directory.")
	logFilePath      = flag.String("log", "./.log", "Path to your log file.")
)

func main() {
	flag.Parse()

	logger.SetFile(*logFilePath)

	logger.ErrorLogger.Fatal(*dbFilePath, *storageDirectory)
}
