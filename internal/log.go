package internal

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var Logger InternalLogger = InternalLogger{
	logging: int(LoggingInfo),
}

type Logging int

const (
	LoggingNone  Logging = 0
	LoggingInfo  Logging = 1
	LoggingError Logging = 2
	LoggingWarn  Logging = 3
	LoggingDebug Logging = 4
)

type InternalLogger struct {
	logging int
}

var flagVerbose bool
var flagQuiet bool

func (logger *InternalLogger) FlagDeclare(cmd *cobra.Command, args ...interface{}) {
	cmd.Flags().BoolVar(&flagVerbose, "verbose", false, "")
	cmd.Flags().BoolVar(&flagQuiet, "quit", false, "")
}

func (logger *InternalLogger) FlagInit(args ...interface{}) {
	if flagVerbose {
		logger.logging = int(LoggingDebug)
	} else if flagQuiet {
		logger.logging = int(LoggingNone)
	}
}

func (logger *InternalLogger) ErrorObj(err error) {
	color.Set(color.FgRed)
	fmt.Println("error: ", err.Error())
	color.Unset()
}

func (logger *InternalLogger) Info(args ...interface{}) {
	if logger.logging >= int(LoggingInfo) {
		fmt.Println(args...)
	}
}

func (logger *InternalLogger) Error(args ...interface{}) {
	if logger.logging >= int(LoggingError) {
		color.Set(color.FgRed)
		fmt.Println(args...)
		color.Unset()
	}
}

func (logger *InternalLogger) Warn(args ...interface{}) {
	if logger.logging >= int(LoggingWarn) {
		color.Set(color.FgYellow)
		fmt.Println(args...)
		color.Unset()
	}
}

func (logger *InternalLogger) Debug(args ...interface{}) {
	if logger.logging >= int(LoggingDebug) {
		color.Set(color.FgMagenta)
		fmt.Println(args...)
		color.Unset()
	}
}
