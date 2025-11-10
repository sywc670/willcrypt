package main

import "fmt"

func debug(a ...any) {
	if c.Debug == true {
		fmt.Println(a...)
	}
}
func debugf(format string, a ...any) {
	if c.Debug == true {
		fmt.Printf(format, a...)
	}
}
