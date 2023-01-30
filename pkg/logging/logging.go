package logging

import "github.com/sirupsen/logrus"

type Level int

const (
	Trace Level = iota
	Debug
	Info
	Warn
	Error
	Fatal
	Panic
	max
)

func VerbosityToLogrusLevel(l int) logrus.Level {
	return logrus.Level(int(max-1) - l)
}
