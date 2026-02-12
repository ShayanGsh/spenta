package main

import (
	"fmt"

	"github.com/rouzbehsbz/spenta/iter"
)

func main() {
	arr := []int{1, 2, 3, 4, 5, 6}

	parIter := iter.SliceParMap(&arr, func(a int) int {
		return a * 2
	})

	parIter.Done()

	fmt.Println(arr)
}
