package main

import (
	"EnablementInternship0622/goroutines"
	"time"
)

func main() {
	for i := 1; i <= 10; i++ {
		go func() {
			time.Sleep(1 * time.Second)
			go func() {
				time.Sleep(1 * time.Second)

			}()
		}()
	}

	goroutines.GatherAndWriteGoroutines()
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
