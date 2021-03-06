//Copied since repository no longer exists
// Tideland Common Go Library - Application Log
//
// Copyright (C) 2012 Frank Mueller / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package applog

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"
)

//--------------------
// LOG LEVEL
//--------------------

// Log levels to control the logging output.
const (
	LevelDebug = iota
	LevelInfo
	LevelWarning
	LevelError
	LevelCritical
)

// logLevel controls the global log level used by the logger.
var level = LevelDebug

// LogLevel returns the global log level and can be used in
// own implementations of the logger interface.
func Level() int {
	return level
}

// SetLogLevel sets the global log level used by the simple
// logger.
func SetLevel(l int) {
	level = l
}

//--------------------
// LOGGER
//--------------------

// Logger is the interface for different logger implementations.
type Logger interface {
	// Debug logs a message at debug level.
	Debug(info, msg string)
	// Info logs a message at info level.
	Info(info, msg string)
	// Warning logs a message at warning level.
	Warning(info, msg string)
	// Error logs a message at error level.
	Error(info, msg string)
	// Critical logs a message at critical level.
	Critical(info, msg string)
}

// logger references the used application logger.
var logger Logger = NewStandardLogger(os.Stdout)

// SetLogger sets a new logger.
func SetLogger(l Logger) {
	logger = l
}

// timeFormat controls how the timestamp of the standard logger is printed.
const timeFormat = "2006-01-02 15:04:05 Z07:00"

// StandardLogger is a simple logger writing to the given writer. It
// doesn't handle the levels differently.
type StandardLogger struct {
	mutex sync.Mutex
	out   io.Writer
}

// NewStandardLogger creates the standard logger.
func NewStandardLogger(out io.Writer) Logger {
	return &StandardLogger{out: out}
}

// Debug logs a message at debug level.
func (sl *StandardLogger) Debug(info, msg string) {
	sl.mutex.Lock()
	defer sl.mutex.Unlock()

	io.WriteString(sl.out, "[D] ")
	io.WriteString(sl.out, time.Now().Format(timeFormat))
	io.WriteString(sl.out, " ")
	io.WriteString(sl.out, info)
	io.WriteString(sl.out, " ")
	io.WriteString(sl.out, msg)
	io.WriteString(sl.out, "\n")
}

// Info logs a message at info level.
func (sl *StandardLogger) Info(info, msg string) {
	sl.mutex.Lock()
	defer sl.mutex.Unlock()

	io.WriteString(sl.out, "[I] ")
	io.WriteString(sl.out, time.Now().Format(timeFormat))
	io.WriteString(sl.out, " ")
	io.WriteString(sl.out, info)
	io.WriteString(sl.out, " ")
	io.WriteString(sl.out, msg)
	io.WriteString(sl.out, "\n")
}

// Warning logs a message at warning level.
func (sl *StandardLogger) Warning(info, msg string) {
	sl.mutex.Lock()
	defer sl.mutex.Unlock()

	io.WriteString(sl.out, "[W] ")
	io.WriteString(sl.out, time.Now().Format(timeFormat))
	io.WriteString(sl.out, " ")
	io.WriteString(sl.out, info)
	io.WriteString(sl.out, " ")
	io.WriteString(sl.out, msg)
	io.WriteString(sl.out, "\n")
}

// Error logs a message at error level.
func (sl *StandardLogger) Error(info, msg string) {
	sl.mutex.Lock()
	defer sl.mutex.Unlock()

	io.WriteString(sl.out, "[E] ")
	io.WriteString(sl.out, time.Now().Format(timeFormat))
	io.WriteString(sl.out, " ")
	io.WriteString(sl.out, info)
	io.WriteString(sl.out, " ")
	io.WriteString(sl.out, msg)
	io.WriteString(sl.out, "\n")
}

// Critical logs a message at critical level.
func (sl *StandardLogger) Critical(info, msg string) {
	sl.mutex.Lock()
	defer sl.mutex.Unlock()

	io.WriteString(sl.out, "[C] ")
	io.WriteString(sl.out, time.Now().Format(timeFormat))
	io.WriteString(sl.out, " ")
	io.WriteString(sl.out, info)
	io.WriteString(sl.out, " ")
	io.WriteString(sl.out, msg)
	io.WriteString(sl.out, "\n")
}

// GoLogger just uses the standard go log package.
type GoLogger struct{}

// Debug logs a message at debug level.
func (gl GoLogger) Debug(info, msg string) {
	log.Println("[D]", info, msg)
}

// Info logs a message at info level.
func (gl GoLogger) Info(info, msg string) {
	log.Println("[I]", info, msg)
}

// Warning logs a message at warning level.
func (gl GoLogger) Warning(info, msg string) {
	log.Println("[W]", info, msg)
}

// Error logs a message at error level.
func (gl GoLogger) Error(info, msg string) {
	log.Println("[E]", info, msg)
}

// Critical logs a message at critical level.
func (gl GoLogger) Critical(info, msg string) {
	log.Println("[C]", info, msg)
}

//--------------------
// LOGGING
//--------------------

// Debugf logs a message at debug level.
func Debugf(format string, args ...interface{}) {
	if level <= LevelDebug {
		ci := retrieveCallInfo()
		fi := fmt.Sprintf(format, args...)

		logger.Debug(ci.verboseFormat(), fi)
	}
}

// Infof logs a message at info level.
func Infof(format string, args ...interface{}) {
	if level <= LevelInfo {
		ci := retrieveCallInfo()
		fi := fmt.Sprintf(format, args...)

		logger.Info(ci.shortFormat(), fi)
	}
}

// Warningf logs a message at warning level.
func Warningf(format string, args ...interface{}) {
	if level <= LevelWarning {
		ci := retrieveCallInfo()
		fi := fmt.Sprintf(format, args...)

		logger.Warning(ci.shortFormat(), fi)
	}
}

// Errorf logs a message at error level.
func Errorf(format string, args ...interface{}) {
	if level <= LevelError {
		ci := retrieveCallInfo()
		fi := fmt.Sprintf(format, args...)

		logger.Error(ci.shortFormat(), fi)
	}
}

// Criticalf logs a message at critical level.
func Criticalf(format string, args ...interface{}) {
	ci := retrieveCallInfo()
	fi := fmt.Sprintf(format, args...)

	logger.Critical(ci.verboseFormat(), fi)
}

//Implement io.Writer
type logWriter struct {
	level  int
	prefix string
}

func (f *logWriter) Write(p []byte) (int, error) {
	if level <= LevelDebug {
		Debugf("%s%s", f.prefix, p)
	} else if level <= LevelInfo {
		Infof("%s%s", f.prefix, p)
	} else if level <= LevelWarning {
		Warningf("%s%s", f.prefix, p)
	} else if level <= LevelError {
		Errorf("%s%s", f.prefix, p)
	} else if level <= LevelCritical {
		Criticalf("%s%s", f.prefix, p)
	}
	return len(p), nil
}

func NewLogWriter(level int, prefix string) io.Writer {
	return &logWriter{level: level, prefix: prefix}
}

//--------------------
// HELPER
//--------------------

// callInfo bundles the info about the call environment
// when a logging statement occured.
type callInfo struct {
	packageName string
	fileName    string
	funcName    string
	line        int
}

// shortFormat returns a string representation in a short variant.
func (ci *callInfo) shortFormat() string {
	return fmt.Sprintf("[%s]", ci.packageName)
}

// verboseFormat returns a string representation in a more verbose variant.
func (ci *callInfo) verboseFormat() string {
	return fmt.Sprintf("[%s] (%s:%s:%d)", ci.packageName, ci.fileName, ci.funcName, ci.line)
}

// retrieveCallInfo
func retrieveCallInfo() *callInfo {
	pc, file, line, _ := runtime.Caller(2)
	_, fileName := path.Split(file)
	parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	pl := len(parts)
	packageName := ""
	funcName := parts[pl-1]

	if parts[pl-2][0] == '(' {
		funcName = parts[pl-2] + "." + funcName
		packageName = strings.Join(parts[0:pl-2], ".")
	} else {
		packageName = strings.Join(parts[0:pl-1], ".")
	}

	return &callInfo{
		packageName: packageName,
		fileName:    fileName,
		funcName:    funcName,
		line:        line,
	}
}
