package fmt

import "fmt"

type Logger struct{}

const (
	errorTag   = "[ERROR]"
	warningTag = "[WARNING]"
	infoTag    = "[INFO]"
)

func (l *Logger) Error(msg string) {
	l.log(errorTag, msg)
}

func (l *Logger) Errorv(msg string, v interface{}) {
	l.logv(errorTag, msg, v)
}

func (l *Logger) Errorf(format string, vs ...interface{}) {
	l.logf(errorTag, format, vs...)
}

func (l *Logger) Warn(msg string) {
	l.log(warningTag, msg)
}

func (l *Logger) Warnv(msg string, v interface{}) {
	l.logv(warningTag, msg, v)
}

func (l *Logger) Warnf(format string, vs ...interface{}) {
	l.logf(warningTag, format, vs...)
}

func (l *Logger) Info(msg string) {
	l.log(infoTag, msg)
}

func (l *Logger) Infov(msg string, v interface{}) {
	l.logv(infoTag, msg, v)
}

func (l *Logger) Infof(format string, vs ...interface{}) {
	l.logf(infoTag, format, vs...)
}

func (*Logger) log(tag, msg string) {
	fmt.Printf("%s %s\n", tag, msg)
}

func (*Logger) logv(tag, msg string, v interface{}) {
	fmt.Printf("%s %s: %v\n", tag, msg, v)
}

func (*Logger) logf(tag, format string, vs ...interface{}) {
	fmt.Printf("%s %s", tag, fmt.Sprintf(format, vs...))
}
