package tests

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/gbenroscience/linkedlist/ds"
)

func TestForEach(t *testing.T) {

	list := ds.NewAnyList[int]()

	add(list, 1000)

	var wg sync.WaitGroup

	wg.Add(2)

	list.ForEach(func(val int) bool {
		fmt.Printf("val-0: %+v\n", val)
		return true
	})

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		list.ForEach(func(val int) bool {
			fmt.Printf("val-1: %+v\n", val)
			return true
		})
		fmt.Println("First goroutine done")
	}(&wg)
	time.Sleep(5 * time.Second)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		list.ForEach(func(val int) bool {
			fmt.Printf("val-2: %+v\n", val)
			return true
		})
		fmt.Println("Second goroutine done")
	}(&wg)

	wg.Wait()

	fmt.Println("All tests done")

}

func add(list *ds.AnyList[int], itemCount int) {
	for i := 0; i < itemCount; i++ {
		list.Add(i)
	}
	fmt.Println("Done adding " + strconv.Itoa(itemCount) + " items to list")
}
