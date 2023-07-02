package main

import (
	"log"
	"os"
	"sync"
	"time"
)

type Logger struct {
    filename string
    *log.Logger
}

var logger *Logger
var once sync.Once

func getLogger() *Logger {
    once.Do(func() {
        logger = createLogger("logs.txt")
    })
    return logger
}

func createLogger(fname string) *Logger {
    file, _ := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)

    return &Logger{
        filename: fname,
        Logger:   log.New(file, time.Now().Format("2006-01-02 15:04:05") + ", ", log.Lshortfile),
    }
}