package log

import (
	"fmt"
	"os"
)

// Infof is Printf prefixed with [INFO]
func Infof(format string, a ...interface{}) {
	fmt.Printf("[INFO] "+format+"\n", a...)
}

// Loadf is Printf prefixed with [LOAD]
func Loadf(format string, a ...interface{}) {
	fmt.Printf("[LOAD] "+format+"\n", a...)
}

// Warnf is Fprintf(os.Stderr) prefixed with [WARN]
func Warnf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "[WARN] "+format+"\n", a...)
}

// Errorf is Fprintf(os.Stderr) prefixed with [ERRO]
func Errorf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "[ERRO] "+format+"\n", a...)
}

// Verbosef is Printf prefixed with [VERB]
func Verbosef(format string, a ...interface{}) {
	fmt.Printf("[VERB] "+format+"\n", a...)
}
