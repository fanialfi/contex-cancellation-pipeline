package filegeneration

import (
	"context"
	"log"
	"os"

	"github.com/fanialfi/contex-cancellation-pipeline/lib"
)

// this function useful for generate random file with context
func GenerateFileWithContext(ctx context.Context) {
	err := os.RemoveAll(lib.TempPath)
	if err != nil {
		log.Println("Error removing directory", err.Error())
	}
	err = os.MkdirAll(lib.TempPath, os.ModePerm)
	if err != nil {
		log.Println("Error creating directory", err.Error())
	}

	done := make(chan int)

	go func() {
		// pipeline 1 : job distribution
		chanFileIndex := generateFileIndexesWithContext(ctx)

		// pipeline 2 : the main logic (creating file)
		createFileWorker := 100
		chanFileResult := createFileWithContext(ctx, chanFileIndex, createFileWorker)

		// track and print output
		counterSuccess := 0
		for fileResult := range chanFileResult {
			if fileResult.Err != nil {
				log.Printf("error creating file %s. stack trace : %s\n", fileResult.FileName, fileResult.Err)
			} else {
				counterSuccess++
			}
		}

		// notif that the process is complete
		done <- counterSuccess
	}()

	select {
	case <-ctx.Done():
		log.Printf("generation process stoped. %s\n", ctx.Err().Error())
	case counterSuccess := <-done:
		log.Printf("%d/%d of total file created\n", counterSuccess, lib.TotalFile)
	}
}
