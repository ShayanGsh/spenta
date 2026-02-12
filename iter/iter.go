package iter

import "sync"

const MIN_CHUNK_SIZE uint = 512

type ParIter struct {
	wg *sync.WaitGroup
}

func (p *ParIter) Done() {
	p.wg.Wait()
}

type ParIterOptions struct {
	MinChunkSize uint
}

func DefaultParIterOptions() *ParIterOptions {
	return &ParIterOptions{
		MinChunkSize: MIN_CHUNK_SIZE,
	}
}

func WithMinChunkSize(size uint) ParIterOptions {
	return ParIterOptions{
		MinChunkSize: size,
	}
}

func BuildParIterOptions(opts []ParIterOptions) ParIterOptions {
	o := DefaultParIterOptions()

	for _, opt := range opts {
		o.MinChunkSize = opt.MinChunkSize
	}

	return *o
}

func ChunkSize(len int, minSize uint) int {
	size := (len + 1) / 2
	if size > int(minSize) {
		return ChunkSize(size, minSize)
	}

	return size
}

func ChunkCount(len, chunkSize int) int {
	return (len + chunkSize - 1) / chunkSize
}

func SliceChunk[T comparable](slice *[]T, minChunkSize uint) (int, int, int) {
	len := len((*slice))
	chunkSize := ChunkSize(len, minChunkSize)
	chunkCount := ChunkCount(len, chunkSize)

	return len, chunkSize, chunkCount
}
