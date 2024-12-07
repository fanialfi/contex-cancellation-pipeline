package filegeneration

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/fanialfi/contex-cancellation-pipeline/lib"
)

// this funciton are used to generate random file
func GenerateFileConc() {
	err := os.RemoveAll(lib.TempPath)
	if err != nil {
		log.Println("Error remove directory :", err.Error())
	}

	err = os.MkdirAll(lib.TempPath, os.ModePerm)
	if err != nil {
		log.Println("Error create directory :", err.Error())
	}

	// pipeline 1 : job distribution
	chanFileIndex := generateFileIndexes()

	// pipeline 2 : main logic (create file)
	createWorker := 12
	chanFileResult := createFile(chanFileIndex, createWorker)

	// track and print output
	counterSuccess := 0
	counterTotal := 0
	for fileResult := range chanFileResult {
		if fileResult.Err != nil {
			log.Printf("error creating file %s\n\tstack trace : %s", fileResult.FileName, fileResult.Err)
		} else {
			counterSuccess++
		}

		counterTotal++
	}
	log.Printf("%d/%d of file was created", counterSuccess, counterTotal)
}

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

func generateFileIndexes() <-chan lib.FileInfo {
	chanOut := make(chan lib.FileInfo)

	go func() {
		for i := 0; i < lib.TotalFile; i++ {
			chanOut <- lib.FileInfo{
				Index:    i,
				FileName: fmt.Sprintf("file-%d.txt", i),
			}
		}
		close(chanOut)
	}()

	return chanOut
}
