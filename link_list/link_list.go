package link_list

import (
	"fmt"
	"sync"
)

type DLL struct {
	root *Node
	last *Node
	sync.RWMutex
}

type Node struct {
	left  *Node
	right *Node
	value int
}

func (dll *DLL) Inset(value int) *Node {
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
