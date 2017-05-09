package hermes

import (
	"testing"

	"github.com/sapcc/hermes/pkg/data"
	"github.com/sapcc/hermes/pkg/storage"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setup() {
	viper.Set("hermes.keystone_driver", "mock")
	viper.Set("hermes.storage_driver", "mock")
}

func Test_GetEvent(t *testing.T) {
	setup()
	eventId := "d5eed458-6666-58ec-ad06-8d3cf6bafca1"
	event, err := GetEvent(eventId, storage.ConfiguredDriver())
	require.Nil(t, err)
	require.NotNil(t, event)
	assert.Equal(t, event.Payload.ID, "d5eed458-6666-58ec-ad06-8d3cf6bafca1")
	assert.NotEmpty(t, event.Payload.EventType)
	assert.NotEmpty(t, event.Payload.EventTime)
	assert.NotEmpty(t, event.Payload.Target.Name)
}

func Test_GetEvents(t *testing.T) {
	setup()
	events, total, err := GetEvents(storage.ConfiguredDriver(), data.Filter{})
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
