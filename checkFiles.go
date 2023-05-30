package goUtils

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// WalkDir 外部调用的遍历目录统计信息的方法
func WalkDir(dirs ...string) string {
	// 保证至少有一个目录需要统计遍历,默认为当前目录
	if len(dirs) == 0 {
		dirs = []string{"."}
	}

	// 初始化变量，channel用于完成Size传递，waitGroup用于等待调度
	fileSizeCh := make(chan int64, 1)
	wg := &sync.WaitGroup{}

	// 启动多个Goroutinue统计信息,取决于len(dirs)
	for _, dir := range dirs {
		wg.Add(1)
		// 并发遍历统计每个目录信息
		go walkDir(dir, fileSizeCh, wg)
	}

	// 启动累计运算的Goroutinue
	// 用户关闭fileSizeCh
	go func(wg *sync.WaitGroup) {
		// 等待统计工作完成
		wg.Wait()
		close(fileSizeCh)
	}(wg)

	// range方式从fileSizeCh中获取文件大小, 统计结果通过channel传递出来
	fileNumCh, sizeTotalCh := make(chan int64, 1), make(chan int64, 1)
	go func(fileSizeCh <-chan int64, fileNumCh, sizeTotalCh chan<- int64) {
		// 统计文件数及文件整体大小
		var fileNum, sizeTotal int64
		for fileSize := range fileSizeCh {
			fileNum++
			sizeTotal += fileSize
		}

		fileNumCh <- fileNum
		sizeTotalCh <- sizeTotal
	}(fileSizeCh, fileNumCh, sizeTotalCh)

	// 整理返回值
	return fmt.Sprintf("%d files %.2f MB\n", <-fileNumCh, float64(<-sizeTotalCh)/1e6)
}

// 遍历并统计某个特定目录的信息，完成递归、统计等
func walkDir(dir string, fileSizeCh chan<- int64, wg *sync.WaitGroup) {
	defer wg.Done()

	// 读取dir下的全部文件并遍历
	for _, fileInfo := range fileInfos(dir) {
		// 根据dir下的文件信息：目录则递归获取信息；不是目录则统计文件大小放入channel
		if fileInfo.IsDir() {
			subDir := filepath.Join(dir, fileInfo.Name())
			// 递归调用,也是并发的
			wg.Add(1)
			go walkDir(subDir, fileSizeCh, wg)
		} else {
			fileSizeCh <- fileInfo.Size() // byte
		}
	}
}

// 获取某个目录下文件信息列表
func fileInfos(dir string) []fs.FileInfo {
	// 读取目录的全部文件
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Println("WalkDir error: ", err)
		return []fs.FileInfo{}
	}

	// 获取文件的文件信息, dirEntry to fileInfo
	infos := make([]fs.FileInfo, 0, len(entries))
	for _, entry := range entries {
		if info, err := entry.Info(); err == nil {
			infos = append(infos, info)
		}
	}

	return infos
}
