package sort

import (
	"fmt"
	"testing"
)

func TestQuickSortSingle(t *testing.T) {
	fmt.Println(QuickSortSingle([]int{3, 6, 8, 2, 9, 0, 1, 7}))
}

func TestQuickSortConcurrency(t *testing.T) {
	fmt.Println(QuickSortConcurrency([]int{3, 6, 8, 2, 9, 0, 1, 7}))
}
