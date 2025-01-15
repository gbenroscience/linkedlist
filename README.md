# linkedlist

This is a thread safe linkedlist implementation for Golang.
It may store any type of object.

 
## Features

1. Thread safety in concurrent access is ensured by using mutexes.

2. Allows greater manipulation using sublists. Sublists are backed by the parent list, so you can manipulate portions of the list as though they were a list!

3. Has several methods that allow the user manipulate the list with ease.


We now have 3 different implementations of the `linkedlist`.
The old `List` is now called a `CList` which implies it is left for backwards compatibility with old Golang code that still uses interfaces instead of generics.

We now have 2 list implementations which take advantage of the shiny new `generics` feature available in newer versions of golang.

One of them implements generics using `comparable` and the other uses the broader `any` interface.

We will now look at initializing the 3 types of `linkedlist`.

### Initializing the old `CList`

All that has changed here is the name of the struct, and its initializing function's name.

To initialize it, do:
```Go
list := ds.NewCList()
```

### Initializing the new `List[T comparable]`

The name of the struct is `List`

To initialize it, do:
```Go
list := ds.NewList()
```


### Initializing the new `AnyList[T any]`

The name of the struct is `AnyList`

To initialize it, do:
```Go
list := ds.NewAnyList()
```

In order to get the best runtime speed from an AnyList instance, make sure you override the `Equals` function(which is a field of the `AnyList` struct). We have added a generic `Equals` function which will run very slow(and perhaps may be incorrect for some cases). So, whenever you calll

```Go
list := ds.NewAnyList()
```
The next line should be:
```Go
list.Equals = func(val1 T, val2 T) bool {
	return //code that returns true when val1 is same as val2
	}
```

This is not necessary for the instance of
```Go 
List[T comparable]
```


**All 3 implementations support the SubList functionality.**

Sublists behave like normal lists too, presenting a view of portions of the list.
e.g.

```Go

package main

func testAdd(n int) *ds.List {

	list := ds.NewList[int]()
	for i := 0; i < n; i++ {
		list.Add(i)
	}

	return list
}

func main(){

list := testAdd(10)//adds 50 ints (from 0 to 9) to the list e.g. [0,1,2,3,4,5,6,7,8,9]

subList, err := list.SubList(2, 8)
if err != nil{
//handle error here
}
//The sublist here now contains [2,3,4,5,6,7]

subList.Clear()
//The Clear command empties the sublist and also clears the portion of the list occupied by the sublist. The list now contains: [0,1,8,9]


}

```

Changes made to the sublist (add , remove, clear, update) are reflected in the parent list.
If you clear the sublist, it becomes detached from its parent.
Changes made to the sublist are no longer propagated to the sublist

If you need to have a sublist of a list independent of the original list, then create the sublist as above and call the Clone method on it e.g:


```Go
freeSubList := subList.Clone()
```

If the <code>freeSubList</code> above is modified, the changes no longer reflect on the main list.


1. Allows quick iteration using the <b>ForEach</b> function

The old way to iterate over the list was:

 ```Go
	for ; ; {
		x = list.Next()
		if x == nil {
			fmt.Printf("Printing list: found ??? %v ", x)
			break
		}
		fmt.Printf("Printing list: found %d\n ", x)
	}
 ```

But we have removed the list.Next() method. The standard way to iterate over the list now is to use the <b>ForEach</b> function.

```Go
  func (list *List) ForEach(function func(val interface{}))
```

### More on Iteration:

To iterate through the list, do not use the list.Get(index) function in a loop as that runs in O(n) time and so will give you O(n<sup>2</sup>) performance.
Instead , say you created the list like this:

```Go
list := ds.NewList()
 ```


And then you populated it like this:

```Go
for i:=0; i<1000;i++{
list.Add(i)
}

fmt.Printf("list now has %d elements\n" , list.Size))
 ```
 Then iterate over it like this:
 
 ```Go
	list.ForEach(func(x interface{}) bool{
	   //do stuff
	   
	   return true// if you want the loop to continue to the end, return true. If you want the function to break out, return false.
	})
 ```


This allows you fetch consecutive items in the list in constant time and so the traversal over the entire list is done in O(n) time

Please note that the iterator is a cyclic one.
If it detects the end of the list, it resets, which allows you to break out of the loop, but once you repeat that loop, it starts all over again from the
begin. So you can iterate repeatedly over the same list


## Using the ForEach function For Iteration

This function is defined as:

```Go
  func (list *List) ForEach(function func(val interface{}) bool)
```

An example would be:

```Go
func Print(x interface{}) bool{
	fmt.Printf("Printing list: found %v\n ", x)
	return true
}

 list.ForEach(Print)

```
The list will iterate over every element in it and call the function on each of them (e.g Print or whatever) 

Alternatively of course, you may do:
```Go
 list.ForEach(func(x interface{})bool {
   fmt.Printf("Printing list: found %v\n ", x)
   return true
 })
```
