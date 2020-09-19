# linkedlist

This is a thread safe linkedlist implementation for Golang.
It may store any type of object.

 
## Features



1. Thread safety in concurrent access is ensured by using mutexes.
Thread safety using channels has been removed due to various issues with it.

2. Allows greater manipulation using sublists. You can manipulate portions of the list as though they were a list!
I would limit this to the basest operations and also, DO NOT create sublists of a sublist.

3. Allows quick iteration using the <b>ForEach</b> function

The idiomatic way to iterate over the list is:

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

But we have added an even more convenient and standard way using the <b>ForEach</b> function

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
	for ; ; {
		x = list.Next()
		if x == nil {
			fmt.Printf("Printing list: found ??? %v ", x)
			break
		}
		fmt.Printf("Printing list: found %d\n ", x)
	}
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
func Print(x interface{}){
	fmt.Printf("Printing list: found %d\n ", x)
	return true
}

 list.ForEach(Print)

```
The list will iterate over every element in it and call the function on each of them (e.g Print or whatever) 

Alternatively of course, you may do:
```Go
 list.ForEach(func(x interface{}) {
   fmt.Printf("Printing list: found %d\n ", x)	 
 })
```
