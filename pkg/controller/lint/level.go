package lint

import "errors"

type ErrorLevel int

const (
	debugLevel ErrorLevel = iota
	infoLevel
	warnLevel
	errorLevel
)

func newErrorLevel(s string) (ErrorLevel, error) {
	switch s {
	case "debug":
		return debugLevel, nil
	case "info":
		return infoLevel, nil
	case "warn":
		return warnLevel, nil
	case "error":
		return errorLevel, nil
	}
	return errorLevel, errors.New("log level is invalid")
}
