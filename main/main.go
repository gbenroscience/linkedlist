package main

import (
	"bytes"
	"fmt"
	"github.com/gbenroscience/linkedlist/ds"
	"github.com/gbenroscience/linkedlist/utils"
	"math"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type Rect struct {
	length  int
	breadth int
	height  int
	color   string
	Square  *Rect
}

func personal() {

	value := &Rect{
		length:  25,
		breadth: 302,
		height:  128,
		color:   "red",
		Square: &Rect{
			length:  5,
			breadth: 8,
			height:  17,
			color:   "blue",
			Square:  nil,
		},
	}
	fmt.Printf("Dimen %d\n", value.length*value.breadth*value.height*value.Square.length*value.Square.breadth*value.Square.height)
}

func approxSquareRoot(x float64) {
	var root float64 = rand.Float64()
	var iter int = 20

	for i := 0; i < iter; i++ {
		root = root - ((root*root)-x)/(2*root)
	}

	fmt.Printf("Approx Root: %f Real Root: %17f\n", root, math.Sqrt(x))

}
func approxCubeRoot(x float64) {
	var root float64 = rand.Float64()
	var iter int = 20

	for i := 0; i < iter; i++ {
		root = root - ((root*root*root)-x)/(3*root*root)
	}

	fmt.Printf("Approx CubeRoot: %7f Real CubeRoot: %7f\n", root, math.Cbrt(x))

}

func appendText(str1 string, str2 string) string {

	var buf bytes.Buffer

	buf.WriteString(str1)
	buf.WriteString(str2)
	result := buf.String()

	fmt.Println(result)

	return result
}

func test() {
	appendText("GOD ", "IS HERE!!!")

	data := new(ds.List)
	data = ds.NewList()

	data.Add(3)
	data.Add(8)
	data.Add("9")

	fmt.Println(data.Get(2))

	testSize := 200
	rand.New(rand.NewSource(time.Now().UnixNano()))
	n := 200000
	start := time.Now().UnixNano()
	for I := 0; I < testSize; I++ {
		data.Add(rand.Intn(n))
	}
	stop := time.Now().UnixNano()

	fmt.Println("Duration for adding " + strconv.Itoa(testSize) + " items to LinkedList is: " + strconv.Itoa(int(stop-start)/1000000) + " ms")

	fmt.Println("Done Massive Element Adding")

	data.Log("MAIN_LIST")

	smallList, err := data.SubList(5, 12)

	if err == nil {

		smallList.Log("SUB_LIST")

		smallList.Clear()

		data.Log("AFTER CLEARING SUB_LIST, MAIN_LIST")
		smallList.Log("SUB_LIST")

		second := ds.NewList()

		second.Add(5)
		second.Add(205)
		second.Add(502)
		second.Add(45)
		second.Add(67)
		second.Add(99)
		second.Add(510)
		second.Add(1020)
		second.Add(5008)
		second.Add(30042)
		second.Add(20014)
		second.Add(30041)
		second.Add(10518902)

		second.Log("SECOND-AAAA")
		s, err := second.SubList(2, 5)
		s.Log("S-AAAA")
		if err == nil {
			s.Clear()
		}
		second.Log("SECOND-AAAA---CLEAR")

		second.Log("SECOND_LIST")

		smallList.AddAll(second)

		smallList.Log("SUB_LIST absorbed SECOND_LIST..result")

		data.Log("MAIN_LIST")

		data.Clear()

		smallList.Log("MAIN LIST CLEARED: SEE SUB_LIST")

		data.Log("MAIN_LIST")
	} else {
		fmt.Printf("%v", err)
	}

}
func Print(x interface{})  bool{
	fmt.Printf("Printing list: found %d\n", x)
	return true
}

func test1() {

	list := ds.NewList()

	var wg sync.WaitGroup

	for j := 1; j <= 4; j++ {

		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				list.Add(i)
			}

			val , err := list.Get(3)
			if err != nil{
				fmt.Printf("error %v\n", err)
			}
			fmt.Printf("Elem @ index %d is %v\n",3, val)
			fmt.Printf("list now has %d elements\n", list.Count())
		}()
	}

	wg.Wait()

	list.ForEach(Print)

	list.RemoveIndex(3)

	list.Log("Checking...")

	time.Sleep(time.Second * 3)

	fmt.Println(list.Get(3))
}

func test2() {

	list := ds.NewList()

	for i := 0; i < 10; i++ {
		list.Add(i)
	}
	fmt.Printf("list now has %d elements\n", list.Count())

	list.ForEach(Print)

	list.RemoveIndex(3)

	list.Log("Checking...")

	time.Sleep(time.Second * 10)

	fmt.Println(list.Get(3))
}
















func testAdd(n int) *ds.List {

	list := ds.NewList()
	for i := 0; i < n; i++ {
		list.Add(i)
	}

	return list
}

func testAddFromArgs(itemsToAdd ...interface{}) *ds.List {

	list := ds.NewList()
	list.AddValues(itemsToAdd...)

	list.Log("testAddFromArgs")

	return list
}

func testAddArray(items []int) *ds.List {

	list := ds.NewList()
	list.Log("testAddArray-Begins")
	defer list.Log("testAddArray-Ends")

	arrInterface := make([]interface{}, len(items))
	for i, v := range items {
		arrInterface[i] = v
	}
	list.AddArray(arrInterface)

	return list
}

func testAddValAtIndex(n int, index int, itemsToAdd ...interface{}) *ds.List {

	list := testAdd(n)
	list.Log("testAddValAtIndex generated items")

	for _, v := range itemsToAdd {
		list.AddVal(v, index)
	}

	list.Log("testAddValAtIndex... added " + strconv.Itoa(len(itemsToAdd)) + " items at index " + strconv.Itoa(index))

	return list
}

func testAddAll(n int, lst *ds.List) *ds.List {

	list := ds.NewList()
	defer list.Log("testAddAll after adding list!")
	for i := 0; i < n; i++ {
		list.Add(i)
	}
	list.Log("testAddAll after adding initial contents")

	list.AddAll(lst)

	return list
}

func testAddListAtIndex(n int, index int, lst *ds.List) {

	list := testAdd(n)
	list.Log("starting testAddListAtIndex")

	list.AddAllAt(index, lst)

	list.Log("testAddListAtIndex... after adding list at index " + strconv.Itoa(index))

}

func testLog(list *ds.List) {

	list.Log("testLog")

}
func testClear(list *ds.List) {

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

func testRemoveIndex(list *ds.List, index int) {
	list.Log("Remove from index: before")
	list.RemoveIndex(index)
	list.Log("Remove from index: after")
}

func testRemoveVal(list *ds.List, val interface{}) {
	list.Log("Remove value: before")
	list.Remove(val)
	list.Log("Remove value: after")
}

func testRemoveAll(list *ds.List) {

	lst := ds.NewList()
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
		val , err := list.Get(index)

		fmt.Printf("error: %v" , err)

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
	sublist , _ := list.SubList(3, 7)
	list.Log("MainList")
	subList.Log("SubList")
	sublist.Log("sublist")

	subList.Set(3 ,90000000)
	subList.Log("SubList...")
	list.Log("MainList")

	subList.Clear()

	list.Log("MainList")
	subList.Log("SubList")
	sublist.Log("sublist")


	subList.AddValues(2,9,8)

	list.Log("MainList")
	subList.Log("SubList")
	sublist.Log("sublist")

	fmt.Println(sublist.Get(0))



}



func main() {
	test1()

	/*	runSuite1()
		runSuite2()
		runSuite3()

		runSuite4()
	*/

	runSuite5()
	time.Sleep(time.Second * 4)

}
