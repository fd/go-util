package log

import (
	"fmt"
	"io"
)

func New(w io.Writer, l Level) Logger {
	return &logger{w, l}
}

type logger struct {
	w     io.Writer
	level Level
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
