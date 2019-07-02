package log

type Logger interface {
	Error(err error)
	Warn(msg string)
	Info(msg string)
}
