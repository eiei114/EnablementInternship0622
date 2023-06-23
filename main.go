package main

import (
	"EnablementInternship0622/goroutines"
	"fmt"
	"runtime"
	"time"
)

func test1(ch chan<- string) {
	for {
		ch <- "test1"
		time.Sleep(2 * time.Second)
	}
}

func test2(ch chan<- string) {
	for {
		ch <- "test2"
		time.Sleep(4 * time.Second)
	}
}

func test3(quit chan<- int) {
	time.Sleep(10 * time.Second)
	quit <- 0
}

func main() {
	c1 := make(chan string)
	c2 := make(chan string)
	quit := make(chan int)
	go test1(c1)
	go test2(c2)
	go test3(quit)

	cnt := 0
	goroutines.GatherAndWriteGoroutines()

	for {
		select {
		case s1 := <-c1:
			fmt.Println(s1)
		case s2 := <-c2:
			fmt.Println(s2)
		case <-quit:
			fmt.Println("quit")
			break
		default:
			cnt = cnt + 1
			fmt.Printf("(cnt: %v)\n", cnt)
			time.Sleep(1 * time.Second)
		}
	}

	buf := make([]byte, 8192)
	for {
		n := runtime.Stack(buf, true) // すべてのゴルーチンのスタックトレースを取得
		fmt.Printf("=== Stack trace ===\n%s\n", buf[:n])
		time.Sleep(1 * time.Second) // 5秒ごとにスタックトレースを取得
	}
}

//必要な情報
//
//時間
//ゴール―チンID
//ゴール―チンの状態（どのような状態があるか確認:running,runnable,sleepまで確認できた。）
//親ルーチンID（親があれば）
//取り合えずbuf[:n]から先ほどの情報をパースして配列などにまとめていく。
//
//まとめた情報をどのようにビジュアライズしていくか確認。
