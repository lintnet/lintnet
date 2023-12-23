package errlevel

import "errors"

type Level int

const (
	Debug Level = iota
	Info
	Warn
	Error
)

func New(s string) (Level, error) {
	switch s {
	case "debug":
		return Debug, nil
	case "info":
		return Info, nil
	case "warn":
		return Warn, nil
	case "error":
		return Error, nil
	}
	return Error, errors.New("log level is invalid")
}
