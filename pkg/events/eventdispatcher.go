package events

import (
	"errors"
	"sync"
)

var errorEventAlreadyRegistered = errors.New("event already registered")

type EventDispatcher struct {
	Handlers map[string][]EventHandlerInterface
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		Handlers: make(map[string][]EventHandlerInterface),
	}
}

func (e *EventDispatcher) Register(eventName string, handler EventHandlerInterface) error {
	if _, ok := e.Handlers[eventName]; ok {
		for _, h := range e.Handlers[eventName] {
			if h == handler {
				return errorEventAlreadyRegistered
			}
		}
	}

	e.Handlers[eventName] = append(e.Handlers[eventName], handler)
	return nil
}

func (e *EventDispatcher) Clear() error {
	e.Handlers = make(map[string][]EventHandlerInterface)
	return nil
}

func (e *EventDispatcher) Has(eventName string, handler EventHandlerInterface) bool {
	if _, ok := e.Handlers[eventName]; ok {
		for _, h := range e.Handlers[eventName] {
			if h == handler {
				return true
			}
		}
	}
	return false
}

func (e *EventDispatcher) Dispatch(eventInterface EventInterface) error {
	if handlers, ok := e.Handlers[eventInterface.GetName()]; ok {
		wg := &sync.WaitGroup{}
		for _, handler := range handlers {
			wg.Add(1)
			go handler.Handler(eventInterface, wg)
		}
		wg.Wait()
	}
	return nil
}

func (e *EventDispatcher) Remove(eventName string, handler EventHandlerInterface) {
	if _, ok := e.Handlers[eventName]; ok {
		for i, h := range e.Handlers[eventName] {
			if h == handler {
				e.Handlers[eventName] = append(e.Handlers[eventName][:i], e.Handlers[eventName][i+1:]...)
			}
		}
	}
}
