package logger

import (
	"fmt"
	"os"
)

// Log function definitions.

func Info(format string, a ...any) {
	fmt.Fprintf(os.Stdout, format+"\n", a...)
}

func Warn(format string, a ...any) {
	fmt.Fprintf(os.Stdout, format+"\n", a...)
}

func Err(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", a...)
}
