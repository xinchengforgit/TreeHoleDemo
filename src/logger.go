package main

import (
	"io"
	"log"
	"os"
)

// 先从easyDemo做起
// 此外得考虑如何保证用户加密的方式
func initLog() {
	logFile, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	// logger
	w := io.MultiWriter(logFile, os.Stdout)
	log.SetOutput(w)
}
