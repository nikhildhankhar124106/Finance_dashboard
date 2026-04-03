package logger

import (
	"log"
	"os"
)

// Logger is a basic wrapper for the standard library logger.
// In a production-ready application, you might want to replace this with
// a structured logger like zap or logrus.
var (
	Info  *log.Logger
	Error *log.Logger
)

func Init() {
	Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
