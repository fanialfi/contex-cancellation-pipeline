package main

import (
	"log"
	"time"

	filegeneration "github.com/fanialfi/contex-cancellation-pipeline/fileGeneration"
)

func main() {
	log.Println("start")
	start := time.Now()

	filegeneration.GenerateFileConc()

	duration := time.Since(start)
	log.Printf("done in %.3f seconds", duration.Seconds())
}
