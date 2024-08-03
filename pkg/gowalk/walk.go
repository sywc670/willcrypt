package gowalk

import (
	"os"
	"path/filepath"
)

type WalkFunc func(filename string, args ...any)

func Walk(filename string, walkFn WalkFunc, args ...any) {

	fileinfo, err := os.Stat(filename)
	if err != nil {
		panic(err)
	}

	if fileinfo.IsDir() {
		entries, err := os.ReadDir(filename)
		if err != nil {
			panic(err)
		}

		for _, entry := range entries {
			subpath := filepath.Join(filename, entry.Name())
			Walk(subpath, walkFn, args)
		}
		return
	}

	walkFn(filename, args)
}
