package feeder

import (

	"testing"
	"github.com/google/go-cmp/cmp"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
	"time"
)


func TestNewActivity(t *testing.T) {

	tt := []struct {
		name    string
		size int
		perPage int
	}{
		{ name: "notifications", size: 5,  perPage: 10 },
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			server, err := miniredis.Run()
			if err != nil {
				panic(err)
			}
			defer server.Close()

			c := redis.NewClient(&redis.Options{
				Addr: server.Addr(),
			})

			client := NewRedisClient(c)

			userID := "okandas"
			feed := NewFeed(tc.name, tc.size, tc.perPage, client)

			got := NewActivity(userID, feed)

			if !cmp.Equal(got.Feed, feed) {
				t.Errorf("user activity created on wrong feed got %+v want %+v", got.Feed, feed)
				return
			}

			if got.UserID != userID {
				t.Errorf("user activity created on with wrong ID got %+v want %+v", got.UserID, userID)
				return
			}

		})
	}
}

func TestActivityWipe(t *testing.T) {

	server, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer server.Close()

	c := redis.NewClient(&redis.Options{
		Addr: server.Addr(),
	})

	client := NewRedisClient(c)

	userID := "okandas"

	feed := NewFeed("notifications", 10, 15, client)

	activity := NewActivity(userID, feed)

	got, err  := activity.Wipe()

	if err != nil {
		t.Errorf("activity wipe error %s", err)
	}

	var want int64 =  0

	if got != want {
		t.Errorf("failed to wipe user feed using activity got %v but wanted %v", got, want)
	}

}

func TestActivityStore(t *testing.T) {

	server, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer server.Close()

	c := redis.NewClient(&redis.Options{
		Addr: server.Addr(),
	})

	client := NewRedisClient(c)

	userID := "okandas"


	feed := NewFeed("notifications", 10, 15, client)

	activity := NewActivity(userID, feed)

	got, err := activity.Store("this is the first value store", time.Now().Unix())

	expected := int64(1)

	if got != expected {
		t.Errorf("failed to store to  user feed using activity got %v but wanted %v", got, expected)
	}

	count, err  := activity.Count()

	if err != nil {
		t.Errorf("activity count error %s", err)
		return
	}


	want :=  1

	if count != want {
		t.Errorf("failed to count user feed events after storing one item got %v but wanted %v", count, want)
		return
	}

}

func TestActivityCount(t *testing.T) {
	// test setup
	server, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer server.Close()

	c := redis.NewClient(&redis.Options{
		Addr: server.Addr(),
	})

	client := NewRedisClient(c)

	userID := "okandas"


	feed := NewFeed("notifications", 10, 15, client)

	activity := NewActivity(userID, feed)

	activity.Store("this is the first value store", time.Now().Unix())
	activity.Store("this is the second value store", time.Now().Unix())



	// testing
	count, err  := activity.Count()

	if err != nil {
		t.Errorf("activity count error %s", err)
		return
	}

	want :=  2

	// test case
	if count != want {
		t.Errorf("failed to count user feed events after storing one item got %v but wanted %v", count, want)
		return
	}

}

func TestActivityUnread(t *testing.T) {
	// test setup
	server, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer server.Close()

	c := redis.NewClient(&redis.Options{
		Addr: server.Addr(),
	})

	client := NewRedisClient(c)

	userID := "okandas"



	feed := NewFeed("notifications", 10, 15, client)

	activity := NewActivity(userID, feed)

	unix := time.Now().Add(time.Duration(-2) * time.Minute).Unix()

	activity.ResetLastRead(unix)

	activity.Store("this is the first value store", time.Now().Add(time.Duration(-3) * time.Minute).Unix())
	activity.Store("this is the second value store", time.Now().Unix())
	activity.Store("this is the third value store", time.Now().Unix())
	activity.Store("this is the four value store", time.Now().Unix())


	want := int64(3)

	got, err := activity.UnreadCount()

	if err != nil {
		t.Errorf("Unread count error %s", err)
	}

	if got != want {
		t.Errorf("got %d want %d", got, want)
	}

}

func TestActivityPagination(t *testing.T) {
	// test setup
	server, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer server.Close()

	c := redis.NewClient(&redis.Options{
		Addr: server.Addr(),
	})

	client := NewRedisClient(c)

	userID := "okandas"



	feed := NewFeed("notifications", 10, 15, client)

	activity := NewActivity(userID, feed)


	activity.Store("this is the first value store", time.Now().Unix())
	activity.Store("this is the second value store", time.Now().Unix())
	activity.Store("this is the third value store", time.Now().Unix())
	activity.Store("this is the four value store", time.Now().Unix())
	activity.Store("this is the five value store", time.Now().Unix())
	activity.Store("this is the six value store", time.Now().Unix())
	activity.Store("this is the seven value store", time.Now().Unix())


	want := 5

	got, err := activity.Paginate(1, 5)

	if err != nil {
		t.Errorf("pagination error %s", err)
	}

	if len(got) != want {
		t.Errorf("got %d want %d", got, want)
	}

}

func TestActivityAll(t *testing.T) {
	// test setup
	server, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer server.Close()

	c := redis.NewClient(&redis.Options{
		Addr: server.Addr(),
	})

	client := NewRedisClient(c)

	userID := "okandas"



	feed := NewFeed("notifications", 10, 15, client)

	activity := NewActivity(userID, feed)


	activity.Store("this is the first value store", time.Now().Unix())
	activity.Store("this is the second value store", time.Now().Unix())
	activity.Store("this is the third value store", time.Now().Unix())
	activity.Store("this is the four value store", time.Now().Unix())
	activity.Store("this is the five value store", time.Now().Unix())
	activity.Store("this is the six value store", time.Now().Unix())
	activity.Store("this is the seven value store", time.Now().Unix())


	want := 7

	got, err := activity.All()

	if err != nil {
		t.Errorf("all error %s", err)
	}

	if len(got) != want {
		t.Errorf("got %d want %d", got, want)
	}

}

