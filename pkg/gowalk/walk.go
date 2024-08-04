package gowalk

import (
	"os"
	"path/filepath"
	"sync"
)

type WalkFunc func(filename string, args ...any)

func Walk(filename string, wg *sync.WaitGroup, walkFn WalkFunc, args ...any) {
	defer wg.Done()

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
			wg.Add(1)
			go Walk(subpath, wg, walkFn, args)
		}
		return
	}
	wg.Add(1)
	go walkFn(filename, args)
}
