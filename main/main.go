package main

import (
	"fmt"

	"github.com/gbenroscience/linkedlist/ds"
)

func runComparableListExamples() {

	list := ds.NewList[int]()
	list.Add(2)
	list.Add(4)
	list.Add(6)

	list.Log("Check-Contents-1")
	list.Add(8)
	list.Log("Check-Contents-2")
	list.AddVal(3, 1)
	list.Log("Check-Contents-3")

	list.Set(0, 10000)
	list.Log("Check-Contents-4")

	var a []int = []int{20, 40, 60, 80, 100}
	list.AddArray(a)
	list.Log("Check-Contents-5")

	subList, err := list.SubList(3, 8)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	subList.Log("Check-sublist-contents-1")
	subList.RemoveIndex(2)
	subList.RemoveIndex(2)
	subList.Log("Check-sublist-contents-2")
	list.Log("Check-Contents-6")

	subList.AddValues(200, 250, 300, 350, 400)
	subList.Log("Check-sublist-contents-3")
	list.Log("Check-Contents-7")
}

func runAnyListExamples() {

	list := ds.NewAnyList[int]()
	list.Equals = func(val1, val2 int) bool {
		return val1 == val2
	}
	list.Add(2)
	list.Add(4)
	list.Add(6)

	list.Log("Check-Contents-1")
	list.Add(8)
	list.Log("Check-Contents-2")
	list.AddVal(3, 1)
	list.Log("Check-Contents-3")

	list.Set(0, 10000)
	list.Log("Check-Contents-4")

	var a []int = []int{20, 40, 60, 80, 100}
	list.AddArray(a)
	list.Log("Check-Contents-5")

	subList, err := list.SubList(3, 8)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	subList.Log("Check-sublist-contents-1")
	subList.RemoveIndex(2)
	subList.RemoveIndex(2)
	subList.Log("Check-sublist-contents-2")
	list.Log("Check-Contents-6")

	subList.AddValues(200, 250, 300, 350, 400)
	subList.Log("Check-sublist-contents-3")
	list.Log("Check-Contents-7")
}

func main() {
	fmt.Println("----------------COMPARABLE LIST EXAMPLES----------------")
	runAnyListExamples()
	fmt.Println("----------------ANY LIST EXAMPLES----------------")
	runComparableListExamples()
}
