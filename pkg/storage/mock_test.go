package storage

import (
	"github.com/sapcc/hermes/pkg/data"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_MockStorage_EventDetail(t *testing.T) {
	eventDetail, error := Mock().GetEvent("d5eed458-6666-58ec-ad06-8d3cf6bafca1", "b3b70c8271a845709f9a03030e705da7")

	assert.Nil(t, error)
	assert.Equal(t, "d5eed458-6666-58ec-ad06-8d3cf6bafca1", eventDetail.Payload.ID)
	assert.Equal(t, "identity.project.deleted", eventDetail.EventType)
	assert.Equal(t, "2017-05-02T12:02:46.726056+0000", eventDetail.Payload.EventTime)
}

func Test_MockStorage_Events(t *testing.T) {
	eventsList, total, error := Mock().GetEvents(data.Filter{}, "b3b70c8271a845709f9a03030e705da7")

	assert.Nil(t, error)
	assert.Equal(t, total, 24)
	assert.Equal(t, len(eventsList), 3)
	assert.Equal(t, "identity.project.deleted", eventsList[0].EventType)
	assert.Equal(t, "095056c9-4cbb-5200-af70-0977dbcf5000", eventsList[1].Payload.ID)
	assert.Equal(t, "2017-05-02T11:45:44.755215+0000", eventsList[2].Payload.EventTime)
}
