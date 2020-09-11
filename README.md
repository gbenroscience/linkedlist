# linkedlist

This is a linkedlist implementation.
It may store any type of object.

For now it is not thread safe, so thread safety has to be added by the client code

You may watch this space for more updates as more work will probably be done on the list.



## Iteration:

To iterate through the list, do not use the list.Get(index) function as that runs in O(n) time and so will give you O(n<sup>2</sup>) performance.
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

