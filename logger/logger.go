package logger

import (
	"log"
	"os"
)

var logger = log.New(os.Stderr, "", 0)

func Errorf(format string, v ...interface{}) {
	logger.Printf("ERROR: "+format, v)
}

func Fatalf(format string, v ...interface{}) {
	logger.Printf("FATAL: "+format, v)
}

func Infof(format string, v ...interface{}) {
	logger.Printf("INFO: "+format, v)
}
