package feeder

import (
	"log"
	"math"
	"strconv"

	"github.com/go-redis/redis"
)

// Redis is a wrapper that implements the Backend Interface and represents a connection to a Redis server
type Redis struct {
	C *redis.Client
}

// HealthCheck checks if our redis server is up and running
func (r Redis) HealthCheck() (string, error) {
	response, err := r.C.Ping().Result()
	return response, err
}

// Store stores an event for a user
func (r Redis) Store(user, value string, at int64) (int64, error) {
	userData := user + ".data"
	// add event
	addResponse, err := r.C.ZAdd(userData, redis.Z{Score: float64(at), Member: value}).Result()

	if err != nil {
		log.Println(err)
		return 0, err
	}
	maxSize := int64(17)

	// check if events aren't exceeding 17 (if 18 remove the oldest event) and response will be 1
	// and events should now be 17
	removeResponse, err := r.C.ZRemRangeByRank(userData, 0, -maxSize-1).Result()

	if err != nil {
		log.Println(err)
		return 0, err
	}

	// increase
	userMeta := user + ".meta"

	change := addResponse - removeResponse
	_, err = r.C.HIncrBy(userMeta, "total_count", change).Result()

	if err != nil {
		log.Println(err)
		return 0, err
	}

	_, err = r.C.HIncrBy(userMeta, "unread_count", change).Result()

	if err != nil {
		log.Println(err)
		return 0, err
	}

	return change, nil
}

// Delete removes an event for a user
func (r Redis) Delete(user, value string, at int64) (int64, error) {
	userData := user + ".data"

	result, err := r.C.ZRem(userData, value).Result()

	if err != nil {
		log.Printf("remove item %s \n", err)
		return 0, err
	}

	return result, nil
}

// Wipe wipes the users feed
func (r Redis) Wipe(user string) (int64, error) {
	userData := user + ".data"


	response, err := r.C.Del(userData).Result()

	if err != nil {
		log.Printf("error failed: wipe delete events %s \n", err)
		return 0, err
	}

	return response, err
}

// Paginate paginates events for the user when peek is true
// do not reset #last_read
func (r Redis) Paginate(user string, page, perPage int) ([]string, error) {
	userData := user + ".data"

	from := (page - 1) * perPage
	to := (page * perPage) - 1

	results, err := r.C.ZRevRange(userData, int64(from), int64(to)).Result()


	return results, err

}

// All returns all events for the user
func (r Redis) All(user string) ([]string, error) {
	userData := user + ".data"
	events, err := r.C.ZRevRange(userData, 0, -1).Result()
	return events, err
}

// Count return the total count of events of user
func (r Redis) Count(user string) (int, error) {
	userMeta := user + ".meta"
	response, err := r.C.HGet(userMeta, "total_count").Result()
	count, err := strconv.Atoi(response)

	if err != nil {
		return 0, err
	}
	return count, err

}

// LastRead is LastRead when feed was last paginated
func (r Redis) LastRead(user string) (int64, error) {
	userMeta := user + ".meta"
	response, err := r.C.HGet(userMeta, "last_read").Result()

	if err != nil {
		return 0, err
	}

	timestamp, err := strconv.Atoi(response)

	if err != nil {
		return 0, err
	}

	return int64(timestamp), err
}

// ResetLastRead resets last read time stamp
// and also reset unread count
func (r Redis) ResetLastRead(user string, at int64) (bool, error) {
	userMeta := user + ".meta"

	res, err := r.C.HSet(userMeta, "last_read", at).Result()

	return res, err
}

// UnRead returns the total count of un read feed items
func (r Redis) UnRead(user string, at int64) (int64, error) {
	userData := user + ".data"
	positiveInf := strconv.FormatFloat(math.Inf(1), 'f', -1, 64)
	lastReadTimeStamp := strconv.FormatInt(at, 10)

	count, err := r.C.ZCount(userData, lastReadTimeStamp, positiveInf).Result()

	if err != nil {
		return 0, err
	}
	return count, err
}

// RecalculateCount count recalculates the length of events
func (r Redis) RecalculateCount(user string) (int64, error) {
	userMeta := user + ".meta"
	userData := user + ".data"

	count, err := r.C.ZCard(userData).Result()

	if err != nil {
		return count, err
	}

	res, err := r.C.HSet(userMeta, "total_count", count).Result()

	var updated int64

	if res == false {
		updated = 0
	}

	return updated, err
}

// NewRedisClient returns an instance of the Redis struct
func NewRedisClient(c *redis.Client) Backend {
	client := Redis{
		C: c,
	}

	return client
}
