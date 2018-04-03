package feeder

import (
	"log"
	"time"
)

// Activity represents a users activity in each feed

// Activity is an activity object that has a users activities per feed
type Activity struct {
	UserID   string   `json:"user_id"`
	Feed     *Feed    `json:"feed"`
}

// Wipe removes all events made by a user/app in a feed
func (a Activity) Wipe() (int64, error) {

	res, err := a.Feed.P.Wipe(a.UserID)

	return res, err
}

func (a Activity) Store(value string, at int64) (int64, error) {

	res, err := a.Feed.P.Store(a.UserID, value, at)

	return res, err
}

func (a Activity) Count() (int, error)  {

	res, err := a.Feed.P.Count(a.UserID)

	return res, err
}

func (a Activity) ResetLastRead(at int64) (bool, error) {

	res, err := a.Feed.P.ResetLastRead(a.UserID, at)

	return res, err
}



func (a Activity) UnreadCount() (int64, error) {

	lastRead, err := a.Feed.P.LastRead(a.UserID)

	if err != nil {
		log.Printf("failed to get last read %s \n", err)
	}

	res, err := a.Feed.P.UnRead(a.UserID, lastRead)

	return res, err
}

func (a Activity) Paginate(page, perPage int) ([]string, error) {

	res, err := a.Feed.P.Paginate(a.UserID, page, perPage)
	if err == nil {
		res, err := a.Feed.P.ResetLastRead(a.UserID, time.Now().Unix())

		if err == nil {
			log.Printf("user has read feed %v", res)
		}

	}

	return res, err
}

func (a Activity) All() ([]string, error) {

	res, err := a.Feed.P.All(a.UserID)

	return res, err
}


// NewActivity instantiates a new user activity
func NewActivity(id string, feed *Feed) *Activity {
	return &Activity{
		UserID: id,
		Feed:   feed,
	}
}

