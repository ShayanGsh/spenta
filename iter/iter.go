package iter

import "sync"

type ParIter struct {
	wg *sync.WaitGroup
}

func (p *ParIter) Done() {
	p.wg.Wait()
}

func ChunkSize(len, workers int) int {
	targetChunks := workers * 4
	size := len / targetChunks

	if size < 1 {
		return 1
	}

	return size
}

func ChunkCount(len, chunkSize int) int {
	return (len + chunkSize - 1) / chunkSize
}

func SliceChunk[T comparable](slice *[]T, workers int) (int, int, int) {
	len := len((*slice))
	chunkSize := ChunkSize(len, workers)
	chunkCount := ChunkCount(len, chunkSize)

	return len, chunkSize, chunkCount
}
