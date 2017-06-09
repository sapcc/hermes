package hermes

import (
	"testing"

	"github.com/sapcc/hermes/pkg/identity"
	"github.com/sapcc/hermes/pkg/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_GetEvent(t *testing.T) {
	eventId := "d5eed458-6666-58ec-ad06-8d3cf6bafca1"
	event, err := GetEvent(eventId, "", identity.Mock{}, storage.Mock{})
	require.Nil(t, err)
	require.NotNil(t, event)
	assert.Equal(t, "d5eed458-6666-58ec-ad06-8d3cf6bafca1", event.Payload.ID)
	assert.NotEmpty(t, event.Payload.EventType)
	assert.NotEmpty(t, event.Payload.EventTime)
	assert.NotEmpty(t, event.Payload.Target.ID)
}

func Test_GetEvents(t *testing.T) {
	events, total, err := GetEvents(&Filter{}, "", identity.Mock{}, storage.Mock{})
	require.Nil(t, err)
	require.NotNil(t, events)
	assert.Equal(t, len(events), 3)
	assert.True(t, total >= len(events))
	for _, event := range events {
		assert.NotEmpty(t, event.ID)
		assert.NotEmpty(t, event.Type)
		assert.NotEmpty(t, event.Time)
	}
	assert.NotEqual(t, events[0].ID, events[1].ID)
	assert.NotEqual(t, events[0].ID, events[2].ID)
}
