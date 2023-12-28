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
	m := map[string]Level{
		"debug": Debug,
		"info":  Info,
		"warn":  Warn,
		"error": Error,
	}
	if l, ok := m[s]; ok {
		return l, nil
	}
	return Error, errors.New("log level is invalid")
}
