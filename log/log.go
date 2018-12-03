package log

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

var infoColor = color.New(color.FgCyan).PrintfFunc()

// Infof is Printf prefixed with [INFO]
func Infof(format string, a ...interface{}) {
	infoColor("[INFO] ")
	fmt.Printf(format+"\n", a...)
}

var loadColor = color.New(color.FgGreen).PrintfFunc()

// Loadf is Printf prefixed with [LOAD]
func Loadf(format string, a ...interface{}) {
	loadColor("[LOAD] ")
	fmt.Printf(format+"\n", a...)
}

var warnColor = color.New(color.FgYellow).FprintfFunc()

// Warnf is Fprintf(os.Stderr) prefixed with [WARN]
func Warnf(format string, a ...interface{}) {
	warnColor(os.Stderr, "[WARN] ")
	fmt.Fprintf(os.Stderr, format+"\n", a...)
}

var errorColor = color.New(color.FgRed).FprintfFunc()

// Errorf is Fprintf(os.Stderr) prefixed with [ERRO]
func Errorf(format string, a ...interface{}) {
	errorColor(os.Stderr, "[ERRO] ")
	fmt.Fprintf(os.Stderr, format+"\n", a...)
}

var verboseColor = color.New(color.FgMagenta).PrintfFunc()

// Verbosef is Printf prefixed with [VERB]
func Verbosef(format string, a ...interface{}) {
	verboseColor("[VERB] ")
	fmt.Printf(format+"\n", a...)
}
