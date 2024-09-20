package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/fsnotify/fsnotify"
)

func main() {
	// 指定要监听的目录
	watchPath := "./testdir" // 这里替换成你要监听的目录路径

	// 创建一个新的Watcher实例
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// 开始监听指定的目录
	err = watcher.Add(watchPath)
	if err != nil {
		log.Fatal(err)
	}

	// 监听文件系统的改变
	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// 检查是否是写入或创建事件
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
					fmt.Println("File modified or created:", event.Name)

					// 打印文件新增内容
					printNewContent(event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// 阻塞主goroutine以保持程序运行
	<-done
}

// printNewContent 打开文件并打印其最新内容
func printNewContent(filePath string) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("Failed to read file: %s", err)
		return
	}
	fmt.Println("New content of the file:", string(content))
}
