package utils

import (
	"log"
	"os"
)

type Logger struct {
	*log.Logger
}

func NewLogger() *Logger {
	return &Logger{Logger: log.New(os.Stdout, "[microservice] ", log.LstdFlags|log.Lshortfile)}
}

func (l *Logger) AsyncLog(msg string) {
	go l.Println(msg)
}
