package log

type Level uint8

const (
	DEFAULT = Level(0)

	DEBUG Level = 1 << iota
	INFO
	NOTICE
	WARN
	ERROR
	FATAL
	UNKNOWN
)

type Logger interface {
	SetLevel(lvl Level)

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

	Namespace() string
	Sub(level Level, namespace string) Logger
}

var DefaultLogger = New(nil, INFO, "")

func SetLevel(lvl Level) {
	DefaultLogger.SetLevel(lvl)
}

func Debug(args ...interface{}) {
	DefaultLogger.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	DefaultLogger.Debugf(format, args...)
}

func Info(args ...interface{}) {
	DefaultLogger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	DefaultLogger.Infof(format, args...)
}

func Notice(args ...interface{}) {
	DefaultLogger.Notice(args...)
}

func Noticef(format string, args ...interface{}) {
	DefaultLogger.Noticef(format, args...)
}

func Warn(args ...interface{}) {
	DefaultLogger.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	DefaultLogger.Warnf(format, args...)
}

func Error(args ...interface{}) {
	DefaultLogger.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	DefaultLogger.Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	DefaultLogger.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	DefaultLogger.Fatalf(format, args...)
}

func Namespace() string {
	return DefaultLogger.Namespace()
}

func Sub(level Level, namespace string) Logger {
	return DefaultLogger.Sub(level, namespace)
}
