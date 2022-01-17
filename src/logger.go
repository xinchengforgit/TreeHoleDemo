package main

import (
	"io"
	"log"
	"os"
)


func InitLog() {
	logFile, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	// logger
	w := io.MultiWriter(logFile, os.Stdout)
	log.SetOutput(w)
}
