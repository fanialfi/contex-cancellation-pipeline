package main

import (
	"context"
	"log"
	"time"

	filegeneration "github.com/fanialfi/contex-cancellation-pipeline/fileGeneration"
)

const timeoutDuration = 3 * time.Second

func main() {
	log.Println("start")
	start := time.Now()

	// pada saat pembuatan context.WIthTimeout memerlukan / membutuhkan Context
	// dan belum ada object Context yang dibuat, maka buat object Context baru dengan menggunakan function context.Background
	// context.Background ini menghasilkan object Context kosong dan tidak memiliki deadline
	// biasanya digunakan untuk inisialisasi object context.Context baru yang akan dichain dengan function context.With...
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()
	filegeneration.GenerateFileWithContext(ctx)

	duration := time.Since(start)
	log.Printf("done in %.3f seconds", duration.Seconds())
}
