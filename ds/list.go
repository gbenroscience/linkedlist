package ds

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// abstractList - An abstraction of a list
type abstractList[T comparable] interface {
	Add(val T)
	AddVal(val T, index int) bool
	AddAll(lst *List[T]) bool
	AddAllAt(index int, lst *List[T]) bool

	Remove(val T) bool
	RemoveIndex(index int) bool
	RemoveAll(lst *List[T]) *List[T]
	Clear() bool

	Set(index int, val T)
	Get(index int) T
	ToArray() []T
	LastElement() T

	SubList(startIndex int, endIndex int) (*List[T], error)

	IsEmpty() bool
	Contains(val T) bool

	IndexOf(val T) int

	Log(optionalLabel string)
}

// lNode - A list node
type lNode[T comparable] struct {
	next *lNode[T]
	prev *lNode[T]
	val  T
}

// List - The List
type List[T comparable] struct {
	size      int
	firstNode *lNode[T]
	lastNode  *lNode[T]
	parent    *List[T]
	parenLen  int
	mu        sync.Mutex
	//Used for rapid iteration over the list's values
	iter *lNode[T]
	//Used for rapid iteration over the list's nodes
	nodeIter *lNode[T]
}

func NewList[T comparable]() *List[T] {
	list := new(List[T])

	list.size = 0
	list.parent = nil
	list.firstNode = nil
	list.lastNode = nil
	list.iter = nil
	list.nodeIter = nil
	list.mu = sync.Mutex{}

	return list
}

func initNode[T comparable](prev *lNode[T], val T, next *lNode[T]) *lNode[T] {
	node := new(lNode[T])
	node.prev = prev
	node.next = next
	node.val = val

	return node
}

func (node *lNode[T]) isNilValOnNode() bool {
	return new(T) == &node.val
}

func (list *List[T]) nextNode() *lNode[T] {

	if list.nodeIter == nil {
		if list.firstNode == nil {
			return nil
		}
		list.nodeIter = list.firstNode
		return list.nodeIter
	} else {
		if list.nodeIter == list.lastNode {
			return nil
		}
		if list.nodeIter.next != nil {
			list.nodeIter = list.nodeIter.next
			return list.nodeIter
		}
		return nil
	}

}

// Call this to reset the nodes iterator
func (list *List[T]) resetNodeIterator() {
	if list.nodeIter != nil {
		list.nodeIter = nil
	}
}

func (list *List[T]) next() *T {

	var next = list.nextNode()
	if next != nil {
		return &next.val
	}

	return nil
}

// Call this to reset the values iterator
func (list *List[T]) resetIterator() {
	if list.iter != nil {
		list.iter = nil
	}
}

func (list *List[T]) ForEach(function func(val T) bool) {

	x := new(T)

	defer list.mu.Unlock()
	list.mu.Lock()
	list.resetIterator()

	for {
		x = list.next()
		if x == nil {
			break
		}
		if !function(*x) {
			break
		}

	}
}

func (list *List[T]) forEachNodeFrom(start *lNode[T], function func(node *lNode[T]) bool) {

	var x *lNode[T]
	list.resetNodeIterator()
	list.nodeIter = start

	for {
		x = list.nextNode()
		if x == nil {
			break
		}
		if !function(x) {
			break
		}

	}
}

func (list *List[T]) forEachNode(function func(node *lNode[T]) bool) {

	var x *lNode[T]
	list.resetNodeIterator()

	for {
		x = list.nextNode()
		if x == nil {
			break
		}
		if !function(x) {
			break
		}

	}
}

// TESTED
func (list *List[T]) ToArray() []T {

	result := make([]T, list.count())

	i := 0
	list.ForEach(func(x T) bool {
		result[i] = x
		i++
		return true
	})

	return result
}

// TESTED
func (list *List[T]) addNode(elem *lNode[T]) {
	oldLastNode := list.lastNode

	list.lastNode = elem
	if oldLastNode == nil {
		list.firstNode = elem
	} else {
		oldLastNode.next = elem
		elem.prev = oldLastNode
	}
	list.incrementSize(1)
}

// TESTED
func (list *List[T]) Add(val T) {
	defer list.mu.Unlock()

	list.mu.Lock()
	list.add(val)

}

// TESTED
func (list *List[T]) add(val T) {
	list.append(val)
}

func (list *List[T]) AddVal(val T, index int) bool {
	defer list.mu.Unlock()

	list.mu.Lock()
	_, _ = list.addVal(val, index)

	return true
}

// TESTED
func (list *List[T]) addVal(val T, index int) (bool, error) {
	node := initNode(nil, val, nil)
	return list.addNodeAt(node, index)

}

func (list *List[T]) AddValues(args ...T) {
	defer list.mu.Unlock()

	list.mu.Lock()
	list.addValues(args...)

}

// TESTED
func (list *List[T]) addValues(args ...T) {
	for _, v := range args {
		list.append(v)
	}
}

func (list *List[T]) AddArray(array []T) {
	defer list.mu.Unlock()

	list.mu.Lock()
	list.addArray(array)
}

// TESTED
func (list *List[T]) addArray(array []T) {
	list.addValues(array...)
}

func (list *List[T]) AddAll(lst *List[T]) bool {
	defer list.mu.Unlock()

	list.mu.Lock()
	_ = list.addAll(lst)

	return true
}

// TESTED
func (list *List[T]) addAll(lst *List[T]) error {
	err := list.addAllAt(list.count(), lst)
	return err
}

func (list *List[T]) AddAllAt(index int, lst *List[T]) bool {
	defer list.mu.Unlock()

	list.mu.Lock()
	_ = list.addAllAt(index, lst)

	return true
}

func (list *List[T]) Clone() *List[T] {

	ls := new(List[T])

	list.forEachNode(func(node *lNode[T]) bool {
		ls.Add(node.val)
		return true
	})

	return ls
}

func (list *List[T]) addAllAt(index int, lst *List[T]) error {

	sz := list.count()
	//empty list
	if lst.firstNode == nil || index > sz || index < 0 {
		return errors.New("bad value for index or badly initialized list")
	}
	dup := lst.Clone()

	numNew := dup.count()
	if numNew == 0 {
		return errors.New("empty parameter sublist found")
	}

	if index == sz {

		lastNode := list.lastNode

		lastNode.next = dup.firstNode
		dup.firstNode.prev = lastNode
		list.lastNode = dup.lastNode

	} else if index == 0 {
		oldFirstNode := list.firstNode

		dup.lastNode.next = oldFirstNode
		oldFirstNode.prev = dup.lastNode
		list.firstNode = dup.firstNode
	} else {
		nodeAtIndex, err := list.getNode(index)
		if err != nil {
			return err
		}

		prevNode := nodeAtIndex.prev

		prevNode.next = dup.firstNode
		dup.firstNode.prev = prevNode

		dup.lastNode.next = nodeAtIndex
		nodeAtIndex.prev = dup.lastNode

	}

	list.incrementSize(numNew)
	return nil
}

/**
 *
 *  Only parent lists should ever call this function!
 */
//TESTED
func (list *List[T]) addNodeAt(elem *lNode[T], index int) (bool, error) {

	sz := list.count()
	if index >= 0 && index <= sz {

		if index == sz {
			list.addNode(elem)
		} else {

			succ, err := list.getNode(index)

			if err != nil {
				return false, err
			}

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
			list.incrementSize(1)
		}

		return true, nil

	} else {
		return false, errors.New("index must lie between 0 and " + strconv.Itoa(list.count()))
	}

}

func (list *List[T]) incrementSize(dx int) {
	list.size += dx
	if list.parent != nil {
		list.parent.size += dx
	}
}

func (list *List[T]) decrementSize(dx int) {
	list.size -= dx
	if list.parent != nil {
		list.parent.size -= dx
	}
}

func (list *List[T]) removeNode(elem *lNode[T]) bool {
	var nilVal T
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

	elem.val = nilVal
	list.decrementSize(1)
	return true
}

func (list *List[T]) removeIndex(index int) bool {

	x := list.firstNode
	for i := 0; i <= index; i++ {
		if index == i {
			succ := list.removeNode(x)
			return succ
		}
		x = x.next
	}
	return false
}

func (list *List[T]) Remove(val T) bool {

	defer list.mu.Unlock()
	list.mu.Lock()
	list.remove(val)

	return true
}

/**
 * Remove the first node that has
 * the same value as the parameter
 */
func (list *List[T]) remove(val T) bool {

	//empty list
	if list.firstNode == nil {
		return false
	}

	x := list.firstNode
	sz := list.count()
	for i := 0; i < sz; i++ {
		if x.val == val {
			succ := list.removeNode(x)

			return succ
		}

		x = x.next
	}

	return false

}

func (list *List[T]) RemoveIndex(index int) bool {
	defer list.mu.Unlock()
	list.mu.Lock()
	return list.removeIndex(index)
}

func (list *List[T]) RemoveAll(lst *List[T]) bool {
	defer list.mu.Unlock()
	list.mu.Lock()
	list.removeAll(lst)

	return true
}

// TESTED
func (list *List[T]) removeAll(lst *List[T]) {

	//empty list
	if lst.firstNode == nil {
		return
	}

	x := lst.firstNode

	sz := lst.count()

	for i := 0; i < sz; i++ {
		list.remove(x.val)
		x = x.next
	}

}

func (list *List[T]) IsEmpty() bool {
	return list.count() == 0 && list.firstNode == nil
}

// SubList ...Creates a view of the list... starting at startIndex and ending at endIndex-1.
// In essence, the element at `endIndex` is not included
func (list *List[T]) SubList(startIndex int, endIndex int) (*List[T], error) {
	defer list.mu.Unlock()
	list.mu.Lock()

	if startIndex < 0 {
		return nil, errors.New("startIndex(" + strconv.Itoa(startIndex) + ") < 0 is not allowed")
	}
	sz := list.count()

	if endIndex > sz {
		panic("endIndex(" + strconv.Itoa(endIndex) + ") > listsize(" + strconv.Itoa(sz) + ") is not allowed")
	}

	if startIndex > endIndex {
		return nil, errors.New("startIndex(" + strconv.Itoa(startIndex) + ") > endIndex(" + strconv.Itoa(endIndex) + ") is not allowed")
	}

	subList := NewList[T]()
	start, end := list.getBoundaryNodes(startIndex, endIndex)
	subList.firstNode = start
	subList.lastNode = end
	subList.parent = list
	subList.parenLen = sz
	subList.size = endIndex - startIndex

	return subList, nil

}

func (list *List[T]) isSubList() bool {
	return list.parent != nil
}

/**
 * Returns the (non-nil) Node at the specified element index.
 */
func (list *List[T]) getNode(index int) (*lNode[T], error) {

	if index < 0 {
		return nil, errors.New("Index=(" + strconv.Itoa(index) + ") < 0 is not allowed")
	}

	sz := list.count()
	if index >= sz {
		return nil, errors.New("Index=(" + strconv.Itoa(index) + ") > list-size=(" + strconv.Itoa(list.count()) + ") is not allowed")
	}

	// NOTE x >> y is same as x ÷ 2^y
	if index < (sz >> 1) {
		x := list.firstNode
		for i := 0; i < index; i++ {
			x = x.next
		}
		return x, nil
	} else {
		x := list.lastNode
		for i := sz - 1; i > index; i-- {
			x = x.prev
		}
		return x, nil
	}

}

// getBoundaryNodes ... Return the nodes at the specified indexes
func (list *List[T]) getBoundaryNodes(start int, end int) (*lNode[T], *lNode[T]) {
	sz := list.count()
	if start >= 0 && start <= end && end <= sz {
		nd, _ := list.getNode(start)
		nd1, _ := list.getNode(end - 1)

		return nd, nd1
	}

	return nil, nil
}

// Get - returns the element at that index in the list
func (list *List[T]) Set(index int, val T) {
	defer list.mu.Unlock()
	list.mu.Lock()
	node, err := list.getNode(index)
	if err == nil {
		node.val = val
	}
}

// Get - returns the element at that index in the list
func (list *List[T]) Get(index int) (T, error) {
	defer list.mu.Unlock()
	list.mu.Lock()
	node, err := list.getNode(index)
	if err == nil {
		return node.val, nil
	}
	var nilVal T
	return nilVal, err
}

func (list *List[T]) getLastNode() *lNode[T] {
	return list.lastNode
}

func (list *List[T]) LastElement() interface{} {
	defer list.mu.Unlock()
	list.mu.Lock()
	return list.getLastNode().val
}

func (list *List[T]) Contains(val T) bool {
	return list.IndexOf(val) != -1

}

func (list *List[T]) IndexOf(val T) int {
	defer list.mu.Unlock()
	list.mu.Lock()

	x := list.firstNode

	if x.val == val {
		return 0
	}
	sz := list.count()
	for i := 0; i < sz; i++ {
		if val == x.val {
			return i
		}

		x = x.next
	}

	return -1

}
func (list *List[T]) indexOfNode(node *lNode[T]) int {

	x := list.firstNode

	if x == node {
		return 0
	}
	sz := list.count()
	for i := 0; i < sz; i++ {

		if node == x {
			return i
		}

		x = x.next
	}

	return -1

}

func (list *List[T]) containsNode(node *lNode[T]) bool {
	return list.indexOfNode(node) != -1
}

/**
 * Links val as first element.
 */
func (list *List[T]) prepend(val T) {
	f := list.firstNode
	newNode := initNode(nil, val, f)
	list.firstNode = newNode
	if f == nil {
		list.lastNode = newNode
	} else {
		f.prev = newNode
	}
	list.incrementSize(1)
}

/**
 * Links val as last element.
 */
func (list *List[T]) append(val T) {

	l := list.lastNode
	newNode := initNode(l, val, nil)
	list.lastNode = newNode
	if l == nil {
		list.firstNode = newNode
	} else {
		l.next = newNode
	}
	list.incrementSize(1)
}

/**
 * Inserts element e before non-null Node succ.
 * ONLY TO BE CALLED BY BASE PARENT LISTS
 * Return a pointer to the new node that was inserted.
 * This will help with spontaneous insertions
 */
func (list *List[T]) insertBefore(e T, succ *lNode[T]) *lNode[T] {

	prev := succ.prev

	newNode := initNode(prev, e, succ)

	succ.prev = newNode
	if prev == nil {
		list.firstNode = newNode
	} else {
		prev.next = newNode
	}
	list.incrementSize(1)
	return newNode
}

/**
 * Inserts element e after non-null Node succ.
 *
 * ONLY TO BE CALLED BY BASE PARENT LISTS
 * Return a pointer to the new node that was inserted.
 * This will help with spontaneous insertions
 */
func (list *List[T]) insertAfter(e T, succ *lNode[T]) *lNode[T] {

	next := succ.next

	newNode := initNode(succ, e, succ.next)

	succ.next = newNode
	if next == nil {
		list.lastNode = newNode
	} else {
		next.prev = newNode
	}
	list.incrementSize(1)

	return newNode
}

func (list *List[T]) Clear() bool {
	defer list.mu.Unlock()

	list.mu.Lock()
	list.clear()

	return true
}

func (list *List[T]) clear() {

	sz := list.count()
	first := list.firstNode
	last := list.lastNode

	if first != nil && last != nil {
		/**
		Must be a sublist embedded inside a list
		[23,9,10,12,34,28,99,55,32]--parent
		     [10,12,34,28]--sublist
		*/
		if first.prev != nil && last.next != nil {

			checkForNodeBeforeFirst := first.prev
			checkForNodeBeyondLast := last.next

			checkForNodeBeforeFirst.next = checkForNodeBeyondLast
			checkForNodeBeyondLast.prev = checkForNodeBeforeFirst

		}

		/**
		Must be a sublist starting from start of parent but parent extends beyond end of sublist e.g.
		[23,9,10,12,34,28,99]--parent
		[23,9,10,12]--sublist
		*/
		if first.prev == nil && last.next != nil {
			if list.parent != nil {
				list.parent.firstNode = last.next
			}
		}
		/**
		Must be a sublist ending on the tail of the parent, but starting somewhere within the parent e.g.
		[23,9,10,12,34,28,99,43,215,28]--parent
		                 [99,43,215,28]--sublist
		*/
		if first.prev != nil && last.next == nil {
			if list.parent != nil {
				list.parent.lastNode = first.prev
			}
		}
	}

	var nilVal T

	list.forEachNode(func(x *lNode[T]) bool {
		next := x.next
		x.val = nilVal
		x.next = nil
		x.prev = nil
		x = next
		return true
	})

	list.firstNode = nil
	list.lastNode = nil
	list.size = 0
	list.iter = nil
	list.nodeIter = nil

	if list.parent != nil {
		list.parent.decrementSize(sz)
	}

}

// Not tested yet
func (list *List[T]) removeLinkedRange(startNode *lNode[T], stopNode *lNode[T]) {

	defer list.mu.Unlock()
	list.mu.Lock()

	if startNode != nil && stopNode != nil {
		prev := startNode.prev
		next := stopNode.next

		if prev != nil {
			prev.next = next
		}
		if next != nil {
			next.prev = prev
		}

		x := startNode

		var nilVal T
		i := 0
		for {
			x.val = nilVal
			x = x.next
			i++
			if x == stopNode {
				x.val = nilVal
				i++
				break
			}
		}
		list.decrementSize(i)

	}

}

func (list *List[T]) Log(optionalLabel string) {
	defer list.mu.Unlock()

	list.mu.Lock()
	list.log(optionalLabel)
}

func (list *List[T]) log(optionalLabel string) {

	x := list.firstNode

	if x == nil {
		fmt.Println(optionalLabel + ":\n[], len: 0")
		return
	}

	counter := 0
	var bld strings.Builder

	bld.WriteString(optionalLabel)
	bld.WriteString(":\n[")

	sz := list.count()

	for ; x != nil; x = x.next {

		bld.WriteString(fmt.Sprintf("%v", x.val))
		counter++
		if x == list.lastNode {
			break
		}
		bld.WriteString(", ")
	}

	bld.WriteString("], len:")
	bld.WriteString(strconv.Itoa(sz))
	bld.WriteString(", confirm-len(")
	bld.WriteString(strconv.Itoa(counter))
	bld.WriteString(")")

	fmt.Println(bld.String())

}

// sync ... Sublists will use this method to synchronize their lengths with their parents.
// The core functionality here will run if the parent list size has changed since when the sublist last checked
func (list *List[T]) sync() {

	//Check for list beheading!...Head removed
	if list.firstNode == nil {
		return
	}

	//Check for list tail docking... the tail was removed
	if list.lastNode == nil {
		return
	}

	if list.firstNode.prev == nil && list.firstNode.isNilValOnNode() && list.firstNode.next == nil {
		list.firstNode = nil
		list.lastNode = nil
		list.iter = nil
		list.nodeIter = nil
		list.parent = nil
		list.parenLen = 0
		list.size = 0
	}

	//Run core sync method functionality only if the list has a parent
	if list.parent != nil {
		parenLen := list.parent.count()
		sizeChanged := parenLen != list.parenLen

		if sizeChanged {

			i := 0
			list.forEachNode(func(x *lNode[T]) bool {
				i++
				return true
			})
			list.size = i
			list.parenLen = parenLen
		}
	}
}

func (list *List[T]) count() int {
	list.sync()
	return list.size
}

func (list *List[T]) Count() int {
	defer list.mu.Unlock()
	list.mu.Lock()
	return list.size
}
