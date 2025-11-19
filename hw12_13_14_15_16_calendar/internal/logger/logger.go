package logger

import (
	"fmt"
	"time"
)

type logLevel int

const (
	DEBUG logLevel = iota
	INFO
	WARN
	ERROR
)

var levelMap = map[string]logLevel{
	"DEBUG": DEBUG,
	"INFO":  INFO,
	"WARN":  WARN,
	"ERROR": ERROR,
}

type Logger struct {
	level logLevel
}

type LogDeps struct {
	Level string
}

func New(deps LogDeps) *Logger {
	return &Logger{
		level: levelMap[deps.Level],
	}
}

func (l Logger) Info(msg string) {
	msg = fmt.Sprintf("INFO: %s: %s", time.Now().Format("02.01.2006 15:04"), msg)
	l.sendMsg(INFO, msg)
}

func (l Logger) Error(msg string) {
	msg = fmt.Sprintf("ERROR: %s: %s", time.Now().Format("02.01.2006 15:04"), msg)
	l.sendMsg(ERROR, msg)
}

func (l Logger) Warn(msg string) {
	msg = fmt.Sprintf("WARN: %s: %s", time.Now().Format("02.01.2006 15:04"), msg)
	l.sendMsg(WARN, msg)
}

func (l Logger) Debug(msg string) {
	msg = fmt.Sprintf("DEBUG: %s: %s", time.Now().Format("02.01.2006 15:04"), msg)
	l.sendMsg(DEBUG, msg)
}

func (l Logger) sendMsg(lvl logLevel, msg string) {
	if lvl < l.level {
		return
	}
	fmt.Println(msg)
}
