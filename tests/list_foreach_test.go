package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/gbenroscience/linkedlist/ds"
)

func TestForEach(t *testing.T) {

	list := ds.NewAnyList()

	add(list, 1000)

	go list.ForEach(func(val interface{}) bool {
		fmt.Printf("val-1: %+v\n", val)
		return true
	})

	go list.ForEach(func(val interface{}) bool {
		fmt.Printf("val-2: %+v\n", val)
		return true
	})

	time.Sleep(time.Second * 10)

	fmt.Println("Tests done")

}

func add(list *ds.AnyList, itemCount int) {
	for i := 0; i < itemCount; i++ {
		list.Add(i)
	}
}
