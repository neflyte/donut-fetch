package internal

import (
	"log"
	"os"
)

const (
	logFlags = log.LstdFlags | log.Lmicroseconds | log.Lshortfile | log.Lmsgprefix
)

func GetLogger(funcName string) *log.Logger {
	return log.New(os.Stderr, funcName+"] ", logFlags)
}
