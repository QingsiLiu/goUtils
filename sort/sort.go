package sort

import "sync"

// QuickSortSingle 快速排序单线程实现
func QuickSortSingle(arr []int) []int {
	if arr == nil || len(arr) < 2 {
		return arr
	}
	quickSortSingle(arr, 0, len(arr)-1)
	return arr
}

func quickSortSingle(arr []int, l, r int) {
	// 判断待排序范围是否合法
	if l < r {
		mid := partition(arr, l, r)
		// 递归排序左边
		quickSortSingle(arr, l, mid-1)
		// 递归排序右边
		quickSortSingle(arr, mid+1, r)
	}
}

// 大小分区，返回参考元素索引
func partition(arr []int, l, r int) int {
	p := l - 1
	for i := l; i <= r; i++ {
		if arr[i] <= arr[r] {
			p++
			arr[p], arr[i] = arr[i], arr[p]
		}
	}
	return p
}

func QuickSortConcurrency(arr []int) []int {
	if arr == nil || len(arr) < 2 {
		return arr
	}

	// 同步的控制
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go quickSortConcurrency(arr, 0, len(arr)-1, wg)
	wg.Wait()

	return arr
}

func quickSortConcurrency(arr []int, l, r int, wg *sync.WaitGroup) {
	defer wg.Done()

	if l < r {
		// 大小分区元素，并获取参考元素索引
		mid := partition(arr, l, r)
		// 并发对左部分排序
		wg.Add(1)
		go quickSortConcurrency(arr, l, mid-1, wg)
		// 并发对右部分排序
		wg.Add(1)
		go quickSortConcurrency(arr, mid+1, r, wg)
	}
}
