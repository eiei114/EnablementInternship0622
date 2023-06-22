package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	go func() {
		for {
			time.Sleep(1 * time.Second)
		}
	}()

	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, true) // すべてのゴルーチンのスタックトレースを取得
		fmt.Printf("=== Stack trace ===\n%s\n", buf[:n])
		time.Sleep(5 * time.Second) // 5秒ごとにスタックトレースを取得
	}
}
