package ds

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

//AbstractList - An abstraction of a list
type AbstractList interface {
	ToArray() []interface{}
	addNode(elem *Node)
	Add(val interface{})
	addNodeAt(elem *Node, index int) bool
	AddVal(val interface{}, index int) bool
	removeNode(elem *Node) bool
	RemoveAll(lst *List)
	AddAll(lst *List)
	IsEmpty() bool
	AddAllAt(index int, lst *List)
	Remove(val interface{}) bool
	RemoveIndex(index int) bool
	SubList(startIndex int, endIndex int) *List
	getNode(index int) *Node
	Get(index int) interface{}
	getLastNode() *Node
	LastElement() interface{}
	IndexOf(val interface{}) int
	Contains(val interface{}) bool
	Clear()
	Log()
}

//Node - A list node
type Node struct {
	next *Node
	prev *Node
	val  interface{}
}

//List - The List
type List struct {
	Size      int
	firstNode *Node
	lastNode  *Node
	subList   *List

	parent *List

	//This indices are -1 if the sublist field is nil
	startIndex int
	endIndex   int

	iter *Node
}

func (list *List) Next() interface{} {

	if list.firstNode == nil || list.lastNode == nil{
		return nil
	}

	if list.firstNode != nil && list.iter == nil{
		list.iter = list.firstNode
		return list.iter.val
	}

	if list.firstNode == list.iter{
		list.iter = list.firstNode.next
		return list.iter.val
	}

	if list.iter != list.firstNode && list.iter != list.lastNode {
		list.iter = list.iter.next
		return list.iter.val
	}
	if list.lastNode == list.iter{
		return nil
	}

	list.reset()
	return nil
}
//Call this to reset the
func (list *List) reset(){
	if list.iter != nil{
		list.iter = nil
	}

}

func NewList() *List {
	list := new(List)
	list.Size = 0
	list.firstNode = nil
	list.lastNode = nil
	list.subList = nil
	list.parent = nil
	list.startIndex = -1
	list.endIndex = -1
	list.iter = nil
	return list
}
func (node *Node) initNode(prev *Node, val interface{}, next *Node) {
	node.prev = prev
	node.next = next
	node.val = val
}

//[]int{1,4,293,4,9}
//TESTED
func (list *List) ToArray() []interface{} {

	if list.isSubList() {

		result := make([]interface{}, list.Size)

		if len(result) == 0 {
			return result
		}

		i := 0
		for x := list.parent.getNode(list.startIndex); i < list.Size; x = x.next {
			result[i] = x.val
			i++
		}
		return result
	} else {

		result := make([]interface{}, list.Size)

		i := 0
		for x := list.firstNode; i < list.Size; x = x.next {
			result[i] = x.val
			i++
		}
		return result
	}

}

//TESTED
func (list *List) addNode(elem *Node) {

	if list.isSubList() {

		list.parent.addNodeAt(elem, list.endIndex)
		list.Size++
		list.endIndex++

	} else {

		oldLastNode := list.lastNode

		list.lastNode = elem
		if oldLastNode == nil {
			list.firstNode = elem
		} else {
			oldLastNode.next = elem
			elem.prev = oldLastNode
		}
		list.Size++

	}
}



//TESTED
func (list *List) AddValues(args ...interface{}) {
	if list.isSubList() {

		endNode := new(Node)
		skipFirst := false

		if list.IsEmpty() {
			if list.startIndex == 0 {
				endNode = list.parent.insertBefore(args[0], list.parent.firstNode)
				list.parent.Size++
				list.Size++
				list.endIndex++
				skipFirst = true
			} else {
				endNode = list.parent.getNode(list.startIndex - 1)
			}

		} else {
			endNode = list.getNode(list.Size - 1)
		}
		i := 0
		if skipFirst {
			i = 1
		}
		for ; i < len(args); i++ {
			endNode = list.parent.insertAfter(args[i], endNode)
			list.parent.Size++
			list.Size++
			list.endIndex++
		}

	} else {
		for i := 0; i < len(args); i++ {
			list.append(args[i])
		}
	}

}

//TESTED
func (list *List) AddArray(array []interface{}) {
	list.AddValues(array)
}

//TESTED
func (list *List) Add(val interface{}) {
	if list.isSubList() {

		if list.Size == 0 {
			list.parent.AddVal(val, list.startIndex)
		} else {
			list.parent.AddVal(val, list.endIndex)
		}
		//list.parent.AddVal(val, list.endIndex+1)
		list.Size++
		list.endIndex++
	} else {
		list.append(val)
	}

}

/**
 *
 *  Only parent lists should ever call this function!
 */
//TESTED
func (list *List) addNodeAt(elem *Node, index int) bool {

	if list.parent == nil {
		if index == list.Size {
			list.addNode(elem)
		} else {

			succ := list.getNode(index)

			// assert succ != nil;
			prev := succ.prev

			elem.prev = succ.prev
			elem.next = succ
			succ.prev = elem

			if prev == nil {
				list.firstNode = elem
			} else {
				prev.next = elem
			}
			list.Size++
		}
		list.syncAdditions(index)

		return true
	}
	return false

}

//TESTED
func (list *List) AddVal(val interface{}, index int) bool {

	if list.isSubList() {
		success := list.parent.AddVal(val, index+list.startIndex)
		list.Size++
		list.endIndex++
		return success

	} else {
		elem := new(Node)
		elem.next = nil
		elem.val = val

		return list.addNodeAt(elem, index)
	}

}

func (list *List) removeNode(elem *Node) bool {

	if list.isSubList() {

		if list.containsNode(elem) {

			succ := list.parent.removeNode(elem)
			list.Size--
			list.endIndex--
			return succ

		}
		return false
	} else {

		next := elem.next
		prev := elem.prev

		if prev == nil {
			list.firstNode = next
		} else {
			prev.next = next
			elem.prev = nil
		}

		if next == nil {
			list.lastNode = prev
		} else {
			next.prev = prev
			elem.next = nil
		}

		elem.val = nil
		list.Size--
		return true
	}

}

//TESTED
func (list *List) RemoveAll(lst *List) {
	//empty list
	if lst.firstNode == nil {
		return
	}

	if list.isSubList() {

		x := new(Node)
		if lst.isSubList() {
			x = lst.getNode(0)
		} else {
			x = lst.firstNode
		}

		for i := 0; i < lst.Size; i++ {
			list.Remove(x.val)
			x = x.next
		}

	} else {

		x := new(Node)
		if lst.isSubList() {
			x = lst.getNode(0)
		} else {
			x = lst.firstNode
		}
		for i := 0; i < lst.Size; i++ {
			list.Remove(x.val)
			x = x.next
		}

	}
}

//TESTED
func (list *List) AddAll(lst *List) error {

	//empty list

	if lst.isSubList() {
		if lst.Size == 0 {
			return errors.New("cant add empty sublist to a list")
		}
	} else {
		if lst.firstNode == nil {
			return errors.New("badly initialized list")
		}
	}

	return list.AddAllAt(list.Size, lst)

}

func (list *List) IsEmpty() bool {
	return list.Size == 0 && list.firstNode == nil
}

func (list *List) AddAllAt(index int, lst *List) error {

	if list.isSubList() {
		err := list.parent.AddAllAt(list.startIndex+index, lst)
		list.Size += lst.Size
		list.endIndex = list.startIndex + list.Size - 1
		return err
	} else {

		//empty list
		if lst.firstNode == nil || index >= list.Size || index < 0 {
			return errors.New("bad value for index or badly initialized list")
		}

		numNew := lst.Size
		if numNew == 0 {
			return errors.New("empty parameter sublist found")
		}
		prev := new(Node)
		succ := new(Node)
		if index == list.Size {
			succ = nil
			prev = list.lastNode
		} else {
			succ = list.getNode(index)
			prev = succ.prev
		}

		x := lst.firstNode

		for i := 0; i < numNew; i++ {
			e := x.val
			newNode := new(Node)
			newNode.initNode(prev, e, nil)
			if prev == nil {
				list.firstNode = newNode
			} else {
				prev.next = newNode
			}
			prev = newNode
			x = x.next
		}

		if succ == nil {
			list.lastNode = prev
		} else {
			prev.next = succ
			succ.prev = prev
		}

		list.Size += numNew
		return nil
	}
}

/**
 * Remove the first node that has
 * the same value as the parameter
 */
func (list *List) Remove(val interface{}) bool {

	if list.isSubList() {
		ind := list.IndexOf(val)
		if ind == -1 {
			return false
		} else {
			list.RemoveIndex(ind)

			list.Size--
			list.endIndex--
			return true
		}

	} else {
		//empty list
		if list.firstNode == nil {
			return false
		}

		x := list.firstNode
		for i := 0; i < list.Size; i++ {
			if x.val == val {
				succ := list.removeNode(x)

				if succ {
					list.syncRemoval(i)
				}

				return succ
			}

			x = x.next
		}

		return false
	}

}

func (list *List) RemoveIndex(index int) bool {

	if list.isSubList() {
		return list.parent.RemoveIndex(index + list.startIndex)
	} else {

		x := list.firstNode
		for i := 0; i <= index; i++ {
			if index == i {
				succ := list.removeNode(x)

				if succ {
					list.syncRemoval(index)
				}

				return succ
			}
			x = x.next

		}
		return false

	}

}

/**
 * Reflects changes in the parent list on the sublist when removals are made from the parent list
 */
func (list *List) syncRemoval(index int) {
	if list.subList != nil {

		if index <= list.subList.startIndex {
			list.subList.startIndex--
			list.subList.endIndex--
		} else if index > list.subList.startIndex && index <= list.subList.endIndex {
			list.subList.endIndex--
			list.subList.Size--
		}

	}
}

/**
 * Reflects changes in the parent list on the sublist when additions are made to the parent list
 */
func (list *List) syncAdditions(index int) {
	if list.subList != nil {

		if index <= list.subList.startIndex {
			list.subList.startIndex++
			list.subList.endIndex++
		} else if index > list.subList.startIndex && index <= list.subList.endIndex {
			list.subList.endIndex++
			list.subList.Size++
		}

	}
}

func (list *List) SubList(startIndex int, endIndex int) (*List, error) {

	if endIndex > list.Size {
		panic(strconv.Itoa(endIndex) + " > " + strconv.Itoa(list.Size) + " is not allowed")
	}

	if startIndex < 0 {
		return nil, errors.New(strconv.Itoa(startIndex) + " < 0 is not allowed")
	}
	if endIndex < 0 {
		return nil, errors.New(strconv.Itoa(endIndex) + " < 0 is not allowed")
	}
	if list.isSubList() {
		return nil, errors.New("sublist chaining not allowed. sublists of sublists cannot be made")
	}

	list.subList = new(List)
	list.subList = NewList()
	list.subList.firstNode = nil
	list.subList.lastNode = nil
	list.subList.parent = list

	list.subList.Size = endIndex - startIndex

	list.subList.startIndex = startIndex
	list.subList.endIndex = endIndex - 1

	return list.subList, nil

}

func (list *List) isSubList() bool {
	return list.parent != nil
}

/**
 * Returns the (non-nil) Node at the specified element index.
 */
func (list *List) getNode(index int) *Node {
	if index < 0 {
		panic("Index=(" + strconv.Itoa(index) + ") < 0 is not allowed")
	}

	if index > list.Size {
		panic("Index=(" + strconv.Itoa(index) + ") > list-Size=(" + strconv.Itoa(list.Size) + ") is not allowed")
	}

	if list.isSubList() {
		return list.parent.getNode(list.startIndex + index)
	}

	// NOTE x >> y is same as x ÷ 2^y
	if index < (list.Size >> 1) {
		x := list.firstNode
		for i := 0; i < index; i++ {
			x = x.next
		}
		return x
	} else {
		x := list.lastNode
		for i := list.Size - 1; i > index; i-- {
			x = x.prev
		}
		return x
	}
}

/**
 * Returns the (non-nil) Node at the specified element index.
 */
func (list *List) getBoundaryNodes() (*Node, *Node) {

	if list.isSubList() {

		start := list.startIndex
		end := list.endIndex

		first := list.getNode(0)
		last := new(Node)

		x := first

		i := start

		if list.Size > list.parent.Size-list.endIndex {

			last = list.getLastNode()

		} else {
			for ; i < end; i++ {
				x = x.next
			}
		}
		last = x
		return first, last
	} else {
		return list.firstNode, list.lastNode
	}

}

//Get - returns the element at that index in the list
func (list *List) Get(index int) interface{} {

	return list.getNode(index).val

}

func (list *List) getLastNode() *Node {
	if list.isSubList() {
		return list.parent.getNode(list.endIndex)
	} else {
		return list.lastNode
	}
}

func (list *List) LastElement() interface{} {
	return list.getLastNode().val
}

func (list *List) IndexOf(val interface{}) int {

	if list.isSubList() {
		startNode := list.getNode(0)

		i := 0
		for x := startNode; i < list.Size; x = x.next {

			if val == x.val {
				return i
			}
			i++
		}
		return -1

	} else {

		x := list.firstNode

		if x.val == val {
			return 0
		}

		for i := 0; i < list.Size; i++ {

			if val == x.val {
				return i
			}

			x = x.next
		}

		return -1
	}

}
func (list *List) indexOfNode(node *Node) int {

	if list.isSubList() {
		startNode := list.getNode(0)

		i := 0
		for x := startNode; i < list.Size; x = x.next {

			if node == x {
				return i
			}
			i++
		}
		return -1

	} else {

		x := list.firstNode

		if x == node {
			return 0
		}

		for i := 0; i < list.Size; i++ {

			if node == x {
				return i
			}

			x = x.next
		}

		return -1
	}

}
func (list *List) Contains(val interface{}) bool {
	return list.IndexOf(val) != -1

}
func (list *List) containsNode(node *Node) bool {
	return list.indexOfNode(node) != -1

}

/**
 * Links val as first element.
 */
func (list *List) prepend(val interface{}) {
	f := list.firstNode
	newNode := new(Node)
	newNode.initNode(nil, val, f)
	list.firstNode = newNode
	if f == nil {
		list.lastNode = newNode
	} else {
		f.prev = newNode
	}
	list.Size++
}

/**
 * Links val as last element.
 */
func (list *List) append(val interface{}) {
	l := list.lastNode
	newNode := new(Node)
	newNode.initNode(l, val, nil)
	list.lastNode = newNode
	if l == nil {
		list.firstNode = newNode
	} else {
		l.next = newNode
	}
	list.Size++
}

/**
 * Returns the (non-NULL) Node at the specified element index.
 * The bounds parameter is a pointer to an array that will hold the
 * pointers to the boundary nodes
 */
func (list *List) getNodesAt(start int, end int) []*Node {




	if start == end || start > end || start < 0 || end < 0 {
		return []*Node{}
	}

	var bounds = make([]*Node , 0)

	first := list.getNode(start)

	last := new(Node)

	var x = first

	i := start

	//The end node is closer to the already found start node than the end node is close to the end of the list.
	if end-start < list.Size-end {

		for ; i < end; i++ {
			x = x.next
		}

		last = x
	} else {
		len := list.Size

		y := list.lastNode
		for i = len - 1; i >= end; i-- {
			y = y.prev
		}
		last = y

	}

	bounds = append(bounds, first)
	bounds = append(bounds, last)

	return bounds

}

/**
 * Inserts element e before non-null Node succ.
 * ONLY TO BE CALLED BY BASE PARENT LISTS
 * Return a pointer to the new node that was inserted.
 * This will help with spontaneous insertions
 */
func (list *List) insertBefore(e interface{}, succ *Node) *Node {

	prev := succ.prev

	newNode := new(Node)
	newNode.initNode(prev, e, succ)

	succ.prev = newNode
	if prev == nil {
		list.firstNode = newNode
	} else {
		prev.next = newNode
	}
	list.Size++
	return newNode
}

/**
 * Inserts element e after non-null Node succ.
 *
 * ONLY TO BE CALLED BY BASE PARENT LISTS
 * Return a pointer to the new node that was inserted.
 * This will help with spontaneous insertions
 */
func (list *List) insertAfter(e interface{}, succ *Node) *Node {

	next := succ.next

	newNode := new(Node)
	newNode.initNode(succ, e, succ.next)

	succ.next = newNode
	if next == nil {
		list.lastNode = newNode
	} else {
		next.prev = newNode
	}
	list.Size++

	return newNode
}
func (list *List) Clear() {

	if list.isSubList() {

		begin, end := list.getBoundaryNodes()

		//The node just before this sublist's first node
		leftBound := begin.prev

		//The node just after this range's last node
		rightBound := end.next

		if leftBound != nil {
			leftBound.next = rightBound
		} else {
			list.parent.firstNode = rightBound
		}

		if rightBound != nil {
			rightBound.prev = leftBound
		} else {
			list.parent.lastNode = leftBound
		}

		list.removeLinkedRange(begin, end)

		list.parent.Size -= list.Size
		list.endIndex = list.startIndex
		list.Size = 0

	} else {
		x := list.firstNode

		for x != nil {
			next := x.next
			x.val = nil
			x.next = nil
			x.prev = nil
			x = next
		}
		list.firstNode = nil
		list.lastNode = nil
		list.Size = 0

		if list.subList != nil {
			list.subList.startIndex = 0
			list.subList.endIndex = 0
			list.subList.Size = 0
		}

	}

}
func (list *List) removeLinkedRange(startNode *Node, stopNode *Node) {

	for x := startNode; x != stopNode; x = x.next {

		x.val = nil
		if x != startNode {
			x.prev = nil
		}

	}

}
func (list *List) Log(optionalLabel string) {

	if list.isSubList() {

		if list.Size == 0 {
			fmt.Println(optionalLabel, ":\n[], len: 0")
			return
		}

		appender := optionalLabel + ":\n ["

		i := 0

		for x := list.getNode(i); i < list.Size; x = x.next {
			dType := reflect.TypeOf(x.val).Kind()
			if dType == reflect.String {
				appender += x.val.(string) + ","
			} else if dType == reflect.Int {
				appender += strconv.Itoa(x.val.(int)) + ","
			} else if dType == reflect.Int8 {
				appender += strconv.Itoa(int(x.val.(int8))) + ","
			} else if dType == reflect.Int16 {
				appender += strconv.Itoa(int(x.val.(int16))) + ","
			} else if dType == reflect.Int32 {
				appender += strconv.Itoa(int(x.val.(int32))) + ","
			} else if dType == reflect.Int64 {
				appender += strconv.Itoa(int(x.val.(int64))) + ","
			}

			i++
		}
		appender += "], len: " + strconv.Itoa(list.Size) + ", startIndex: " + strconv.Itoa(list.startIndex) + ", endIndex: " + strconv.Itoa(list.endIndex)
		fmt.Println(appender)

	} else {

		if list.firstNode == nil {
			fmt.Println(optionalLabel + ":\n[], len: 0")
			return
		}
		counter := 0
		currentNode := list.firstNode
		dType := reflect.TypeOf(currentNode.val).Kind()

		appender := optionalLabel + ":\n"
		if dType == reflect.String {
			appender += "[" + currentNode.val.(string) + ","
		} else if dType == reflect.Int {
			appender += "[" + strconv.Itoa(currentNode.val.(int)) + ","
		} else if dType == reflect.Int8 {
			appender += "[" + strconv.Itoa(int(currentNode.val.(int8))) + ","
		} else if dType == reflect.Int16 {
			appender += "[" + strconv.Itoa(int(currentNode.val.(int16))) + ","
		} else if dType == reflect.Int32 {
			appender += "[" + strconv.Itoa(int(currentNode.val.(int32))) + ","
		} else if dType == reflect.Int64 {
			appender += "[" + strconv.Itoa(int(currentNode.val.(int64))) + ","
		}

		for currentNode.next != nil && counter < list.Size {
			dType = reflect.TypeOf(currentNode.next.val).Kind()
			if dType == reflect.String {
				appender += currentNode.next.val.(string) + ","
			} else if dType == reflect.Int {
				appender += strconv.Itoa(currentNode.next.val.(int)) + ","
			} else if dType == reflect.Int8 {
				appender += strconv.Itoa(int(currentNode.next.val.(int8))) + ","
			} else if dType == reflect.Int16 {
				appender += strconv.Itoa(int(currentNode.next.val.(int16))) + ","
			} else if dType == reflect.Int32 {
				appender += strconv.Itoa(int(currentNode.next.val.(int32))) + ","
			} else if dType == reflect.Int64 {
				appender += strconv.Itoa(int(currentNode.next.val.(int64))) + ","
			}

			currentNode = currentNode.next
			counter++
		}
		appender += "], len: " + strconv.Itoa(list.Size)
		fmt.Println(appender)

	}

}
