package main

import (
	"fmt"
	"sync"

	"github.com/reugn/async"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(5)
	for range 5 {
		go func() {
			id, _ := async.GoroutineID()
			fmt.Println(id)
			wg.Done()
		}()
	}
	wg.Wait()
}
