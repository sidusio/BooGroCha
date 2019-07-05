package diectory

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrNoServices        = Error("no booking services")
	ErrAllServicesFailed = Error("all booking services failed")
)
