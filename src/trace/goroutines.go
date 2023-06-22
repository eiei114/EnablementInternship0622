package main

import (
	"fmt"
	"runtime"
	"sync"
)

type Goroutine struct {
	Id        int
	ParentId  int
	Stack     string
	Children  []*Goroutine
	isRunning bool
}

var (
	mutex      sync.Mutex
	goroutines = make(map[int]*Goroutine)
	nextId     = 1
)

func StartGoroutine(parentId int, f func()) int {
	mutex.Lock()
	defer mutex.Unlock()

	buf := make([]byte, 1024)
	runtime.Stack(buf, false)
	stack := string(buf)

	g := &Goroutine{
		Id:        nextId,
		ParentId:  parentId,
		Stack:     stack,
		Children:  []*Goroutine{},
		isRunning: true,
	}
	goroutines[nextId] = g

	if parent, ok := goroutines[parentId]; ok {
		parent.Children = append(parent.Children, g)
	}

	nextId++

	go func() {
		f()
		g.isRunning = false
	}()

	return g.Id
}

func PrintAll() {
	mutex.Lock()
	defer mutex.Unlock()

	for _, g := range goroutines {
		if g.isRunning {
			fmt.Printf("Goroutine %d (parent: %d) is running. Stack trace:\n%s\n", g.Id, g.ParentId, g.Stack)
		} else {
			fmt.Printf("Goroutine %d (parent: %d) has stopped.\n", g.Id, g.ParentId)
		}
	}
}
