package hermes

import (
	"testing"

	"github.com/sapcc/hermes/pkg/data"
	"github.com/sapcc/hermes/pkg/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_GetEvent(t *testing.T) {
	storage := storage.Mock()
	eventId := "d5eed458-6666-58ec-ad06-8d3cf6bafca1"
	event, err := GetEvent(eventId, storage)
	require.Nil(t, err)
	require.NotNil(t, event)
	assert.Equal(t, event.ID, "d5eed458-6666-58ec-ad06-8d3cf6bafca1")
	assert.NotEmpty(t, event.Type)
	assert.NotEmpty(t, event.Time)
}

func Test_GetEvents(t *testing.T) {
	storage := storage.Mock()
	events, total, err := GetEvents(storage, data.Filter{})
	require.Nil(t, err)
	require.NotNil(t, events)
	assert.Equal(t, len(events), 3)
	assert.True(t, total >= len(events))
	for i := range events {
		assert.NotEmpty(t, events[i].ID)
		assert.NotEmpty(t, events[i].Type)
		assert.NotEmpty(t, events[i].Time)
	}
	assert.NotEqual(t, events[0].ID, events[1].ID)
	assert.NotEqual(t, events[0].ID, events[2].ID)
}
