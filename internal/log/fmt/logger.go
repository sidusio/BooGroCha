package fmt

import "fmt"

type Logger struct{}

func (*Logger) Error(err error) {
	fmt.Println(err)
}

func (*Logger) Warn(msg string) {
	fmt.Println(msg)
}

func (*Logger) Info(msg string) {
	fmt.Println(msg)
}
