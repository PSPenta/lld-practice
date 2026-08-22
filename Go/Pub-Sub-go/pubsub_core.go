package main

import (
	"fmt"
	"sync"
	"time"
)

type topicStore struct {
	subscribers []func(any)
}

var (
	coreMu     sync.RWMutex
	corePubSub = make(map[string]*topicStore)
)

func Produce(topic string, data any) {
	coreMu.RLock()
	store := corePubSub[topic]
	coreMu.RUnlock()
	if store == nil {
		return
	}
	for _, sub := range store.subscribers {
		sub(data)
	}
}

type Consumer struct {
	topic string
}

func Consume(topic string) *Consumer {
	coreMu.Lock()
	if corePubSub[topic] == nil {
		corePubSub[topic] = &topicStore{}
	}
	coreMu.Unlock()
	return &Consumer{topic: topic}
}

func (c *Consumer) On(event string, callback func(any)) {
	if event == "data" {
		coreMu.Lock()
		corePubSub[c.topic].subscribers = append(corePubSub[c.topic].subscribers, callback)
		coreMu.Unlock()
	}
}

func demoCorePubSub() {
	data := 0
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for range ticker.C {
			Produce("myTopic", data)
			data++
			if data > 10 {
				close(done)
				return
			}
		}
	}()

	consumer := Consume("myTopic")
	consumer.On("data", func(d any) {
		fmt.Println("got data", d)
	})
	<-done
}
