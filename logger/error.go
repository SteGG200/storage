package logger

import (
	"log"
	"os"
)

var ErrorLogger *log.Logger

func init() {
	ErrorLogger = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
