package trie

import "sync"

type Node[Data any] struct {
	mu       sync.RWMutex
	children map[string]*Node[Data]

	Data Data
}

func NewNode[Data any]() *Node[Data] {
	return &Node[Data]{
		children: make(map[string]*Node[Data]),
	}
}

func (n *Node[Data]) Insert(key string) *Node[Data] {
	n.children[key] = NewNode[Data]()
	return n.children[key]
}

func (n *Node[Data]) Get(key string) *Node[Data] {
	n.mu.RLock()
	defer n.mu.RUnlock()

	child, ok := n.children[key]
	if !ok {
		return nil
	}
	return child
}

func (n *Node[Data]) GetAll() map[string]*Node[Data] {
	return n.children
}

//func (n *Node[Data]) SetData(key string, data Data) bool {
//	child, ok := n.children[key]
//	if !ok {
//		return false
//	}
//	child.data = data
//	return true
//}
//
//func (n *Node[Data]) GetData(key string) (Data, bool) {
//	n.mu.RLock()
//	defer n.mu.RUnlock()
//
//	child, ok := n.children[key]
//	if !ok {
//		var empty Data
//		return empty, false
//	}
//	return child.data, true
//}
