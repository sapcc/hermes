package storage

import (
	"github.com/sapcc/hermes/pkg/data"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_MockStorage_EventDetail(t *testing.T) {
	eventDetail, error := Mock().GetEvent("d5eed458-6666-58ec-ad06-8d3cf6bafca1")

	assert.Nil(t, error)
	assert.Equal(t, eventDetail.ID, "d5eed458-6666-58ec-ad06-8d3cf6bafca1")
	assert.Equal(t, eventDetail.Type, "identity.project.deleted")
	assert.Equal(t, eventDetail.Time, "2017-05-02T12:02:46.726056+0000")
}

func Test_MockStorage_Events(t *testing.T) {
	eventsList, total, error := Mock().GetEvents(data.Filter{})

	assert.Nil(t, error)
	assert.Equal(t, total, 24)
	assert.Equal(t, len(eventsList), 3)
	assert.Equal(t, eventsList[0].Type, "identity.project.deleted")
	assert.Equal(t, eventsList[1].ID, "095056c9-4cbb-5200-af70-0977dbcf5000")
	assert.Equal(t, eventsList[2].Time, "2017-05-02T11:45:44.755215+0000")
}
