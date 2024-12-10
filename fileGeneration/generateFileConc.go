package filegeneration

import (
	"log"
	"os"

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
