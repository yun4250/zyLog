package zylog

import (
	"log"
	"os"
	"sync"
)

var debugStdout = false
var stdoutLogger *log.Logger
var soLock sync.RWMutex

func isStdOutEmpty() bool {
	soLock.RLock()
	defer soLock.RUnlock()
	return stdoutLogger == nil
}

func stdout() *log.Logger {
	if isStdOutEmpty() {
		soLock.Lock()
		defer soLock.Unlock()
		stdoutLogger = log.New(OpenOrCreate("stdout"), "", log.LstdFlags)
	}
	return stdoutLogger
}

func UseSystemOutput() {
	stdout().SetOutput(os.Stdout)
}

func OpenPrintDebug() {
	debugStdout = true
}

func ClosePrintDebug() {
	debugStdout = false
}

func Panic(format string, a ...interface{}) {
	stdout().Printf(format, a...)
	os.Exit(1)
}

func Print(format string, a ...interface{}) {
	stdout().Printf(format, a...)
}

func PrintDebug(format string, a ...interface{}) {
	if debugStdout {
		stdout().Printf(format, a...)
	}
}
