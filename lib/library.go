package lib

import (
	"math/rand"
	"time"
)

const (
	TotalFile     = 3000
	ContentLength = 100 << 10 // content length sebesar 512KB
)

var TempPath = "/dev/shm/cancellation"

type FileInfo struct {
	Index       int
	FileName    string
	WorkerIndex int
	Err         error
}

func RandomString(length int) string {
	randomizer := rand.New(rand.NewSource(time.Now().UnixNano()))
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ ")

	b := make([]rune, length)

	for i := range b {
		b[i] = letters[randomizer.Intn(len(letters))]
	}

	return string(b)
}
