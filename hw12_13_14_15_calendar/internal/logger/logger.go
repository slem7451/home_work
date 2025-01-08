package logger

import (
	"log"
	"strings"
)

const (
	errLevel   = "error"
	warnLevel  = "warn"
	infoLevel  = "info"
	debugLevel = "debug"
)

var logLevels = map[string]int{
	errLevel:   0,
	warnLevel:  1,
	infoLevel:  2,
	debugLevel: 3,
}

type Logger struct {
	level int
}

func New(level string) *Logger {
	numLevel, ok := logLevels[strings.ToLower(level)]
	if !ok {
		panic("Unknown log level")
	}

	return &Logger{level: numLevel}
}

func (l Logger) Info(msg string) {
	if l.level >= logLevels[infoLevel] {
		log.Println("INFO:", msg)
	}
}

func (l Logger) Error(msg string) {
	if l.level >= logLevels[errLevel] {
		log.Println("ERROR:", msg)
	}
}

func (l Logger) Warn(msg string) {
	if l.level >= logLevels[warnLevel] {
		log.Println("WARNING:", msg)
	}
}

func (l Logger) Debug(msg string) {
	if l.level >= logLevels[debugLevel] {
		log.Println("DEBUG:", msg)
	}
}
