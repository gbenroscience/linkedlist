package ds

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// abstractAnyList - An abstraction of a list
type abstractAnyList[T any] interface {
	Add(val T)
	AddVal(val T, index int) bool
	AddAll(lst *AnyList[T]) bool
	AddAllAt(index int, lst *AnyList[T]) bool

	Remove(val T) bool
	RemoveIndex(index int) bool
	RemoveAll(lst *AnyList[T]) *AnyList[T]
	Clear() bool

	Set(index int, val T)
	Get(index int) T
	ToArray() []T
	LastElement() T

	SubList(startIndex int, endIndex int) (*AnyList[T], error)

	IsEmpty() bool
	Contains(val T) bool

	IndexOf(val T) int

	Log(optionalLabel string)
}

// node - A list node
type node[T any] struct {
	next *node[T]
	prev *node[T]
	val  T
}

// AnyList - The AnyList
type AnyList[T any] struct {
	size      int
	firstNode *node[T]
	lastNode  *node[T]
	parent    *AnyList[T]
	parenLen  int
	mu        sync.Mutex
	//Used for rapid iteration over the list's values
	iter *node[T]
	//Used for rapid iteration over the list's nodes
	nodeIter *node[T]
	// Every instance had better override this function after calling the NewAnyList function in order to gain speed in the Remove, IndexOf and other relevant function
	Equals func(val1 T, val2 T) bool
}

func NewAnyList[T any]() *AnyList[T] {
	list := new(AnyList[T])

	list.size = 0
	list.parent = nil
	list.firstNode = nil
	list.lastNode = nil
	list.iter = nil
	list.nodeIter = nil
	list.mu = sync.Mutex{}

	list.Equals = func(val1 T, val2 T) bool {
		return fmt.Sprintf("%v", val1) == fmt.Sprintf("%v", val2)
	}

	return list
}

func init_node[T any](prev *node[T], val T, next *node[T]) *node[T] {
	node := new(node[T])
	node.prev = prev
	node.next = next
	node.val = val

	return node
}

func (node *node[T]) isNilValOnNode() bool {
	return new(T) == &node.val
}

func (list *AnyList[T]) nextNode() *node[T] {

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
func (list *AnyList[T]) resetNodeIterator() {
	if list.nodeIter != nil {
		list.nodeIter = nil
	}
}

func (list *AnyList[T]) next() *T {

	var next = list.nextNode()
	if next != nil {
		return &next.val
	}

	return nil
}

// Call this to reset the values iterator
func (list *AnyList[T]) resetIterator() {
	if list.iter != nil {
		list.iter = nil
	}
}

func (list *AnyList[T]) ForEach(function func(val T) bool) {

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

func (list *AnyList[T]) forEachNodeFrom(start *node[T], function func(node *node[T]) bool) {

	var x *node[T]
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

func (list *AnyList[T]) forEachNode(function func(node *node[T]) bool) {

	var x *node[T]
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
func (list *AnyList[T]) ToArray() []T {

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
func (list *AnyList[T]) addNode(elem *node[T]) {
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
func (list *AnyList[T]) Add(val T) {
	defer list.mu.Unlock()

	list.mu.Lock()
	list.add(val)

}

// TESTED
func (list *AnyList[T]) add(val T) {
	list.append(val)
}

func (list *AnyList[T]) AddVal(val T, index int) bool {
	defer list.mu.Unlock()

	list.mu.Lock()
	_, _ = list.addVal(val, index)

	return true
}

// TESTED
func (list *AnyList[T]) addVal(val T, index int) (bool, error) {
	node := init_node(nil, val, nil)
	return list.addNodeAt(node, index)

}

func (list *AnyList[T]) AddValues(args ...T) {
	defer list.mu.Unlock()

	list.mu.Lock()
	list.addValues(args...)

}

// TESTED
func (list *AnyList[T]) addValues(args ...T) {
	for _, v := range args {
		list.append(v)
	}
}

func (list *AnyList[T]) AddArray(array []T) {
	defer list.mu.Unlock()

	list.mu.Lock()
	list.addArray(array)
}

// TESTED
func (list *AnyList[T]) addArray(array []T) {
	list.addValues(array...)
}

func (list *AnyList[T]) AddAll(lst *AnyList[T]) bool {
	defer list.mu.Unlock()

	list.mu.Lock()
	_ = list.addAll(lst)

	return true
}

// TESTED
func (list *AnyList[T]) addAll(lst *AnyList[T]) error {
	err := list.addAllAt(list.count(), lst)
	return err
}

func (list *AnyList[T]) AddAllAt(index int, lst *AnyList[T]) bool {
	defer list.mu.Unlock()

	list.mu.Lock()
	_ = list.addAllAt(index, lst)

	return true
}

func (list *AnyList[T]) Clone() *AnyList[T] {

	ls := new(AnyList[T])

	list.forEachNode(func(node *node[T]) bool {
		ls.Add(node.val)
		return true
	})

	return ls
}

func (list *AnyList[T]) addAllAt(index int, lst *AnyList[T]) error {

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
func (list *AnyList[T]) addNodeAt(elem *node[T], index int) (bool, error) {

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

func (list *AnyList[T]) incrementSize(dx int) {
	list.size += dx
	if list.parent != nil {
		list.parent.size += dx
	}
}

func (list *AnyList[T]) decrementSize(dx int) {
	list.size -= dx
	if list.parent != nil {
		list.parent.size -= dx
	}
}

func (list *AnyList[T]) removeNode(elem *node[T]) bool {

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
func (list *AnyList[T]) removeIndex(index int) bool {

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

func (list *AnyList[T]) Remove(val T) bool {

	defer list.mu.Unlock()
	list.mu.Lock()
	list.remove(val)

	return true
}

/**
 * Remove the first node that has
 * the same value as the parameter
 */
func (list *AnyList[T]) remove(val T) bool {

	//empty list
	if list.firstNode == nil {
		return false
	}

	x := list.firstNode
	sz := list.count()
	for i := 0; i < sz; i++ {
		if list.Equals(x.val, val) { // if x.val == val{
			succ := list.removeNode(x)

			return succ
		}

		x = x.next
	}

	return false

}

func (list *AnyList[T]) RemoveIndex(index int) bool {
	defer list.mu.Unlock()
	list.mu.Lock()
	return list.removeIndex(index)
}

func (list *AnyList[T]) RemoveAll(lst *AnyList[T]) bool {
	defer list.mu.Unlock()
	list.mu.Lock()
	list.removeAll(lst)

	return true
}

// TESTED
func (list *AnyList[T]) removeAll(lst *AnyList[T]) {

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

func (list *AnyList[T]) IsEmpty() bool {
	return list.count() == 0 && list.firstNode == nil
}

// SubList ...Creates a view of the list... starting at startIndex and ending at endIndex-1.
// In essence, the element at `endIndex` is not included
func (list *AnyList[T]) SubList(startIndex int, endIndex int) (*AnyList[T], error) {
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

	subList := NewAnyList[T]()
	subList.Equals = list.Equals
	start, end := list.getBoundaryNodes(startIndex, endIndex)
	subList.firstNode = start
	subList.lastNode = end
	subList.parent = list
	subList.parenLen = sz
	subList.size = endIndex - startIndex

	return subList, nil

}

func (list *AnyList[T]) isSubList() bool {
	return list.parent != nil
}

/**
 * Returns the (non-nil) Node at the specified element index.
 */
func (list *AnyList[T]) getNode(index int) (*node[T], error) {

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
func (list *AnyList[T]) getBoundaryNodes(start int, end int) (*node[T], *node[T]) {
	sz := list.count()
	if start >= 0 && start <= end && end <= sz {
		nd, _ := list.getNode(start)
		nd1, _ := list.getNode(end - 1)

		return nd, nd1
	}

	return nil, nil
}

// Get - returns the element at that index in the list
func (list *AnyList[T]) Set(index int, val T) {
	defer list.mu.Unlock()
	list.mu.Lock()
	node, err := list.getNode(index)
	if err == nil {
		node.val = val
	}
}

// Get - returns the element at that index in the list
func (list *AnyList[T]) Get(index int) (T, error) {
	defer list.mu.Unlock()
	list.mu.Lock()
	node, err := list.getNode(index)
	if err == nil {
		return node.val, nil
	}
	var nilVal T
	return nilVal, err
}

func (list *AnyList[T]) getLastNode() *node[T] {
	return list.lastNode
}

func (list *AnyList[T]) LastElement() interface{} {
	defer list.mu.Unlock()
	list.mu.Lock()
	return list.getLastNode().val
}

func (list *AnyList[T]) Contains(val T) bool {
	return list.IndexOf(val) != -1

}

func (list *AnyList[T]) IndexOf(val T) int {
	defer list.mu.Unlock()
	list.mu.Lock()

	x := list.firstNode

	if list.Equals(x.val, val) { //if x.val == val
		return 0
	}
	sz := list.count()
	for i := 0; i < sz; i++ {
		if list.Equals(val, x.val) { //if val == x.val
			return i
		}

		x = x.next
	}

	return -1

}
func (list *AnyList[T]) indexOfNode(node *node[T]) int {

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

func (list *AnyList[T]) containsNode(node *node[T]) bool {
	return list.indexOfNode(node) != -1
}

/**
 * Links val as first element.
 */
func (list *AnyList[T]) prepend(val T) {
	f := list.firstNode
	newNode := init_node(nil, val, f)
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
func (list *AnyList[T]) append(val T) {

	l := list.lastNode
	newNode := init_node(l, val, nil)
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
func (list *AnyList[T]) insertBefore(e T, succ *node[T]) *node[T] {

	prev := succ.prev

	newNode := init_node(prev, e, succ)

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
func (list *AnyList[T]) insertAfter(e T, succ *node[T]) *node[T] {

	next := succ.next

	newNode := init_node(succ, e, succ.next)

	succ.next = newNode
	if next == nil {
		list.lastNode = newNode
	} else {
		next.prev = newNode
	}
	list.incrementSize(1)

	return newNode
}

func (list *AnyList[T]) Clear() bool {
	defer list.mu.Unlock()

	list.mu.Lock()
	list.clear()

	return true
}

func (list *AnyList[T]) clear() {

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

	list.forEachNode(func(x *node[T]) bool {
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
func (list *AnyList[T]) removeLinkedRange(startNode *node[T], stopNode *node[T]) {

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

func (list *AnyList[T]) Log(optionalLabel string) {
	defer list.mu.Unlock()

	list.mu.Lock()
	list.log(optionalLabel)
}

func (list *AnyList[T]) log(optionalLabel string) {

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
func (list *AnyList[T]) sync() {

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
			list.forEachNode(func(x *node[T]) bool {
				i++
				return true
			})
			list.size = i
			list.parenLen = parenLen
		}
	}
}

func (list *AnyList[T]) count() int {
	list.sync()
	return list.size
}

func (list *AnyList[T]) Count() int {
	defer list.mu.Unlock()
	list.mu.Lock()
	return list.size
}
