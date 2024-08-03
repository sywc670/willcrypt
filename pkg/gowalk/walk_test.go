package gowalk

import (
	"testing"
)

func TestWalk(t *testing.T) {
	Walk("testdir", func(filename string, a ...any) { print(filename) })
}
