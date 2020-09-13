package ds

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"sync"
)

const (
	ChanSize = 1000
	Shutdown = 1
)

func (list *List) process() {

	defer func() {
		close(list.writeChan)
	}()
	for {
		select {

		case x := <-list.writeChan:
			if x.shutdown == Shutdown{
				return
			}
			x.callback()
			break
		}
	}

}

//AbstractList - An abstraction of a list
type AbstractList interface {
	Add(val interface{})
	AddVal(val interface{}, index int) bool
	AddAll(lst *List)
	AddAllAt(index int, lst *List)

	Remove(val interface{}) bool
	RemoveIndex(index int) bool
	RemoveAll(lst *List)
	Clear()

	Get(index int) interface{}
	ToArray() []interface{}
	LastElement() interface{}

	SubList(startIndex int, endIndex int) *List

	IsEmpty() bool
	Contains(val interface{}) bool

	IndexOf(val interface{}) int

	Log()
}

//Node - A list node
type Node struct {
	next *Node
	prev *Node
	val  interface{}
}

type DataTypePair struct {
	callback func()
	shutdown int
}

//List - The List
type List struct {
	size      int
	firstNode *Node
	lastNode  *Node
	subList   *List
	parent    *List
	//This indices MUST BE -1 if the sublist field is nil
	startIndex int
	endIndex   int
	mu sync.Mutex
	//Used for rapid iteration over the list
	iter *Node
	//write data into list...used by the Add methods
	writeChan chan DataTypePair
}

func NewList() *List {
	list := new(List)
	list.size = 0
	list.firstNode = nil
	list.lastNode = nil
	list.subList = nil
	list.parent = nil
	list.startIndex = -1
	list.endIndex = -1
	list.iter = nil
	list.mu = sync.Mutex{}
	list.writeChan = make(chan DataTypePair, ChanSize)

	go list.process()
	return list
}
func (node *Node) initNode(prev *Node, val interface{}, next *Node) {
	node.prev = prev
	node.next = next
	node.val = val
}

func (list *List) Next() interface{} {

	if list.iter != nil {
		if list.iter.next != nil {
			list.iter = list.iter.next
			return list.iter.val
		}
		list.resetIterator()
		return nil
	} else {
		if list.firstNode != nil {
			list.iter = list.firstNode
			return list.Next()
		}
		return nil
	}

}

//Call this to reset the
func (list *List) resetIterator() {
	if list.iter != nil {
		list.iter = nil
	}

}

func (list *List) ForEach(function func(val interface{})) {

	defer list.mu.Unlock()
	list.mu.Lock()
	var x interface{}
	for ; ; {
		x = list.Next()
		if x == nil {
			break
		}
		function(x)
	}

}

//[]int{1,4,293,4,9}
//TESTED
func (list *List) ToArray() []interface{} {

	if list.isSubList() {

		result := make([]interface{}, list.size)

		if len(result) == 0 {
			return result
		}

		i := 0
		for x := list.parent.getNode(list.startIndex); i < list.size; x = x.next {
			result[i] = x.val
			i++
		}
		return result
	} else {

		result := make([]interface{}, list.size)

		i := 0
		for x := list.firstNode; i < list.size; x = x.next {
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
		list.size++
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
		list.size++

	}
}

//TESTED
func (list *List) Add(val interface{}) {
	list.writeChan <- DataTypePair{
		callback: func() {
			list.add(val)
		},
	}
}

//TESTED
func (list *List) add(val interface{}) {

	if list.isSubList() {

		if list.size == 0 {
			list.parent.addVal(val, list.startIndex)
		} else {
			list.parent.addVal(val, list.endIndex)
		}
		//list.parent.addVal(val, list.endIndex+1)
		list.size++
		list.endIndex++
	} else {
		list.append(val)
	}

}

func (list *List) AddVal(val interface{}, index int) bool {

	args := make([]interface{}, 0)
	args = append(args, val, index)
	list.writeChan <- DataTypePair{
		callback: func() {
			list.addVal(val, index)
		},
	}

	return true
}

//TESTED
func (list *List) addVal(val interface{}, index int) bool {

	if list.isSubList() {
		success := list.parent.addVal(val, index+list.startIndex)
		list.size++
		list.endIndex++
		return success

	} else {
		elem := new(Node)
		elem.next = nil
		elem.val = val

		return list.addNodeAt(elem, index)
	}

}

func (list *List) AddValues(args ...interface{}) {
	list.writeChan <- DataTypePair{
		callback: func() {
			list.addValues(args)
		},
	}
}

//TESTED
func (list *List) addValues(args ...interface{}) {

	if list.isSubList() {

		endNode := new(Node)
		skipFirst := false

		if list.IsEmpty() {
			if list.startIndex == 0 {
				endNode = list.parent.insertBefore(args[0], list.parent.firstNode)
				list.parent.size++
				list.size++
				list.endIndex++
				skipFirst = true
			} else {
				endNode = list.parent.getNode(list.startIndex - 1)
			}

		} else {
			endNode = list.getNode(list.size - 1)
		}
		i := 0
		if skipFirst {
			i = 1
		}
		for ; i < len(args); i++ {
			endNode = list.parent.insertAfter(args[i], endNode)
			list.parent.size++
			list.size++
			list.endIndex++
		}

	} else {
		for i := 0; i < len(args); i++ {
			list.append(args[i])
		}
	}

}

func (list *List) AddArray(array []interface{}) {
	list.writeChan <- DataTypePair{
		callback: func() {
			list.addArray(array)
		},
	}
}

//TESTED
func (list *List) addArray(array []interface{}) {
	list.addValues(array)
}

func (list *List) AddAll(lst *List) bool {

	list.writeChan <- DataTypePair{
		callback: func() {
			list.addAll(lst)
		},
	}

	return true
}

//TESTED
func (list *List) addAll(lst *List) error {

	//empty list

	if lst.isSubList() {
		if lst.size == 0 {
			return errors.New("cant add empty sublist to a list")
		}
	} else {
		if lst.firstNode == nil {
			return errors.New("badly initialized list")
		}
	}

	return list.addAllAt(list.size, lst)

}

func (list *List) AddAllAt(index int, lst *List) bool {
	list.writeChan <- DataTypePair{
		callback: func() {
			list.addAllAt(index, lst)
		},
	}
	return true
}

func (list *List) addAllAt(index int, lst *List) error {

	if list.isSubList() {
		err := list.parent.addAllAt(list.startIndex+index, lst)
		list.size += lst.size
		list.endIndex = list.startIndex + list.size - 1
		return err
	} else {

		//empty list
		if lst.firstNode == nil || index >= list.size || index < 0 {
			return errors.New("bad value for index or badly initialized list")
		}

		numNew := lst.size
		if numNew == 0 {
			return errors.New("empty parameter sublist found")
		}
		prev := new(Node)
		succ := new(Node)
		if index == list.size {
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

		list.size += numNew
		return nil
	}
}

/**
 *
 *  Only parent lists should ever call this function!
 */
//TESTED
func (list *List) addNodeAt(elem *Node, index int) bool {

	if list.parent == nil {
		if index == list.size {
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
			list.size++
		}
		list.syncAdditions(index)

		return true
	}
	return false

}

func (list *List) removeNode(elem *Node) bool {

	if list.isSubList() {

		if list.containsNode(elem) {

			succ := list.parent.removeNode(elem)
			list.size--
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
		list.size--
		return true
	}

}
func (list *List) Remove(val interface{}) bool {

	list.writeChan <- DataTypePair{
		callback: func() {
			list.remove(val)
		},
	}
	return true
}

/**
 * Remove the first node that has
 * the same value as the parameter
 */
func (list *List) remove(val interface{}) bool {

	if list.isSubList() {
		ind := list.IndexOf(val)
		if ind == -1 {
			return false
		} else {
			list.removeIndex(ind)

			list.size--
			list.endIndex--
			return true
		}

	} else {
		//empty list
		if list.firstNode == nil {
			return false
		}

		x := list.firstNode
		for i := 0; i < list.size; i++ {
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

	list.writeChan <- DataTypePair{
		callback: func() {
			list.removeIndex(index)
		},
	}
	return true
}

func (list *List) removeIndex(index int) bool {

	if list.isSubList() {
		return list.parent.removeIndex(index + list.startIndex)
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

func (list *List) RemoveAll(lst *List) bool {

	list.writeChan <- DataTypePair{
		callback: func() {
			list.removeAll(lst)
		},
	}
	return true
}

//TESTED
func (list *List) removeAll(lst *List) {

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

		for i := 0; i < lst.size; i++ {
			list.remove(x.val)
			x = x.next
		}

	} else {

		x := new(Node)
		if lst.isSubList() {
			x = lst.getNode(0)
		} else {
			x = lst.firstNode
		}
		for i := 0; i < lst.size; i++ {
			list.remove(x.val)
			x = x.next
		}

	}
}

func (list *List) IsEmpty() bool {
	return list.size == 0 && list.firstNode == nil
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
			list.subList.size--
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
			list.subList.size++
		}

	}
}

func (list *List) SubList(startIndex int, endIndex int) (*List, error) {
	defer list.mu.Unlock()
	list.mu.Lock()

	if endIndex > list.size {
		panic(strconv.Itoa(endIndex) + " > " + strconv.Itoa(list.size) + " is not allowed")
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

	list.subList.size = endIndex - startIndex

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

	if index > list.size {
		panic("Index=(" + strconv.Itoa(index) + ") > list-size=(" + strconv.Itoa(list.size) + ") is not allowed")
	}

	if list.isSubList() {
		return list.parent.getNode(list.startIndex + index)
	}

	// NOTE x >> y is same as x ÷ 2^y
	if index < (list.size >> 1) {
		x := list.firstNode
		for i := 0; i < index; i++ {
			x = x.next
		}
		return x
	} else {
		x := list.lastNode
		for i := list.size - 1; i > index; i-- {
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

		if list.size > list.parent.size-list.endIndex {

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
	defer list.mu.Unlock()
	list.mu.Lock()
	return list.getLastNode().val
}

func (list *List) Contains(val interface{}) bool {
	return list.IndexOf(val) != -1

}

func (list *List) IndexOf(val interface{}) int {
	defer list.mu.Unlock()
	list.mu.Lock()
	if list.isSubList() {
		startNode := list.getNode(0)

		i := 0
		for x := startNode; i < list.size; x = x.next {

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

		for i := 0; i < list.size; i++ {

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
		for x := startNode; i < list.size; x = x.next {

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

		for i := 0; i < list.size; i++ {

			if node == x {
				return i
			}

			x = x.next
		}

		return -1
	}

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
	list.size++
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
	list.size++
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

	var bounds = make([]*Node, 0)

	first := list.getNode(start)

	last := new(Node)

	var x = first

	i := start

	//The end node is closer to the already found start node than the end node is close to the end of the list.
	if end-start < list.size-end {

		for ; i < end; i++ {
			x = x.next
		}

		last = x
	} else {
		len := list.size

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
	list.size++
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
	list.size++

	return newNode
}

func (list *List) Clear() bool {

	list.writeChan <- DataTypePair{
		callback: func() {
			list.clear()
		},
	}
	return true
}

func (list *List) clear() {

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

		list.parent.size -= list.size
		list.endIndex = list.startIndex
		list.size = 0

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
		list.size = 0

		if list.subList != nil {
			list.subList.startIndex = 0
			list.subList.endIndex = 0
			list.subList.size = 0
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

func (list *List) Log(optionalLabel string) bool {

	list.writeChan <- DataTypePair{
		callback: func() {
			list.log(optionalLabel)
		},
	}

	return true
}

func (list *List) log(optionalLabel string) {

	if list.isSubList() {

		if list.size == 0 {
			fmt.Println(optionalLabel, ":\n[], len: 0")
			return
		}

		appender := optionalLabel + ":\n ["

		i := 1

		for x := list.getNode(i); i < list.size; x = x.next {
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

		appender += "], len: " + strconv.Itoa(list.size) + " confirm-len(" + strconv.Itoa(i) + "), startIndex: " + strconv.Itoa(list.startIndex) + ", endIndex: " + strconv.Itoa(list.endIndex)
		fmt.Println(appender)

	} else {

		if list.firstNode == nil {
			fmt.Println(optionalLabel + ":\n[], len: 0")
			return
		}
		counter := 1
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

		for currentNode.next != nil && counter < list.size {
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
		appender += "], len: " + strconv.Itoa(list.size) + " , confirm-len(" + strconv.Itoa(counter) + ")"
		fmt.Println(appender)

	}

}

func (list *List) Count() int {
	return list.size
}

func (list *List) Close() error{
	list.writeChan <- DataTypePair{
		callback: nil,
		shutdown: Shutdown,
	}

	return nil
}
