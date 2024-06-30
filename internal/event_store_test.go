package internal

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/google/uuid"
)

type User struct {
	Name string
}

func TestEventStoreSuite(t *testing.T) {

	t.Run("returns true when adding a unique event", func(t *testing.T) {
		eventStore := NewEventStore()

		event1 := NewEvent(uuid.New().String(), 1, generateDefaultPayload())

		ok := eventStore.AddEvent(event1)

		assert.True(t, ok)
	})

	t.Run("returns false when adding an event with the same id and version", func(t *testing.T) {
		eventStore := NewEventStore()
		evtId := uuid.New().String()
		event1 := NewEvent(evtId, 1, generateDefaultPayload())
		eventStore.AddEvent(event1)

		event2 := NewEvent(evtId, 1, generatePayload(&User{Name: "Mike"}))
		ok := eventStore.AddEvent(event2)

		assert.False(t, ok)
	})

	t.Run("returns true when adding an event with the same id but different version", func(t *testing.T) {
		eventStore := NewEventStore()
		evtId := uuid.New().String()
		event1 := NewEvent(evtId, 1, generateDefaultPayload())
		eventStore.AddEvent(event1)

		event2 := NewEvent(evtId, 2, generateDefaultPayload())
		ok := eventStore.AddEvent(event2)

		assert.True(t, ok)
	})

	t.Run("retrieves a valid event", func(t *testing.T) {
		eventStore := NewEventStore()

		event1 := NewEvent(uuid.New().String(), 1, generateDefaultPayload())

		eventStore.AddEvent(event1)

		got, _ := eventStore.RetrieveEvent(event1.EventId, 1)
		want := event1

		assert.Equal(t, want, got)

	})

	t.Run("retrieves a valid stream", func(t *testing.T) {
		eventStore := NewEventStore()

		eventId := uuid.New().String()
		event1 := NewEvent(eventId, 1, generatePayload(&User{Name: "John"}))
		eventStore.AddEvent(event1)
		event2 := NewEvent(eventId, 2, generatePayload(&User{Name: "John H."}))
		eventStore.AddEvent(event2)
		event3 := NewEvent(eventId, 3, generatePayload(&User{Name: "John H. Watson"}))
		eventStore.AddEvent(event3)

		got, _ := eventStore.RetrieveEventStream(eventId)

		assert.Contains(t, got, event1)
		assert.Contains(t, got, event2)
		assert.Contains(t, got, event3)
	})

}

func generatePayload(v any) string {
	b, err := json.Marshal(v)

	if err != nil {
		panic(err)
	}

	return string(b)
}

func generateDefaultPayload() string {
	user := &User{Name: "Frank"}

	return generatePayload(user)
}
