package logger

import (
	"log"
	"os"
)

var InfoLogger *log.Logger

func init() {
	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}
