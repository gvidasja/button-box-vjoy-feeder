package log

import (
	"io"
	"log"
	"os"
)

var debug bool = false

func SetFile(path string) {
	logFile, _ := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModePerm)
	log.SetOutput(io.MultiWriter(logFile, os.Stdout))
}

func SetDebug(d bool) {
	debug = d
}

func Error(v ...interface{}) {
	logWithLevel(levelError, v...)
}

func Info(v ...interface{}) {
	logWithLevel(levelInfo, v...)
}

func Debug(v ...interface{}) {
	if debug {
		logWithLevel(levelDebug, v...)
	}
}

func Panic(v ...interface{}) {
	log.Panic(v...)
}

func logWithLevel(l level, v ...interface{}) {
	log.Println(append(l, v...)...)
}

type level []interface{}

var (
	levelDebug = level{"[DEBUG]"}
	levelInfo  = level{"[INFO]"}
	levelError = level{"[ERROR]"}
)
