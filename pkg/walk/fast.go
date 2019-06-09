// Package walk contains methods for file system walking.
package walk

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type (
	walkArguments struct {
		path string
		info os.FileInfo
		err  error
	}
)

// Fast concurrently calls filepath.Walk using the provided walk function for the provided
// directory. The maximum number of goroutines to use is specified by the 'concurrency' parameter.
// In larger repos, it may be beneficial to increase this limit.
func Fast(root string, walkFn filepath.WalkFunc, concurrency int) error {
	// Ensure a sensible value is set
	if concurrency <= 0 {
		concurrency = 1
	}

	var filesGroup sync.WaitGroup
	var skip []string

	files := make(chan *walkArguments)
	stop := make(chan bool, 1)
	errs := make(chan error)

	for i := 0; i < concurrency; i++ {
		go func() {
			for {
				select {
				case file, ok := <-files:
					if !ok {
						continue
					}

					err := walkFn(file.path, file.info, file.err)

					if err == filepath.SkipDir {
						skip = append(skip, file.path)
					} else if err != nil {
						errs <- err
					}

					filesGroup.Done()
				case <-stop:
					return
				}
			}
		}()
	}

	var walkErr error

	go func() {
		select {
		case walkErr = <-errs:
			stop <- true
		case <-stop:
			return
		}
	}()

	var walkerWg sync.WaitGroup
	walkerWg.Add(1)

	go func() {
		err := filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
			select {
			case <-stop:
				close(files)
				return errors.New("walker received interrupt signal")
			default:
				for _, d := range skip {
					if strings.HasPrefix(p, d) {
						return nil
					}
				}

				filesGroup.Add(1)
				files <- &walkArguments{path: p, info: info, err: err}
				return nil
			}
		})

		if err != nil && err != io.EOF {
			errs <- err
		}

		walkerWg.Done()
	}()

	walkerWg.Wait()

	if walkErr == nil {
		filesGroup.Wait()
		stop <- true
	}

	return walkErr
}
