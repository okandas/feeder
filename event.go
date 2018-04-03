package feeder

import "time"

// Action is the actual event
type Action struct {
	Value  string    `json:"value"`
	At     time.Time `json:"at"`
	Method string    `json:"method"`
}

// Event are the user actions a user does within the app that appear on the feed
type Event struct {
	UserID   string    `json:"user_id"`
	Activity []Action  `json:"activity"`
	Count    int       `json:"count"`
	LastRead int64 	    `json:"last_read"`
}

// NewEvent creates and instantiates a new user event
func NewEvent(userID string) *Event {
	event := &Event{
		UserID:   userID,
		Activity: []Action{},
	}

	return event
}
