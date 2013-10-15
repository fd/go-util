package log

import (
	"fmt"
	"io"
	"strings"
	"time"
)

type Level uint8

const (
	DEBUG Level = 1 << iota
	INFO
	NOTICE
	WARN
	ERROR
	FATAL
	UNKNOWN
)

const timeFormat = "2006-01-02 15:04:05.000"

var level_codes = map[Level]byte{
	DEBUG:   'D',
	INFO:    'I',
	NOTICE:  'N',
	WARN:    'W',
	ERROR:   'E',
	FATAL:   'F',
	UNKNOWN: 'U',
}

type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Notice(args ...interface{})
	Noticef(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
}

type logger struct {
	w     io.Writer
	level Level
}

func New(w io.Writer, l Level) Logger {
	return &logger{w, l}
}

func (l *logger) Debug(args ...interface{}) {
	l.emit(DEBUG, args...)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.emitf(DEBUG, format, args...)
}

func (l *logger) Info(args ...interface{}) {
	l.emit(INFO, args...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.emitf(INFO, format, args...)
}

func (l *logger) Notice(args ...interface{}) {
	l.emit(NOTICE, args...)
}

func (l *logger) Noticef(format string, args ...interface{}) {
	l.emitf(NOTICE, format, args...)
}

func (l *logger) Warn(args ...interface{}) {
	l.emit(WARN, args...)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.emitf(WARN, format, args...)
}

func (l *logger) Error(args ...interface{}) {
	l.emit(ERROR, args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.emitf(ERROR, format, args...)
}

func (l *logger) Fatal(args ...interface{}) {
	l.emit(FATAL, args...)
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	l.emitf(FATAL, format, args...)
}

func (l *logger) format(level Level, s string) string {
	// format the time of the message
	now := time.Now().UTC().Format(timeFormat)

	// clean up whitespace, indent subsequent lines
	s = strings.TrimSpace(strings.Replace(s, "\n", "\n  ", -1))

	// get the level code
	code := level_codes[level]
	if code == 0 {
		code = level_codes[UNKNOWN]
	}

	// format the message
	s = fmt.Sprintf("%s [%c] %s\n", now, code, s)

	return s
}

func (l *logger) emit(level Level, args ...interface{}) {
	if level >= l.level {
		entry := l.format(level, fmt.Sprint(args...))
		l.w.Write([]byte(entry))
	}
}

func (l *logger) emitf(level Level, format string, args ...interface{}) {
	if level >= l.level {
		entry := l.format(level, fmt.Sprintf(format, args...))
		l.w.Write([]byte(entry))
	}
}
