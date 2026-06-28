package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

func NewLogger(level zerolog.Level) zerolog.Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout}
	logger := zerolog.New(output).
		Level(level).
		With().
		Timestamp().
		Logger()

	return logger
}

func NewDebugLogger() zerolog.Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout}
	logger := zerolog.New(output).
		Level(zerolog.DebugLevel).
		With().
		Timestamp().
		Logger()

	return logger
}

func SetupLogger() {
	output := zerolog.ConsoleWriter{Out: os.Stdout}
	output.FormatLevel = func(i interface{}) string {
		var levelColor int
		var level string

		if ll, ok := i.(string); ok {
			const (
				colorBlack = iota + 30
				colorRed
				colorGreen
				colorYellow
				colorBlue
				colorMagenta
				colorCyan
				colorWhite

				colorBold     = 1
				colorDarkGray = 90
			)
			level = ll
			switch ll {
			case "trace":
				levelColor = colorMagenta
			case "debug":
				levelColor = colorBlue // Changed from default yellow to blue
			case "info":
				levelColor = colorGreen
			case "warn":
				levelColor = colorYellow
			case "error":
				levelColor = colorRed
			case "fatal":
				levelColor = colorRed | colorBold // Bold red for fatal
			case "panic":
				levelColor = colorRed | colorBold // Bold red for panic
			default:
				levelColor = colorWhite
			}
		}

		// Format with color
		if levelColor != 0 {
			return fmt.Sprintf("\x1b[%dm%s\x1b[0m", levelColor, strings.ToUpper(level))
		}
		return strings.ToUpper(level)
	}
}
