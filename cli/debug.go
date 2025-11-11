package main

import "fmt"

func debug(a ...any) {
	if cfg.Verbose == true {
		fmt.Println(a...)
	}
}
func debugf(format string, a ...any) {
	if cfg.Verbose == true {
		fmt.Printf(format, a...)
	}
}
