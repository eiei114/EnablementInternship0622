package goroutines

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

var relationships map[string]bool

func parseGoroutines(stack string) []Goroutine {
	common := regexp.MustCompile(`goroutine (\d+) \[(.+)\]:`)
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
	for relationship := range relationships {
		lines = append(lines, relationship)
	}
	lines = append(lines, "}")

	ioutil.WriteFile(fileName, []byte(strings.Join(lines, "\n")), 0644)
}

func GatherAndWriteGoroutines() {
	relationships = make(map[string]bool)

	timer := time.NewTimer(20 * time.Second)
	buf := make([]byte, 8192)

	for {
		select {
		case <-timer.C:
			writeDotFile("goroutines2.dot")
			fmt.Println("Done!!")
			return
		default:
			n := runtime.Stack(buf, true)
			fmt.Printf("=== Stack trace ===\n%s\n", buf[:n])

			goroutines := parseGoroutines(string(buf[:n]))

			for _, g := range goroutines {
				relationship := ""
				if g.Status == "chan send" {
					relationship = fmt.Sprintf("\"%s\" -> \"%s\" [label = \"%s\" color = red];", g.ID, g.ParentID, g.Status)
					relationships[relationship] = true
					relationship = fmt.Sprintf("\"%s\" -> \"%s\";", g.ParentID, g.ID)
					relationships[relationship] = true
				} else if g.Status == "chan receive" {
					relationship = fmt.Sprintf("\"%s\" -> \"%s\" [label = \"%s\" color = blue];", g.ParentID, g.ID, g.Status)
					relationships[relationship] = true
					relationship = fmt.Sprintf("\"%s\" -> \"%s\";", g.ParentID, g.ID)
					relationships[relationship] = true
				} else if g.ParentID != "" {
					relationship = fmt.Sprintf("\"%s\" -> \"%s\";", g.ParentID, g.ID)
					relationships[relationship] = true
				}
				time.Sleep(1 * time.Millisecond)
			}
		}
	}
}
