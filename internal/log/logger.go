package log

type Logger interface {
	Error(msg string)
	Errorv(msg string, v interface{})
	Errorf(format string, vs ...interface{})
	Warn(msg string)
	Warnv(msg string, v interface{})
	Warnf(format string, vs ...interface{})
	Info(msg string)
	Infov(msg string, v interface{})
	Infof(format string, vs ...interface{})
}
