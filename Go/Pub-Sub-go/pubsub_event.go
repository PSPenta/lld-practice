package main

import (
	"fmt"
	"sync"
	"time"
)

type EventEmitter struct {
	mu           sync.Mutex
	listeners    map[string][]func(any)
	maxListeners int
}

func NewEventEmitter() *EventEmitter {
	return &EventEmitter{
		listeners:    make(map[string][]func(any)),
		maxListeners: 5,
	}
}

func (e *EventEmitter) SetMaxListeners(n int) {
	e.maxListeners = n
}

func (e *EventEmitter) On(topic string, callback func(any)) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.listeners[topic] = append(e.listeners[topic], callback)
}

func (e *EventEmitter) Emit(topic string, data any) {
	e.mu.Lock()
	handlers := append([]func(any){}, e.listeners[topic]...)
	e.mu.Unlock()
	for _, cb := range handlers {
		cb(data)
	}
}

var emitterPubSub = NewEventEmitter()

func ProduceEvent(topic string, data any) {
	emitterPubSub.Emit(topic, data)
}

type EventConsumer struct {
	topic string
}

func ConsumeEvent(topic string) *EventConsumer {
	return &EventConsumer{topic: topic}
}

func (c *EventConsumer) On(event string, callback func(any)) {
	if event == "data" {
		emitterPubSub.On(c.topic, callback)
	}
}

func demoEventEmitterPubSub() {
	emitterPubSub.SetMaxListeners(5)
	data := 0
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for range ticker.C {
			ProduceEvent("myTopic", data)
			data++
			if data > 10 {
				close(done)
				return
			}
		}
	}()

	consumer := ConsumeEvent("myTopic")
	consumer.On("data", func(d any) {
		fmt.Println("got data", d)
	})
	<-done
}
