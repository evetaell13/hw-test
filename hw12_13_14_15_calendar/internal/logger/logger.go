package logger

import (
	"fmt"
	"os"
	"strings"
)

const (
	levelError = iota
	levelWarn
	levelInfo
	levelDebug
)

type Logger struct { // TODO
	level int      // error / warn / info / debug
	file  *os.File // TODO close file
}

func New(level string) *Logger {
	var l int
	switch strings.ToLower(level) {
	case "error":
		l = levelError
	case "warn":
		l = levelWarn
	case "info":
		l = levelInfo
	case "debug":
		l = levelDebug
	default:
		return nil
	}
	f, err := os.CreateTemp("", "")
	if err != nil {
		return nil
	}
	return &Logger{
		level: l,
		file:  f,
	}
}

func (l Logger) Info(msg string) {
	fmt.Println(msg)
	l.log(msg, levelInfo)
}

func (l Logger) Error(msg string) {
	// TODO
	l.log(msg, levelError)
}

func (l Logger) Warn(msg string) {
	l.log(msg, levelWarn)
}

func (l Logger) Debug(msg string) {
	l.log(msg, levelDebug)
}

func (l Logger) log(msg string, level int) {
	if l.level >= level {
		fmt.Println(msg)
	}
}
