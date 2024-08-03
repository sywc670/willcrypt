package main

import "fmt"

func debug(a ...any) {
	if Debug == true {
		fmt.Println(a...)
	}
}
func debugf(format string, a ...any) {
	if Debug == true {
		fmt.Printf(format, a...)
	}
}
