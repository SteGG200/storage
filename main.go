package main

import (
	"flag"

	"github.com/SteGG200/storage/logger"
)

var (
	dbFilePath       = flag.String("database", "./data.db", "Path to your database file.")
	storageDirectory = flag.String("storage", "./storage", "Path to your storage directory.")
	logFilePath      = flag.String("log", "std", "Path to your log file. If not provided, the default is standard output.")
)

func main() {
	flag.Parse()

	logger.SetFile(*logFilePath)

	logger.InfoLogger.Println(*dbFilePath, *storageDirectory)
}
