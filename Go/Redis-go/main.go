package main

import (
	"fmt"
	"sync"
	"time"
)

type CacheNode struct {
	key, value any
	ttl        int64
	next, prev *CacheNode
}

type Redis struct {
	mu         sync.Mutex
	allNodes   map[any]*CacheNode
	size       int
	capacity   int
	head       *CacheNode
	tail       *CacheNode
	defaultTTL int64
}

func NewRedis(capacity int) *Redis {
	return &Redis{
		allNodes:   make(map[any]*CacheNode),
		capacity:   capacity,
		defaultTTL: 3600,
	}
}

func (r *Redis) Set(key, value any, expiry int64) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if expiry == 0 {
		expiry = r.defaultTTL
	}
	node := &CacheNode{
		key:   key,
		value: value,
		ttl:   time.Now().UnixMilli() + expiry*1000,
	}
	r.removeNode(r.allNodes[key])
	r.addNode(node)
	return true
}

func (r *Redis) Get(key any) any {
	r.mu.Lock()
	defer r.mu.Unlock()
	node, ok := r.allNodes[key]
	if !ok {
		return nil
	}
	if node.ttl > time.Now().UnixMilli() {
		r.removeNode(node)
		r.addNode(node)
		return node.value
	}
	r.removeNode(node)
	return nil
}

func (r *Redis) addNode(node *CacheNode) {
	if r.head == nil {
		r.head = node
		r.tail = node
		r.size++
		r.allNodes[node.key] = node
		return
	}
	node.next = r.head
	r.head.prev = node
	r.head = node
	r.size++
	r.allNodes[node.key] = node
	if r.size > r.capacity {
		r.evict()
	}
}

func (r *Redis) removeNode(node *CacheNode) {
	if node == nil {
		return
	}
	prev, next := node.prev, node.next
	if prev != nil {
		prev.next = next
	} else {
		r.head = next
	}
	if next != nil {
		next.prev = prev
	} else {
		r.tail = prev
	}
	delete(r.allNodes, node.key)
	r.size--
}

func (r *Redis) evict() {
	if expired := r.getExpiredNode(); expired != nil {
		r.removeNode(expired)
		return
	}
	r.removeNode(r.tail)
}

func (r *Redis) getExpiredNode() *CacheNode {
	now := time.Now().UnixMilli()
	node := r.tail
	for node != nil && node.ttl > now {
		node = node.prev
	}
	return node
}

func main() {
	redis := NewRedis(10)
	for i := 1; i <= 16; i++ {
		key := fmt.Sprintf("test %d", i)
		redis.Set(key, i, 0)
		fmt.Println(redis.Get(key))
	}
	fmt.Println(redis.Get("test 16"))
	fmt.Println(redis.Get("test 1"))
}
