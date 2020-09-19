package main

import (
	"fmt"
	"github.com/gbenroscience/linkedlist/dsn"
	"github.com/gbenroscience/linkedlist/utils"
	"strconv"
	"time"
)

func testAdd(n int) *dsn.List {

	list := dsn.NewList()
	for i := 0; i < n; i++ {
		list.Add(i)
	}

	return list
}

func testAddFromArgs(itemsToAdd ...interface{}) *dsn.List {

	list := dsn.NewList()
	list.AddValues(itemsToAdd...)

	list.Log("testAddFromArgs")

	return list
}

func testAddArray(items []int) *dsn.List {

	list := dsn.NewList()
	list.Log("testAddArray-Begins")
	defer list.Log("testAddArray-Ends")

	arrInterface := make([]interface{}, len(items))
	for i, v := range items {
		arrInterface[i] = v
	}
	list.AddArray(arrInterface)

	return list
}

func testAddValAtIndex(n int, index int, itemsToAdd ...interface{}) *dsn.List {

	list := testAdd(n)
	list.Log("testAddValAtIndex generated items")

	for _, v := range itemsToAdd {
		list.AddVal(v, index)
	}

	list.Log("testAddValAtIndex... added " + strconv.Itoa(len(itemsToAdd)) + " items at index " + strconv.Itoa(index))

	return list
}

func testAddAll(n int, lst *dsn.List) *dsn.List {

	list := dsn.NewList()
	defer list.Log("testAddAll after adding list!")
	for i := 0; i < n; i++ {
		list.Add(i)
	}
	list.Log("testAddAll after adding initial contents")

	list.AddAll(lst)

	return list
}

func testAddListAtIndex(n int, index int, lst *dsn.List) {

	list := testAdd(n)
	list.Log("starting testAddListAtIndex")

	list.AddAllAt(index, lst)

	list.Log("testAddListAtIndex... after adding list at index " + strconv.Itoa(index))

}

func testLog(list *dsn.List) {

	list.Log("testLog")

}
func testClear(list *dsn.List) {

	list.Clear()

	list.Log("testClear")

}

func runSuite1() {

	list := testAdd(910)
	testLog(list)
	testClear(list)

	testAddFromArgs(22, 54, 45, 67, 76, 9098, -9802, 2345, 12, 21)

	listFromAddValAtIndex := testAddValAtIndex(10, 3, 22, 33, 44, 99, 88, 77, -909, 299)

	testAddListAtIndex(20, 8, listFromAddValAtIndex)

	rnd := utils.NewRnd()

	arr := rnd.GenerateRndArray(32, 120, false)
	lst := testAddArray(arr)

	lst.Log("Before being added to new list")

	testAddAll(20, lst)


}

func testRemoveIndex(list *dsn.List, index int) {
	list.Log("Remove from index: before")
	list.RemoveIndex(index)
	list.Log("Remove from index: after")
}

func testRemoveVal(list *dsn.List, val interface{}) {
	list.Log("Remove value: before")
	list.Remove(val)
	list.Log("Remove value: after")
}

func testRemoveAll(list *dsn.List) {

	lst := dsn.NewList()
	lst.AddValues(0, 1, 2, 3, 4, 101, 800)
	lst.Log("removables")
	list.Log("Remove list: before")
	list.RemoveAll(lst)
	list.Log("Remove list: after")
}

func runSuite2() {
	list := testAdd(910)
	testRemoveIndex(list, 12)
	testRemoveVal(list, 102)
	testRemoveAll(list)

}

func testGetAtIndex() {
	rnd := utils.NewRnd()
	sz := 910
	list := testAdd(sz)

	fmt.Printf("List-len: %d\n", list.Count())
	for i := 0; i < 10; i++ {
		index := rnd.NextInt(sz)
		fmt.Printf("index: %d\n", index)
		val := list.Get(index)

		fmt.Printf("Retrieved %v at index %d\n ", val, index)
	}

}

func testIndexOf() {
	rnd := utils.NewRnd()
	sz := 910
	list := testAdd(sz)

	rndArr := rnd.GenerateRndArray(200, 500, false)

	for _, v := range rndArr {
		index := list.IndexOf(v)
		fmt.Printf("Index of %v in list is %d\n ", v, index)
	}

}

func runSuite3() {

	testGetAtIndex()
	testIndexOf()

}

func testSubListCreate(parenSize int, start int, end int) {

	list := testAdd(parenSize)

	list.Log("MainList")

	subList, err := list.SubList(start, end)

	if err != nil {
		fmt.Printf("error: %v", err)
	}
	subList.Log("SubList")

}

func testForEach(parenSize int, start int, end int) {

	list := testAdd(parenSize)

	list.Log("MainList")

	list.ForEach(func(val interface{}) bool {
		fmt.Printf("See list elem: %v\n", val)
		return true
	})

	subList, err := list.SubList(start, end)
	if err != nil {
		fmt.Printf("error: %v", err)
	}

	subList.ForEach(func(val interface{}) bool {
		fmt.Printf("See sub-list elem: %v\n", val)
		return true
	})

}

func testToArray(parenSize int, start int, end int) {

	list := testAdd(parenSize)

	list.Log("MainList")

	list.ForEach(func(val interface{}) bool {
		fmt.Printf("See list elem: %v\n", val)
		return true
	})

	subList, err := list.SubList(start, end)
	if err != nil {
		fmt.Printf("error: %v", err)
	}

	ar := list.ToArray()
	fmt.Println("list.ToArray:")
	for _, v := range ar {
		fmt.Printf("%v, ", v)
	}

	arr := subList.ToArray()
	fmt.Println("\nsublist.ToArray:")
	for _, v := range arr {
		fmt.Printf("%v, ", v)
	}

}

func testSubListModify(parenSize int, start int, end int) {

	fmt.Println("TEST_SUBLIST_MODIFY...")
	list := testAdd(parenSize)

	list.Log("MainList")

	subList, err := list.SubList(start, end)
	if err != nil {
		fmt.Printf("error: %v", err)
	}

	subList.Log("subList_before_modifying")

	subList.Clear()

	subList.Log("subList_after_modifying")

	list.Log("MainList_after_modifying_sublist")

}

func testSubListWoesWhenParentChanged(parenSize int, start int, end int) {

	fmt.Println("testSubListWoesWhenParentChanged...")
	list := testAdd(parenSize)

	list.Log("MainList...")

	subList, err := list.SubList(start, end)
	if err != nil {
		fmt.Printf("error: %v", err)
	}

	subList.Log("subList_before_modifying")

	list.Remove(start)

	list.Log("MainList after removing sublist_head...")

	subList.Log("subList_after_being_beheaded")
	subList.Close()

}

func testSubListClear(parenSize int, start int, end int) {

	fmt.Println("testSubListClear...")
	list := testAdd(parenSize)

	list.Log("MainList...")

	subList, err := list.SubList(start, end)
	if err != nil {
		fmt.Printf("error: %v", err)
	}

	subList.Log("subList_before_clearing")

    subList.Clear()

	subList.Log("subList_after_clearing")


	list.Log("MainList... after clearing sublist")

	list.Clear()

	list.Log("MainList... after clearing")



}

func runSuite4() {

	testSubListCreate(50, 20, 40)

	testForEach(50, 20, 40)

	testToArray(50, 20, 40)

	testSubListModify(50, 20, 40)


	testSubListWoesWhenParentChanged(50, 20, 40)

	testSubListClear(50, 22, 30)
}

func runSuite5(){

	list := testAdd(10)
	subList , _ := list.SubList(3, 7)
	list.Log("MainList")
	subList.Log("SubList")

	subList.Clear()

	list.Log("MainList")
	subList.Log("SubList")


	subList.AddValues(2,9,8)

	list.Log("MainList")
	subList.Log("SubList")



}
func main() {

/*	runSuite1()
	runSuite2()
	runSuite3()

	runSuite4()
*/

	runSuite5()
	time.Sleep(time.Second * 4)

}
