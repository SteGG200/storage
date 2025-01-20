package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/SteGG200/storage/logger"
	"github.com/SteGG200/storage/server"
	"github.com/SteGG200/storage/server/config"
	"github.com/joho/godotenv"
)

var (
	dbFilePath  = flag.String("database", "./data.db", "Path to your database file.")
	storagePath = flag.String("storage", "./storage", "Path to your storage directory.")
	logFilePath = flag.String("log", "std", "Path to your log file. If not provided, the default is standard output.")
)

func main() {
	flag.Parse()

	logger.SetFile(*logFilePath)

	if err := godotenv.Load(".env"); err != nil {
		logger.ErrorLogger.Fatal(err)
	}

	serverConfigs := []config.ConfigFunc{
		config.SetVerbose(),
		config.SetStoragePath(*storagePath),
	}

	server := server.New(dbFilePath, serverConfigs...)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	port, ok := os.LookupEnv("PORT")

	if !ok {
		port = "3000"
	}

	go server.Start(port)

	received_signal := <-stop

	logger.InfoLogger.Printf("Received signal %s. Shutting down...\n", received_signal)
}
