package main

import (
	"EnablementInternship0622/goroutines"
	"fmt"
	"sync"
	"time"
)

func player(table chan int) {
	for {
		ball := <-table
		ball++
		time.Sleep(100 * time.Millisecond)
		table <- ball
	}
}

func sendMessage(ch chan<- string, message string, sleepDuration time.Duration) {
	for {
		ch <- message
		time.Sleep(sleepDuration)
	}
}

func main() {
	c1 := make(chan string)
	c2 := make(chan string)
	quit := make(chan int)
	var wg sync.WaitGroup
	wg.Add(5)

	go func() {
		sendMessage(c1, "test1", 2*time.Second)
		wg.Done()
	}()

	go func() {
		sendMessage(c2, "test2", 4*time.Second)
		wg.Done()
	}()

	go func() {
		time.Sleep(10 * time.Second)
		quit <- 0
		wg.Done()
	}()

	var Ball int
	table := make(chan int)

	for i := 0; i < 3; i++ {
		go func() {
			player(table)
			wg.Done()
		}()
	}

	table <- Ball
	time.Sleep(1 * time.Second)
	<-table

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
			wg.Wait()
			return
		default:
			cnt++
			fmt.Printf("(cnt: %v)\n", cnt)
			time.Sleep(1 * time.Second)
		}
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
