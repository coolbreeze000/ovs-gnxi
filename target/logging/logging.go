package logging

import (
	"fmt"
	"os"
	"sync"

	"github.com/op/go-logging"
)

var (
	once     sync.Once
	instance map[string]*Logger
)

type Logger struct {
	tag           string
	consoleWriter *logging.Logger
}

func New(tag string) *Logger {
	fmt.Println("TEST")

	once.Do(func() {
		if instance == nil {
			instance = make(map[string]*Logger)
		}
	})

	if logger, ok := instance[tag]; ok {
		return logger
	}

	if _, ok := instance[tag]; !ok {
		instance[tag] = &Logger{tag: tag}
		logFormatConsole := logging.MustStringFormatter(
			`%{color:reset}%{color:bold}%{time:2006-01-02 15:04:05.000000} %{level} - %{shortfile} %{shortfunc}:%{color:reset}%{color} %{message} %{color:reset}`,
		)
		writer := logging.NewLogBackend(os.Stdout, "", 0)
		writerFormatter := logging.NewBackendFormatter(writer, logFormatConsole)
		writerLeveled := logging.AddModuleLevel(writerFormatter)
		writerLeveled.SetLevel(logging.DEBUG, tag)
		logging.SetBackend(writerLeveled)
		instance[tag].consoleWriter = logging.MustGetLogger(tag)
		instance[tag].consoleWriter.ExtraCalldepth = 1
	}

	return instance[tag]
}

func (l *Logger) Fatalf(format string, a ...interface{}) {
	l.consoleWriter.Fatalf(format, a...)
}

func (l *Logger) Fatal(message interface{}) {
	l.consoleWriter.Fatal(message)
}

func (l *Logger) Errorf(format string, a ...interface{}) {
	l.consoleWriter.Errorf(format, a...)
}

func (l *Logger) Error(message interface{}) {
	l.consoleWriter.Error(message)
}

func (l *Logger) Warningf(format string, a ...interface{}) {
	l.consoleWriter.Warningf(format, a...)
}

func (l *Logger) Warning(message interface{}) {
	l.consoleWriter.Warning(message)
}

func (l *Logger) Infof(format string, a ...interface{}) {
	l.consoleWriter.Infof(format, a...)
}

func (l *Logger) Info(message interface{}) {
	l.consoleWriter.Info(message)
}

func (l *Logger) Debugf(format string, a ...interface{}) {
	l.consoleWriter.Debugf(format, a...)
}

func (l *Logger) Debug(message interface{}) {
	l.consoleWriter.Debug(message)
}
