package main

import "fmt"

type Handler func(data any)

type handlerEntry struct {
	id      int
	handler Handler
}

type PubSub struct {
	events        map[string][]handlerEntry
	nextHandlerID int
}

func NewPubSub() *PubSub {
	return &PubSub{events: make(map[string][]handlerEntry)}
}

type Subscription struct {
	unsubscribe func()
}

func (s Subscription) Unsubscribe() {
	s.unsubscribe()
}

func (p *PubSub) Subscribe(event string, cb Handler) Subscription {
	p.nextHandlerID++
	id := p.nextHandlerID
	p.events[event] = append(p.events[event], handlerEntry{id: id, handler: cb})
	return Subscription{
		unsubscribe: func() { p.unsubscribeByID(event, id) },
	}
}

func (p *PubSub) unsubscribeByID(event string, id int) {
	handlers := p.events[event]
	filtered := handlers[:0]
	for _, h := range handlers {
		if h.id != id {
			filtered = append(filtered, h)
		}
	}
	p.events[event] = filtered
}

func (p *PubSub) Publish(event string, data any) {
	for _, entry := range p.events[event] {
		entry.handler(data)
	}
}

func demoMapPubSub() {
	pubsub := NewPubSub()

	consumer1 := pubsub.Subscribe("movie", func(data any) {
		fmt.Println("from 1st movie", data)
	})
	pubsub.Subscribe("movie", func(data any) {
		fmt.Println("from 2nd movie", data)
	})
	pubsub.Subscribe("music", func(data any) {
		fmt.Println("from 1st music", data)
	})

	pubsub.Publish("movie", "Fight Club")
	consumer1.Unsubscribe()
	pubsub.Publish("movie", "Fight Club 2")
}
