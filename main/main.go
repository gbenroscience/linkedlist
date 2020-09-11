package main

import (
	"bytes"
	"fmt"
	"github.com/gbenroscience/linkedlist/ds"
	"math"
	"math/rand"
	"strconv"
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

	fmt.Printf("Approx CubeRoot: %17d Real CubeRoot: %17d\n", root, math.Cbrt(x))

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
func Print(x interface{}){
	a += 1
	fmt.Printf("Printing list: found %d\n ", x)
}

var(
	a = 0
)
func main() {

	list := ds.NewList()

	for i := 0; i < 1000; i++ {
		list.Add(i)
	}

	fmt.Printf("list now has %d elements\n", list.Size)


/*
	var x interface{}

	for ; ; {
		x = list.Next()
		if x == nil {
			fmt.Printf("Printing list: found ??? %v \n", x)
			break
		}
		fmt.Printf("Printing list: found %d\n ", x)
	}*/

 list.ForEach(Print)

 fmt.Printf("a = %d\n" , a)


	fmt.Println("...............................................................................................................")
	list.Log("Checking...")

}
