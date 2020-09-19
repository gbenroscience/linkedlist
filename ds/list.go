package ds



import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

const (
	Shutdown = 1
)



// AbstractList - An abstraction of a list
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
	parent    *List
	parenLen  int
	mu        sync.Mutex
	//Used for rapid iteration over the list's values
	iter *Node
	//Used for rapid iteration over the list's nodes
	nodeIter *Node
}

func NewList() *List {
	list := new(List)

	list.size = 0
	list.parent = nil
	list.firstNode = nil
	list.lastNode = nil
	list.iter = nil
	list.nodeIter = nil
	list.mu = sync.Mutex{}

	return list
}


func initNode(prev *Node, val interface{}, next *Node) *Node {
	node := new(Node)
	node.prev = prev
	node.next = next
	node.val = val

	return node
}

func (list *List) nextNode() *Node {

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

//Call this to reset the nodes iterator
func (list *List) resetNodeIterator() {
	if list.nodeIter != nil {
		list.nodeIter = nil
	}
}

func (list *List) next() interface{} {

	if list.iter == nil {
		if list.firstNode == nil {
			return nil
		}
		list.iter = list.firstNode
		return list.iter.val
	} else {
		if list.iter == list.lastNode {
			return nil
		}
		if list.iter.next != nil {
			list.iter = list.iter.next
			return list.iter.val
		}
		return nil
	}
}

//Call this to reset the values iterator
func (list *List) resetIterator() {
	if list.iter != nil {
		list.iter = nil
	}
}

func (list *List) ForEach(function func(val interface{}) bool) {

	defer list.mu.Unlock()
	list.mu.Lock()
	var x interface{}
	list.resetIterator()

	for ; ; {
		x = list.next()
		if x == nil {
			break
		}
		if !function(x) {
			break
		}

	}
}

func (list *List) forEachNodeFrom(start *Node, function func(node *Node) bool) {

	var x *Node
	list.resetNodeIterator()
	list.nodeIter = start

	for ; ; {
		x = list.nextNode()
		if x == nil {
			break
		}
		if !function(x) {
			break
		}

	}
}

func (list *List) forEachNode(function func(node *Node) bool) {

	var x *Node
	list.resetNodeIterator()

	for ; ; {
		x = list.nextNode()
		if x == nil {
			break
		}
		if !function(x) {
			break
		}

	}
}

//TESTED
func (list *List) ToArray() []interface{} {

	result := make([]interface{}, list.count())

	i := 0
	list.ForEach(func(x interface{}) bool {
		result[i] = x
		i++
		return true
	})

	return result
}

//TESTED
func (list *List) addNode(elem *Node) {
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

//TESTED
func (list *List) Add(val interface{}) {
	defer list.mu.Unlock()

	list.mu.Lock()
	list.add(val)

}

//TESTED
func (list *List) add(val interface{}) {
	list.append(val)
}

func (list *List) AddVal(val interface{}, index int) bool {
	defer list.mu.Unlock()

	list.mu.Lock()
	_, _ = list.addVal(val, index)

	return true
}

//TESTED
func (list *List) addVal(val interface{}, index int) (bool, error) {
	node := initNode(nil, val, nil)
	return list.addNodeAt(node, index)

}

func (list *List) AddValues(args ...interface{}) {
	defer list.mu.Unlock()

	list.mu.Lock()
	list.addValues(args...)

}

//TESTED
func (list *List) addValues(args ...interface{}) {
	for _, v := range args {
		list.append(v)
	}
}

func (list *List) AddArray(array []interface{}) {
	defer list.mu.Unlock()

	list.mu.Lock()
	list.addArray(array)
}

//TESTED
func (list *List) addArray(array []interface{}) {
	list.addValues(array...)
}

func (list *List) AddAll(lst *List) bool {
	defer list.mu.Unlock()

	list.mu.Lock()
	_ = list.addAll(lst)


	return true
}

//TESTED
func (list *List) addAll(lst *List) error {
	err := list.addAllAt(list.count(), lst)
	return err
}

func (list *List) AddAllAt(index int, lst *List) bool {
	defer list.mu.Unlock()

	list.mu.Lock()
	_ = list.addAllAt(index, lst)

	return true
}

func (list *List) Clone() *List {

	ls := new(List)

	x := list.firstNode

	for ; x != nil; {
		ls.append(x.val)
		x = x.next
	}

	return ls
}

func (list *List) addAllAt(index int, lst *List) error {

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
		nodeAtIndex := list.getNode(index)

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
func (list *List) addNodeAt(elem *Node, index int) (bool, error) {

	sz := list.count()
	if index >= 0 && index <= sz {

		if index == sz {
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
			list.incrementSize(1)
		}

		return true, nil

	} else {
		return false, errors.New("index must lie between 0 and " + strconv.Itoa(list.count()))
	}

}

func (list *List) incrementSize(dx int) {
	list.size += dx
	if list.parent != nil {
		list.parent.size += dx
	}
}

func (list *List) decrementSize(dx int) {
	list.size -= dx
	if list.parent != nil {
		list.parent.size -= dx
	}
}

func (list *List) removeNode(elem *Node) bool {

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
	list.decrementSize(1)
	return true

}
func (list *List) removeIndex(index int) bool {

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

func (list *List) Remove(val interface{}) bool {

	defer list.mu.Unlock()
	list.mu.Lock()
	list.remove(val)

	return true
}

/**
 * Remove the first node that has
 * the same value as the parameter
 */
func (list *List) remove(val interface{}) bool {

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

func (list *List) RemoveIndex(index int) bool {
	defer list.mu.Unlock()
	list.mu.Lock()
	return list.removeIndex(index)
}

func (list *List) RemoveAll(lst *List) bool {
	defer list.mu.Unlock()
	list.mu.Lock()
	list.removeAll(lst)

	return true
}

//TESTED
func (list *List) removeAll(lst *List) {

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

func (list *List) IsEmpty() bool {
	return list.count() == 0 && list.firstNode == nil
}

// SubList ...Creates a view of the list... starting at startIndex and ending at endIndex-1.
// In essence, the element at `endIndex` is not included
func (list *List) SubList(startIndex int, endIndex int) (*List, error) {
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

	subList := NewList()
	start, end := list.getBoundaryNodes(startIndex, endIndex)
	subList.firstNode = start
	subList.lastNode = end
	subList.parent = list
	subList.parenLen = sz
	subList.size = endIndex - startIndex

	return subList, nil

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

	sz := list.count()
	if index >= sz {
		panic("Index=(" + strconv.Itoa(index) + ") > list-size=(" + strconv.Itoa(list.count()) + ") is not allowed")
	}

	// NOTE x >> y is same as x ÷ 2^y
	if index < (sz >> 1) {
		x := list.firstNode
		for i := 0; i < index; i++ {
			x = x.next
		}
		return x
	} else {
		x := list.lastNode
		for i := sz - 1; i > index; i-- {
			x = x.prev
		}
		return x
	}

}

// getBoundaryNodes ... Return the nodes at the specified indexes
func (list *List) getBoundaryNodes(start int, end int) (*Node, *Node) {
	sz := list.count()
	if start >= 0 && start <= end && end <= sz {
		return list.getNode(start), list.getNode(end - 1)
	}

	return nil, nil
}

//Get - returns the element at that index in the list
func (list *List) Get(index int) interface{} {
	defer list.mu.Unlock()
	list.mu.Lock()
	return list.getNode(index).val
}

func (list *List) getLastNode() *Node {
	return list.lastNode
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
func (list *List) indexOfNode(node *Node) int {

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

func (list *List) containsNode(node *Node) bool {
	return list.indexOfNode(node) != -1
}

/**
 * Links val as first element.
 */
func (list *List) prepend(val interface{}) {
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
func (list *List) append(val interface{}) {

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
func (list *List) insertBefore(e interface{}, succ *Node) *Node {

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
func (list *List) insertAfter(e interface{}, succ *Node) *Node {

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

func (list *List) Clear() bool {
	defer list.mu.Unlock()

	list.mu.Lock()
	list.clear()

	return true
}

func (list *List) clear() {

	sz := list.count()
	first := list.firstNode
	last := list.lastNode

	if first != nil && last != nil {

		defer func() {
			if list.parent != nil{
				//list.close()
			}
		}()

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

	list.forEachNode(func(x *Node) bool {
		next := x.next
		x.val = nil
		x.next = nil
		x.prev = nil
		x = next
		return true
	})
	if list.parent != nil {
		list.parent.decrementSize(sz)
	}

	list.firstNode = nil
	list.lastNode = nil
	list.iter = nil
	list.size = 0

}

//Not tested yet
func (list *List) removeLinkedRange(startNode *Node, stopNode *Node) {

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

		i := 0
		for {
			x.val = nil
			x = x.next
			i++
			if x == stopNode {
				x.val = nil
				i++
				break
			}
		}
		list.decrementSize(i)

	}

}

func (list *List) Log(optionalLabel string) bool {
	defer list.mu.Unlock()

	list.mu.Lock()
	list.log(optionalLabel)


	return true
}

func (list *List) log(optionalLabel string) {

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
		bld.WriteString(", ")
		counter++
		if x == list.lastNode {
			break
		}
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
func (list *List) sync() {

	//Check for list beheading!...Head removed
	if list.firstNode == nil {
		panic("Oops. This list was beheaded prior to this action! List beheading is not supported for sublists!")
		return
	}

	//Check for list tail docking... the tail was removed
	if list.lastNode == nil {
		panic("Oops. This list was tail-docked(the tail was removed) prior to this action! Tail docking is not supported for sublists!")
		return
	}

	//Run core sync method functionality only if the list has a parent
	if list.parent != nil {
		parenLen := list.parent.count()
		sizeChanged := parenLen != list.parenLen

		if sizeChanged {

			i := 0
			list.forEachNode(func(x *Node) bool {
				i++
				return true
			})
			list.size = i
			list.parenLen = parenLen
		}
	}
}

func (list *List) count() int {
	list.sync()
	return list.size
}


func (list *List) Count() int {
	defer list.mu.Unlock()
	list.mu.Lock()
	return list.size
}