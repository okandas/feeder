package feeder

import (
	"testing"
)

func TestNewEvent(t *testing.T) {
	tt := []struct {
		description string
		user        string
	}{
		{
			description: "creates a new event",
			user:        "okandas",
		},
	}

	for _, tc := range tt {

		t.Run(tc.description, func(t *testing.T) {

			got := NewEvent(tc.user)

			if got.UserID != tc.user {
				t.Errorf("event created with wrong user: got %s want %s", got.UserID, tc.user)
			}

			if got.Activity == nil {
				t.Errorf("event created with a nil activity slice - should not happen we want empty slice")
			}

		})
	}
}
