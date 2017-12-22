package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_MockStorage_EventDetail(t *testing.T) {
	// both params are ignored
	eventDetail, error := Mock{}.GetEvent("d5eed458-6666-58ec-ad06-8d3cf6bafca1", "b3b70c8271a845709f9a03030e705da7")
	assert.Nil(t, error)
	tt := []struct {
		name     string
		jsonPath string
		expected string
	}{
		{"eventDetail.ID", eventDetail.ID, "7be6c4ff-b761-5f1f-b234-f5d41616c2cd"},
		{"eventDetail.Action", eventDetail.Action, "create/role_assignment"},
		{"eventDetail.EventTime", eventDetail.EventTime, "2017-11-17T08:53:32.667973+00:00"},
		{"eventDetail.Outcome", eventDetail.Outcome, "success"},
		{"eventDetail.EventType", eventDetail.EventType, "activity"},
		{"eventDetail.Attachments[0].Name", eventDetail.Attachments[0].Name, "role_id"},
		{"eventDetail.Reason.ReasonType", eventDetail.Reason.ReasonType, "HTTP"},
		{"eventDetail.Reason.ReasonCode", eventDetail.Reason.ReasonCode, "409"},
		{"eventDetail.Initiator.Name", eventDetail.Initiator.Name, "test_admin"},
		{"eventDetail.Target.Addresses[0].URL", eventDetail.Target.Addresses[0].URL, "https://network-3.example.com/v2.0/security-group-rules/uuid"},
	}

	for _, tc := range tt {
		t.Run(tc.jsonPath, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.jsonPath)
		})
	}
}

func Test_MockStorage_Events(t *testing.T) {
	eventsList, total, error := Mock{}.GetEvents(&EventFilter{}, "b3b70c8271a845709f9a03030e705da7")

	assert.Nil(t, error)
	assert.Equal(t, total, 4)
	assert.Equal(t, len(eventsList), 4)
	assert.Equal(t, "success", eventsList[0].Outcome)
	assert.Equal(t, "f6f0ebf3-bf59-553a-9e38-788f714ccc46", eventsList[1].ID)
	assert.Equal(t, "2017-11-06T10:15:56.984390+00:00", eventsList[2].EventTime)
}

func Test_MockStorage__Attributes(t *testing.T) {
	attributesList, error := Mock{}.GetAttributes(&AttributeFilter{}, "b3b70c8271a845709f9a03030e705da7")

	assert.Nil(t, error)
	assert.Equal(t, len(attributesList), 6)
	assert.Equal(t, "compute/server", attributesList[0])
	assert.Equal(t, "network/floatingip", attributesList[4])
}
