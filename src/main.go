package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"runtime"
	"strings"
	"time"
)

type Goroutine struct {
	ID       string
	Status   string
	ParentID string
}

var relationships []string

func parseGoroutines(stack string) []Goroutine {
	common := regexp.MustCompile(`goroutine (\d+) \[(\w+)\]:`)
	child := regexp.MustCompile(`created by .+ in goroutine (\d+)`)

	lines := strings.Split(stack, "\n")
	var goroutines []Goroutine
	var lastGoroutine *Goroutine

	for _, line := range lines {
		if match := common.FindStringSubmatch(line); match != nil {
			g := Goroutine{
				ID:     match[1],
				Status: match[2],
			}
			goroutines = append(goroutines, g)
			lastGoroutine = &goroutines[len(goroutines)-1]
		} else if match := child.FindStringSubmatch(line); match != nil && lastGoroutine != nil {
			lastGoroutine.ParentID = match[1]
		}
	}

	return goroutines
}

func writeDotFile(fileName string) {
	lines := []string{"digraph G {"}
	lines = append(lines, relationships...)
	lines = append(lines, "}")

	ioutil.WriteFile(fileName, []byte(strings.Join(lines, "\n")), 0644)
}

func main() {

	go func() {
		time.Sleep(1 * time.Second)
	}()

	go func() {
		time.Sleep(1 * time.Second)
	}()

	go func() {
		time.Sleep(1 * time.Second)
		go func() {
			time.Sleep(1 * time.Second)
			go func() {
				time.Sleep(1 * time.Second)
			}()
		}()
	}()

	timer := time.NewTimer(20 * time.Second)

	buf := make([]byte, 1024)
	for {
		select {
		case <-timer.C:
			writeDotFile("goroutines.dot") // Write the .dot file here.
			return
		default:
			n := runtime.Stack(buf, true)
			goroutines := parseGoroutines(string(buf[:n]))
			for _, g := range goroutines {
				fmt.Printf("Goroutine ID: %s, Status: %s, Parent ID: %s\n", g.ID, g.Status, g.ParentID)
				if g.ParentID != "" {
					relationships = append(relationships, fmt.Sprintf("\"%s\" -> \"%s\";", g.ParentID, g.ID))
				}
			}
			time.Sleep(5 * time.Second)
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
