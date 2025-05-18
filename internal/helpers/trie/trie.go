package trie

import "fmt"

type Node[Data any] struct {
	Data     Data
	children map[string]*Node[Data]
}

func NewNode[Data any]() *Node[Data] {
	return &Node[Data]{
		children: make(map[string]*Node[Data]),
	}
}

func (n *Node[Data]) Replace(other *Node[Data]) {
	n.Data = other.Data
	n.children = other.children
}

func (n *Node[Data]) Insert(key string) *Node[Data] {
	n.children[key] = NewNode[Data]()
	return n.children[key]
}

func (n *Node[Data]) Get(key string) *Node[Data] {
	child, ok := n.children[key]
	if !ok {
		return nil
	}
	return child
}

func (n *Node[Data]) GetAll() map[string]*Node[Data] {
	return n.children
}

func (n *Node[Data]) Print() {
	if n == nil {
		return
	}

	fmt.Println(n.Data)
	for name, child := range n.children {
		fmt.Println(name)
		child.Print()
	}
}
