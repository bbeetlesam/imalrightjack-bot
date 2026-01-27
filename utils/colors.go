package utils

import "log"

const (
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
	ColorBold 	= "\033[1m"
	ColorReset  = "\033[0m"

	LogTypeInfo = "info"
	LogTypeUpdt = "updt"
	LogTypeWarn = "warn"
	LogTypeErrs = "errs"
	LogTypeDbwr = "dbwr"
)

func ColoriseLog(logType string, msg string) string {
	var prefix, color string

	switch logType {
	case LogTypeInfo:
		prefix = "INFO"
		color = ColorBlue
	case LogTypeUpdt:
		prefix = "UPDT"
		color = ColorCyan
	case LogTypeWarn:
		prefix = "WARN"
		color = ColorYellow
	case LogTypeErrs:
		prefix = "ERRS"
		color = ColorRed
	case LogTypeDbwr:
		prefix = "DBWR"
		color = ColorGreen
	default:
		return msg
	}

	return color + ColorBold + prefix + ColorReset + " " + msg
}

func LogColor(logType string, msg string) {
	prefix := ColoriseLog(logType, "")
	message := prefix + msg

	log.Println(message)
}

func LogColorf(logType string, format string, v ...any) {
	prefix := ColoriseLog(logType, "")
	message := prefix + format

	log.Printf(message, v...)
}
