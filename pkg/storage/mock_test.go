package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_MockStorage_EventDetail(t *testing.T) {
	eventDetail, error := Mock{}.GetEvent("d5eed458-6666-58ec-ad06-8d3cf6bafca1", "b3b70c8271a845709f9a03030e705da7")

	assert.Nil(t, error)
	assert.Equal(t, "7be6c4ff-b761-5f1f-b234-f5d41616c2cd", eventDetail.ID)
	assert.Equal(t, "create/role_assignment", eventDetail.Action)
	assert.Equal(t, "2017-11-17T08:53:32.667973+0000", eventDetail.EventTime)
	assert.Equal(t, "success", eventDetail.Outcome)
}

func Test_MockStorage_Events(t *testing.T) {
	eventsList, total, error := Mock{}.GetEvents(&EventFilter{}, "b3b70c8271a845709f9a03030e705da7")

	assert.Nil(t, error)
	assert.Equal(t, total, 4)
	assert.Equal(t, len(eventsList), 4)
	assert.Equal(t, "success", eventsList[0].Outcome)
	assert.Equal(t, "f6f0ebf3-bf59-553a-9e38-788f714ccc46", eventsList[1].ID)
	assert.Equal(t, "2017-11-06T10:15:56.984390+0000", eventsList[2].EventTime)
}
