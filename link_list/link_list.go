package link_list

import (
	"fmt"
	"sync"
	"unsafe"
)

type DLL struct {
	root *Node
	last *Node
	sync.RWMutex
}

type Node struct {
	left  *Node
	right *Node
	value Value
}

type Value struct {
	pointer unsafe.Pointer //we need 2 thing memory page and location in memory, we can get that data from pointer
	key     string         //link my hash table
}

func NewValue(pointer unsafe.Pointer, key string) Value {
	return Value{
		pointer: pointer,
		key:     key,
	}
}

// is impassible to make data race here?
func (dll *DLL) GetLRUFreeSpace(lru *Node, blockSize int) []byte {
	dll.Lock()
	defer dll.Unlock()

	ptr := lru.value.pointer
	return unsafe.Slice((*byte)(ptr), blockSize)
}

func (dll *DLL) Inset(value Value) *Node {
	dll.Lock()
	defer dll.Unlock()

	newNode, root := &Node{}, dll.root

	if dll.root != nil {
		newNode.right = root
		root.left = newNode
	} else { //if list empty then last element is first element
		dll.last = newNode
	}

	newNode.value = value
	dll.root = newNode

	return newNode
}

func (dll *DLL) Delete(node *Node) {
	if node == nil {
		return
	}

	dll.Lock()
	defer dll.Unlock()

	if node.left != nil {
		node.left.right = node.right
	} else {
		dll.root = node.right
	}

	if node.right != nil {
		node.right.left = node.left
	} else {
		dll.last = node.left
	}

	node.left = nil
	node.right = nil
}

func (dll *DLL) Remove() {
	dll.Lock()
	defer dll.Unlock()

	if dll.last == nil {
		return
	}

	if dll.last.left == nil {
		dll.last = nil
		dll.root = nil
		return
	}

	dll.last = dll.last.left
	dll.last.right = nil
}

func (dll *DLL) LastNode() *Node {
	return dll.last
}

func (dll *DLL) Read(node *Node) {
	dll.Lock()
	defer dll.Unlock()

	if node == dll.root {
		return
	}

	node.left.right = node.right
	node.right.left = node.left

	node.right = dll.root
	node.left = nil

	dll.root.left = node

	dll.root = node
}

func (dll *DLL) ReadAll() {
	for root := dll.root; root != nil; root = root.right {
		fmt.Println(root.value)
	}
}

func (dll *DLL) ReadBack() {
	for current := dll.last; current != nil; current = current.left {
		fmt.Println(current.value)
	}
}
