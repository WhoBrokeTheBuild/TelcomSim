package main

import "fmt"

func Infof(format string, a ...interface{}) {
	fmt.Printf("[INFO] "+format+"\n", a...)
}

func Loadf(format string, a ...interface{}) {
	fmt.Printf("[LOAD] "+format+"\n", a...)
}

func Warnf(format string, a ...interface{}) {
	fmt.Printf("[WARN] "+format+"\n", a...)
}

func Errorf(format string, a ...interface{}) {
	fmt.Printf("[ERRO] "+format+"\n", a...)
}

func Verbosef(format string, a ...interface{}) {
	fmt.Printf("[VERB] "+format+"\n", a...)
}
