package filegeneration

import (
	"context"
	"os"
	"path/filepath"
	"sync"

	"github.com/fanialfi/contex-cancellation-pipeline/lib"
)

func createFile(chanIn <-chan lib.FileInfo, totalWorker int) <-chan lib.FileInfo {
	chanOut := make(chan lib.FileInfo)
	wg := new(sync.WaitGroup)

	wg.Add(totalWorker)
	go func() {
		// dispath to N worker
		for workerIndex := 0; workerIndex < totalWorker; workerIndex++ {
			// do the job (create file)
			go func(workerIndex int) {
				for job := range chanIn {
					filePath := filepath.Join(lib.TempPath, job.FileName)
					contentFile := lib.RandomString(lib.ContentLength)
					err := os.WriteFile(filePath, []byte(contentFile), os.ModePerm)

					// log.Printf("worker %d working on %s file generation", workerIndex, job.FileName)

					// fan-in
					chanOut <- lib.FileInfo{
						FileName:    job.FileName,
						Index:       job.Index,
						WorkerIndex: workerIndex,
						Err:         err,
					}
				}
				wg.Done()
			}(workerIndex)
		}
	}()

	go func() {
		wg.Wait()
		close(chanOut)
	}()

	return chanOut
}

func createFileWithContext(ctx context.Context, chanIn <-chan lib.FileInfo, totalWorker int) <-chan lib.FileInfo {
	chanOut := make(chan lib.FileInfo)

	wg := new(sync.WaitGroup)
	wg.Add(totalWorker)

	go func() {
		for workerIndex := 0; workerIndex < totalWorker; workerIndex++ {
			go func(workerIndex int) {
				for job := range chanIn {
					select {
					case <-ctx.Done():
						break
					default:
						filepath := filepath.Join(lib.TempPath, job.FileName)
						content := lib.RandomString(lib.ContentLength)
						err := os.WriteFile(filepath, []byte(content), os.ModePerm)

						// log.Println("worker", workerIndex, "working", job.FileName, "file generation")

						chanOut <- lib.FileInfo{
							Index:       job.Index,
							FileName:    job.FileName,
							WorkerIndex: workerIndex,
							Err:         err,
						}
					}
				}

				wg.Done()
			}(workerIndex)
		}
	}()

	go func() {
		wg.Wait()
		close(chanOut)
	}()

	return chanOut
}
