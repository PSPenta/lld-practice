package main

import "fmt"

type FIFOQueue struct {
	capacity int
	front    int
	rear     int
	queue    []any
}

func NewFIFOQueue(capacity int) *FIFOQueue {
	return &FIFOQueue{
		capacity: capacity,
		front:    -1,
		rear:     -1,
		queue:    make([]any, capacity),
	}
}

func (q *FIFOQueue) Enqueue(value any) any {
	rear := (q.rear + 1) % q.capacity
	if q.queue[rear] != nil {
		fmt.Println("queue is full")
		return -1
	}
	q.queue[rear] = value
	q.rear = rear
	return value
}

func (q *FIFOQueue) Dequeue() any {
	front := (q.front + 1) % q.capacity
	if q.queue[front] == nil {
		fmt.Println("queue is empty")
		return -1
	}
	value := q.queue[front]
	q.queue[front] = nil
	q.front = front
	return value
}

func main() {
	fifo := NewFIFOQueue(5)
	fifo.Enqueue(1)
	fifo.Enqueue(2)
	fifo.Enqueue(3)
	fifo.Enqueue(4)
	fifo.Enqueue(5)
	fifo.Enqueue(6)
}
