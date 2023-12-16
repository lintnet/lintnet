package lint

import "errors"

type LogLevel int

const (
	debugLevel LogLevel = iota
	infoLevel
	warnLevel
	errorLevel
)

func newLogLevel(s string) (LogLevel, error) {
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
