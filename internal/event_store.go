package internal

import (
	"slices"
	"time"
)

type EventStore struct {
	Events map[string][]*Event
}

type Event struct {
	EventId      string
	EventVersion int
	Payload      string
	CreatedAt    time.Time
}

func NewEvent(eventId string, eventVersion int, payload string) *Event {
	return &Event{
		EventId:      eventId,
		EventVersion: eventVersion,
		Payload:      payload,
		CreatedAt:    time.Now(),
	}
}

func (e *EventStore) AddEvent(event *Event) bool {
	stream := e.Events[event.EventId]
	if stream == nil {
		e.Events[event.EventId] = make([]*Event, 0)
		e.Events[event.EventId] = append(e.Events[event.EventId], event)
		event.EventVersion = 1
		return true
	} else {
		if !slices.ContainsFunc(stream, func(evt *Event) bool { return evt.EventVersion == event.EventVersion }) {
			e.Events[event.EventId] = append(e.Events[event.EventId], event)
			return true
		}
		return false
	}
}

func (e *EventStore) RetrieveEvent(eventId string, eventVersion int) (*Event, bool) {
	if eventVersion <= 0 {
		return nil, false
	}

	event := e.Events[eventId][eventVersion-1]
	if event != nil {
		return event, true
	}
	return nil, false
}

func (e *EventStore) RetrieveEventStream(eventId string) ([]*Event, bool) {
	events := e.Events[eventId]
	if len(events) != 0 {
		return events, true
	}
	return make([]*Event, 0), false
}

func NewEventStore() *EventStore {
	return &EventStore{
		Events: make(map[string][]*Event),
	}
}
