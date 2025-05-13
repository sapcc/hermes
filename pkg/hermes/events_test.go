// Copyright 2022 SAP SE
// SPDX-FileCopyrightText: 2025 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package hermes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sapcc/hermes/pkg/storage"
)

func Test_GetEvent(t *testing.T) {
	eventID := "7be6c4ff-b761-5f1f-b234-f5d41616c2cd"
	event, err := GetEvent(eventID, "", storage.Mock{})
	require.Nil(t, err)
	require.NotNil(t, event)
	assert.Equal(t, "7be6c4ff-b761-5f1f-b234-f5d41616c2cd", event.ID)
	assert.NotEmpty(t, event.Outcome)
	assert.NotEmpty(t, event.EventTime)
	assert.NotEmpty(t, event.Action)
}

func Test_GetEvents(t *testing.T) {
	events, total, err := GetEvents(&EventFilter{}, "", storage.Mock{})
	require.Nil(t, err)
	require.NotNil(t, events)
	assert.Equal(t, len(events), 4)
	assert.True(t, total >= len(events))
	for _, event := range events {
		assert.NotEmpty(t, event.ID)
		assert.NotEmpty(t, event.Outcome)
		assert.NotEmpty(t, event.Time)
		assert.NotEmpty(t, event.Initiator.ID)
		assert.NotEmpty(t, event.Initiator.Name)
		assert.NotEmpty(t, event.Initiator.TypeURI)
	}
	assert.NotEqual(t, events[0].ID, events[1].ID)
	assert.NotEqual(t, events[0].ID, events[2].ID)
}

func Test_GetAttributes(t *testing.T) {
	attributes, err := GetAttributes(&AttributeFilter{}, "", storage.Mock{})
	require.Nil(t, err)
	require.NotNil(t, attributes)
	assert.Equal(t, len(attributes), 6)
}
