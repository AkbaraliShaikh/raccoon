package connection

import (
	"testing"

	"github.com/odpf/raccoon/identification"
	"github.com/stretchr/testify/assert"
)

func TestConnectionPerGroup(t *testing.T) {
	t.Run("Should return all the group on the table with the count", func(t *testing.T) {
		table := NewTable(10)
		table.Store(identification.Identifier{ID: "user1", Group: "group1"})
		table.Store(identification.Identifier{ID: "user2", Group: "group1"})
		table.Store(identification.Identifier{ID: "user3", Group: "group1"})
		table.Store(identification.Identifier{ID: "user1", Group: "group2"})
		table.Store(identification.Identifier{ID: "user2", Group: "group2"})
		assert.Equal(t, map[string]int{"group1": 3, "group2": 2}, table.TotalConnectionPerGroup())
	})
}

func TestStore(t *testing.T) {
	t.Run("Should store new connection", func(t *testing.T) {
		table := NewTable(10)
		err := table.Store(identification.Identifier{ID: "user1", Group: ""})
		assert.NoError(t, err)
		assert.True(t, table.Exists(identification.Identifier{ID: "user1"}))
	})

	t.Run("Should return max connection reached error when connection is maxed", func(t *testing.T) {
		table := NewTable(0)
		err := table.Store(identification.Identifier{ID: "user1", Group: ""})
		assert.Error(t, err, errMaxConnectionReached)
	})

	t.Run("Should return duplicated error when connection already exists", func(t *testing.T) {
		table := NewTable(2)
		err := table.Store(identification.Identifier{ID: "user1", Group: ""})
		assert.NoError(t, err)
		err = table.Store(identification.Identifier{ID: "user1", Group: ""})
		assert.Error(t, err, errConnDuplicated)
	})

	t.Run("Should remove connection when identifier match", func(t *testing.T) {
		table := NewTable(10)
		table.Store(identification.Identifier{ID: "user1", Group: ""})
		table.Remove(identification.Identifier{ID: "user1", Group: ""})
		assert.False(t, table.Exists(identification.Identifier{ID: "user1", Group: ""}))
	})
}

func TestStoreEvent(t *testing.T) {
	t.Run("Should store new event for a connection", func(t *testing.T) {
		table := NewTable(10)
		table.Store(identification.Identifier{ID: "user1", Group: ""})
		table.StoreEvent(identification.Identifier{ID: "user1", Group: ""}, "event-id-1")

		assert.True(t, table.HasEvent(identification.Identifier{ID: "user1", Group: ""}, "event-id-1"))
	})

	t.Run("Should not store new event if the connection is not active", func(t *testing.T) {
		table := NewTable(10)
		table.StoreEvent(identification.Identifier{ID: "user1", Group: ""}, "event-id-1")

		assert.False(t, table.HasEvent(identification.Identifier{ID: "user1", Group: ""}, "event-id-1"))
	})

	t.Run("Should store multiple unique events for a connetion", func(t *testing.T) {
		table := NewTable(10)
		table.Store(identification.Identifier{ID: "user1", Group: ""})
		table.StoreEvent(identification.Identifier{ID: "user1", Group: ""}, "event-id-1")
		table.StoreEvent(identification.Identifier{ID: "user1", Group: ""}, "event-id-2")

		assert.True(t, table.HasEvent(identification.Identifier{ID: "user1", Group: ""}, "event-id-1"))
		assert.True(t, table.HasEvent(identification.Identifier{ID: "user1", Group: ""}, "event-id-2"))
	})

	t.Run("Should store multiple unique events for multiple connetion", func(t *testing.T) {
		table := NewTable(10)
		table.Store(identification.Identifier{ID: "user1", Group: ""})
		table.Store(identification.Identifier{ID: "user2", Group: ""})

		table.StoreEvent(identification.Identifier{ID: "user1", Group: ""}, "event-id-1")
		table.StoreEvent(identification.Identifier{ID: "user1", Group: ""}, "event-id-2")
		table.StoreEvent(identification.Identifier{ID: "user2", Group: ""}, "event-id-1")
		table.StoreEvent(identification.Identifier{ID: "user2", Group: ""}, "event-id-2")

		assert.True(t, table.HasEvent(identification.Identifier{ID: "user1", Group: ""}, "event-id-1"))
		assert.True(t, table.HasEvent(identification.Identifier{ID: "user1", Group: ""}, "event-id-2"))
		assert.True(t, table.HasEvent(identification.Identifier{ID: "user2", Group: ""}, "event-id-1"))
		assert.True(t, table.HasEvent(identification.Identifier{ID: "user2", Group: ""}, "event-id-2"))
	})

	t.Run("Should remove all the events if connetion is removed or not active", func(t *testing.T) {
		table := NewTable(10)

		table.Store(identification.Identifier{ID: "user1", Group: ""})
		table.Store(identification.Identifier{ID: "user2", Group: ""})

		table.StoreEvent(identification.Identifier{ID: "user1", Group: ""}, "event-id-1")
		table.StoreEvent(identification.Identifier{ID: "user1", Group: ""}, "event-id-2")
		table.StoreEvent(identification.Identifier{ID: "user2", Group: ""}, "event-id-1")
		table.StoreEvent(identification.Identifier{ID: "user2", Group: ""}, "event-id-2")

		table.Remove(identification.Identifier{ID: "user1", Group: ""})
		table.Remove(identification.Identifier{ID: "user2", Group: ""})

		assert.False(t, table.HasEvent(identification.Identifier{ID: "user1", Group: ""}, "event-id-1"))
		assert.False(t, table.HasEvent(identification.Identifier{ID: "user1", Group: ""}, "event-id-2"))
		assert.False(t, table.HasEvent(identification.Identifier{ID: "user2", Group: ""}, "event-id-1"))
		assert.False(t, table.HasEvent(identification.Identifier{ID: "user2", Group: ""}, "event-id-2"))
	})
}
