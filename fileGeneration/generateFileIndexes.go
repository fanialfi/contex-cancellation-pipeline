package filegeneration

import (
	"context"
	"fmt"

	"github.com/fanialfi/contex-cancellation-pipeline/lib"
)

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

func generateFileIndexesWithContext(ctx context.Context) <-chan lib.FileInfo {
	chanOut := make(chan lib.FileInfo)

	go func() {
		for i := 0; i < lib.TotalFile; i++ {
			select {
			case <-ctx.Done():
				break
			default:
				chanOut <- lib.FileInfo{
					Index:    i,
					FileName: fmt.Sprintf("file-%d.txt", i),
				}
			}
		}

		close(chanOut)
	}()

	return chanOut
}
