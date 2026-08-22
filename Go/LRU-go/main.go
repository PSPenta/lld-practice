package main

import "fmt"

type Node struct {
	key, value any
	next, prev  *Node
}

type LRUCache struct {
	capacity int
	cache    map[any]*Node
	head     *Node
	tail     *Node
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[any]*Node),
	}
}

func (l *LRUCache) Get(key any) any {
	node, ok := l.cache[key]
	if !ok {
		return -1
	}
	l.removeNode(node)
	l.addNode(node)
	return node.value
}

func (l *LRUCache) Set(key, value any) int {
	if existing, ok := l.cache[key]; ok {
		existing.value = value
		if existing != l.head {
			l.removeNode(existing)
			l.addNode(existing)
		}
	} else {
		node := &Node{key: key, value: value}
		l.cache[key] = node
		l.addNode(node)
		if len(l.cache) > l.capacity {
			tail := l.tail
			delete(l.cache, tail.key)
			l.removeNode(tail)
		}
	}
	return 1
}

func (l *LRUCache) removeNode(node *Node) {
	if node == l.head {
		l.head = node.next
	}
	if node == l.tail {
		l.tail = node.prev
	}
	if node.prev != nil {
		node.prev.next = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	}
}

func (l *LRUCache) addNode(node *Node) {
	if l.head == nil {
		l.head = node
		l.tail = node
		return
	}
	node.next = l.head
	l.head.prev = node
	l.head = node
}

func main() {
	lru := NewLRUCache(3)
	lru.Set(1, "A")
	lru.Set(2, "B")
	lru.Set(3, "C")
	lru.Set(4, "D")
	fmt.Println(lru.Get(1))
	fmt.Println(lru.Get(4))
}
