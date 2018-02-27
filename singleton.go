package zylog

import (
	"log"
	"os"
	"sync"
	"path/filepath"
)

var UseStdout = false
var DebugPrint = false
var stdoutLogger = log.New(os.Stdout, "", log.LstdFlags)
var soLock sync.Mutex
var ok = false

func stdout() *log.Logger {
	if !UseStdout && !ok {
		soLock.Lock()
		defer soLock.Unlock()
		absPath, _ := filepath.Abs("stdout")
		if file, e := os.OpenFile(absPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777); e == nil {
			stdoutLogger.SetOutput(file)
			ok = true
			if DebugPrint {
				stdoutLogger.Printf("OpenFile success: %s\n", absPath)
			}
		} else {
			stdoutLogger.Printf("OpenFile failed: %s\n", e.Error())
		}
	}
	return stdoutLogger
}

func Panic(format string, a ...interface{}) {
	stdout().Printf(format, a...)
	os.Exit(1)
}

func Print(format string, a ...interface{}) {
	stdout().Printf(format, a...)
}

func PrintDebug(format string, a ...interface{}) {
	if DebugPrint {
		stdout().Printf(format, a...)
	}
}
