package golog

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

type GoLogger struct {
	Mutex         sync.Mutex
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
	FatalLogger   *log.Logger
	DebugLogger   *log.Logger
}

func getStack(all bool) []string {
	buf := make([]byte, 1<<8)
	for {
		n := runtime.Stack(buf, all)
		if n < len(buf) {
			break
		}
		buf = make([]byte, len(buf)*2)
	}
	return strings.Split(string(buf), "\n")
}

func NewGoLogger(output *os.File) *GoLogger {
	logger := new(GoLogger)

	logger.InfoLogger = log.New(output, "[ INFO  ]", log.LstdFlags)
	logger.WarningLogger = log.New(output, "[WARNING]", log.LstdFlags)
	logger.ErrorLogger = log.New(output, "[ ERROR ]", log.LstdFlags)
	logger.FatalLogger = log.New(output, "[ FATAL ]", log.LstdFlags)
	logger.DebugLogger = log.New(output, "[ DEBUG ]", log.LstdFlags)

	return logger
}

func (l *GoLogger) Infoln(v ...interface{}) {
	l.Mutex.Lock()
	l.InfoLogger.Println(v...)
	l.Mutex.Unlock()
}

func (l *GoLogger) Infof(format string, v ...interface{}) {
	l.Infoln(fmt.Sprintf(format, v...))
}

func (l *GoLogger) Warningln(v ...interface{}) {
	l.Mutex.Lock()
	l.WarningLogger.Println(v...)
	l.Mutex.Unlock()
}

func (l *GoLogger) Warningf(format string, v ...interface{}) {
	l.Warningln(fmt.Sprintf(format, v...))
}

func (l *GoLogger) errorln(v ...interface{}) {
	_, file, line, _ := runtime.Caller(2)
	caller := fmt.Sprintf("%s:%d", filepath.Base(file), line)
	v = append([]interface{}{caller}, v...)
	l.Mutex.Lock()
	l.ErrorLogger.Println(v...)
	l.Mutex.Unlock()
}

func (l *GoLogger) Errorln(v ...interface{}) {
	l.errorln(v...)
}

func (l *GoLogger) Errorf(format string, v ...interface{}) {
	l.errorln(fmt.Sprintf(format, v...))
}

func (l *GoLogger) Fatalln(v ...interface{}) {
	l.Mutex.Lock()
	l.FatalLogger.Println(v...)
	for _, v := range getStack(false) {
		l.FatalLogger.Println(v)
	}
	l.Mutex.Unlock()
	os.Exit(1)
}

func (l *GoLogger) Fatalf(format string, v ...interface{}) {
	l.Fatalln(fmt.Sprintf(format, v...))
}

func (l *GoLogger) debugln(v ...interface{}) {
	_, file, line, _ := runtime.Caller(2)
	caller := fmt.Sprintf("%s:%d", filepath.Base(file), line)
	v = append([]interface{}{caller}, v...)
	l.Mutex.Lock()
	l.DebugLogger.Println(v...)
	l.Mutex.Unlock()
}

func (l *GoLogger) Debugln(debug bool, v ...interface{}) {
	if !debug {
		return
	}
	l.debugln(v...)
}

func (l *GoLogger) Debugf(debug bool, format string, v ...interface{}) {
	if !debug {
		return
	}
	l.debugln(fmt.Sprintf(format, v...))
}
