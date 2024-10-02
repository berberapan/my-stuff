package models

import (
	"testing"
	"time"

	"github.com/berberapan/my-stuff/internal/assert"
)

func TestItemModel(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}
	db := newTestDB(t)
	m := ItemModel{db}

	testAllItems := []struct {
		name          string
		id            int
		expectedSlice []Item
		errorMsg      error
	}{
		{
			name:          "User with no items",
			id:            10,
			expectedSlice: []Item{},
			errorMsg:      nil,
		},
		{
			name: "User with item",
			id:   1,
			expectedSlice: []Item{
				{
					ID:              1,
					Name:            "Laptop",
					Description:     "15-inch screen",
					Accessories:     "Charger",
					Place:           "Office",
					AdditionalNotes: "yay -Syu",
					UserID:          1,
					CreatedAt:       time.Now(),
				},
			},
			errorMsg: nil,
		},
	}

	for _, tt := range testAllItems {
		t.Run(tt.name, func(t *testing.T) {
			items, err := m.AllItems(tt.id)
			assert.Equal(t, len(items), len(tt.expectedSlice))
			assert.Equal(t, err, tt.errorMsg)
		})
	}

	testInsert := []struct {
		name            string
		itemName        string
		description     string
		accessories     string
		place           string
		additionalNotes string
		userID          int
		errorMsg        error
	}{
		{
			name:            "Insert with all fields",
			itemName:        "example",
			description:     "big",
			accessories:     "circle",
			place:           "corner",
			additionalNotes: "no",
			userID:          1,
			errorMsg:        nil,
		},
		{
			name:            "Insert with just name",
			itemName:        "example2",
			description:     "",
			accessories:     "",
			place:           "",
			additionalNotes: "",
			userID:          1,
			errorMsg:        nil,
		},
	}

	for _, tt := range testInsert {
		t.Run(tt.name, func(t *testing.T) {
			err := m.Insert(tt.itemName, tt.description, tt.accessories, tt.place, tt.additionalNotes, tt.userID)
			assert.Equal(t, err, tt.errorMsg)
		})
	}
}
