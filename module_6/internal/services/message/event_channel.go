package message

import (
	"context"
	"module6/internal/logger"
	"sync"
	"sync/atomic"
)

type EventChannel struct {
	subscribes map[chan<- *Message]struct{}
	mx         sync.RWMutex
	msgs       chan *Message
	isRunning  atomic.Bool
}

func StartEventChannel() *EventChannel {
	result := &EventChannel{
		subscribes: make(map[chan<- *Message]struct{}),
		msgs:       make(chan *Message, 100),
	}
	result.isRunning.Store(true)
	go result.run()
	return result
}

func (e *EventChannel) PublishMessage(message *Message) {
	e.msgs <- message
}

func (e *EventChannel) Subscribe(c chan<- *Message) {
	e.mx.Lock()
	defer e.mx.Unlock()
	e.subscribes[c] = struct{}{}
}

func (e *EventChannel) Unsubscribe(c chan<- *Message) {
	e.mx.Lock()
	defer e.mx.Unlock()
	delete(e.subscribes, c)
}

func (e *EventChannel) IsEmpty() bool {
	e.mx.RLock()
	defer e.mx.RUnlock()
	return len(e.subscribes) == 0
}

func (e *EventChannel) Stop() {
	e.isRunning.Store(false)
}

func (e *EventChannel) run() {
	for msg := range e.msgs {
		if !e.isRunning.Load() {
			return
		}

		e.mx.RLock()
		for sub := range e.subscribes {
			sub <- msg
		}
		e.mx.RUnlock()
	}
	logger.GetLogger().Info(context.Background(), "message channel stopped")
}
