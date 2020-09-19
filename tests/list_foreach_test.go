package tests

import (
	"fmt"
	"github.com/gbenroscience/linkedlist/ds"
	"testing"
	"time"
)

func  ForEachTest(t *testing.T)  {


	list := ds.NewList()

	add(list , 1000)



	go list.ForEach(func(val interface{}) {
		fmt.Printf("val-1: %+v",val)
	})

	go list.ForEach(func(val interface{}) {
		fmt.Printf("val-2: %+v",val)
	})





	time.Sleep(time.Second * 10)

	fmt.Println("Tests done")






}

func add(list *ds.List, itemCount int){
	for i:=0;i<itemCount;i++{
		list.Add(i)
	}
}