package feeder

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
	"math/rand"
)

func TestNewRedisClient(t *testing.T) {

	server, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer server.Close()

	c := redis.NewClient(&redis.Options{
		Addr: server.Addr(),
	})

	client := NewRedisClient(c)

	want := "PONG"
	got, err := client.HealthCheck()

	if err != nil {
		t.Errorf("health check error %s", err)
	}

	if got != want {
		t.Errorf("we did not get a valid response from the serve got %s want %s", got, want)
	}

}

func TestStore(t *testing.T) {

	server, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer server.Close()

	c := redis.NewClient(&redis.Options{
		Addr: server.Addr(),
	})

	client := NewRedisClient(c)

	user := "okandas"
	value := "storing event for user"
	at := time.Now().Unix()

	want := int64(1)
	got, err := client.Store(user, value, at)

	if err != nil {
		t.Errorf("store action error %s", err)
	}

	if got != want {
		t.Errorf("store failed to store an event for user got %d %d", got, want)
	}

}

func TestDelete(t *testing.T) {

	server, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer server.Close()

	c := redis.NewClient(&redis.Options{
		Addr: server.Addr(),
	})

	client := NewRedisClient(c)

	// test setup
	user := "okandas"
	value := "storing event for user"
	at := time.Now().Unix()
	client.Store(user, value, at)

	want := int64(1)
	got, err := client.Delete(user, value, at)

	if err != nil {
		t.Errorf("error %s", err)
		return
	}

	if got != want {
		t.Errorf("got %d wanted %d", got, want)
	}

}

func TestWipe(t *testing.T) {

	server, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer server.Close()

	c := redis.NewClient(&redis.Options{
		Addr: server.Addr(),
	})

	client := NewRedisClient(c)

	// test setup
	user := "okandas"
	value := "storing event for user"
	at := time.Now().Unix()
	client.Store(user, value, at)

	var want int64 = 1

	got, err  := client.Wipe(user)

	if err != nil {
		t.Errorf("error %s", err)
		return
	}

	if got != want {
		t.Errorf("got %d wanted %d", got, want)
	}

}

func TestAll(t *testing.T) {

	server, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer server.Close()

	c := redis.NewClient(&redis.Options{
		Addr: server.Addr(),
	})

	client := NewRedisClient(c)

	// test setup
	user := "okandas"
	value := "storing first event for user"
	at := time.Now().Unix()
	client.Store(user, value, at)

	value = "storing second event for user"
	at = time.Now().Unix()
	client.Store(user, value, at)

	want := 2
	got, err := client.All(user)

	if err != nil {
		t.Errorf("error %s", err)
		return
	}

	if len(got) != want {
		t.Errorf("got %d wanted %d", len(got), want)
	}

}

func TestPaginate(t *testing.T) {

	server, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer server.Close()

	c := redis.NewClient(&redis.Options{
		Addr: server.Addr(),
	})

	client := NewRedisClient(c)

	// test setup
	user := "okandas"
	values := []string{"first", "second", "third", "fourth", "fifth", "sixth"}

	for _, value := range values {
		at := time.Now().Unix()
		client.Store(user, value, at)
	}

	want := 5
	got, err := client.Paginate(user, 1, 5)

	if err != nil {
		t.Errorf("error %s", err)
		return
	}

	if len(got) != want {
		t.Errorf("got %d wanted %d", len(got), want)
	}

}

func TestTotalCount(t *testing.T) {

	server, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer server.Close()

	c := redis.NewClient(&redis.Options{
		Addr: server.Addr(),
	})

	client := NewRedisClient(c)

	// test setup
	user := "okandas"
	values := []string{"first", "second", "third", "fourth", "fifth", "sixth"}

	for _, value := range values {
		at := time.Now().Unix()
		client.Store(user, value, at)
	}

	want := 6
	got, err := client.Count(user)

	if err != nil {
		t.Errorf("error %s", err)
		return
	}

	if  got != want {
		t.Errorf("got %d wanted %d", got, want)
	}

}

func TestUnReadCount(t *testing.T) {

	server, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer server.Close()

	c := redis.NewClient(&redis.Options{
		Addr: server.Addr(),
	})

	client := NewRedisClient(c)

	// test setup
	user := "okandas"

	timeLastRead := time.Now().Add(3 * time.Second)


	_, err = client.ResetLastRead(user, timeLastRead.Unix())

	if err != nil {
		t.Errorf("error %s", err)
		return
	}

	values := []string{"first", "second", "third", "fourth", "fifth", "sixth"}

	for _, value := range values {
		addTime := time.Duration(rand.Int())

		at := time.Now().Add(addTime * time.Second).Unix()
		client.Store(user, value, at)
	}

	lastRead, err := client.LastRead(user)

	if err != nil {
		t.Errorf("error %s", err)
		return
	}

	var want int64 = 2
	got, err := client.UnRead(user, lastRead)

	if err != nil {
		t.Errorf("error %s", err)
		return
	}

	if  got != want {
		t.Errorf("got %d wanted %d", got, want)
	}

}

func TestResetLastRead(t *testing.T) {

	server, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer server.Close()

	c := redis.NewClient(&redis.Options{
		Addr: server.Addr(),
	})

	client := NewRedisClient(c)

	// test setup
	user := "okandas"


	want := time.Now().Unix()
	done, err := client.ResetLastRead(user, want)

	if err != nil {
		t.Errorf("error %s", err)
		return
	}

	reset := true

	if done != reset {
		t.Errorf("failed to reset last read")
		return
	}

	got, err  := client.LastRead(user)

	if err != nil {
		t.Errorf("error %s", err)
		return
	}


	if got != want {
		t.Errorf("failed to reset last time read got %d wanted %d", got, want)
	}

}

func TestRecalculateCount(t *testing.T) {

	server, err := miniredis.Run()

	if err != nil {
		panic(err)
	}
	defer server.Close()

	c := redis.NewClient(&redis.Options{
		Addr: server.Addr(),
	})

	client := NewRedisClient(c)

	// test setup
	user := "okandas"
	values := []string{"first", "second"}


	for _, value := range values {
		at := time.Now().Unix()
		client.Store(user, value, at)
	}

	var want int64 = 0


	got, err  := client.RecalculateCount(user)


	if err != nil {
		t.Errorf("error %s", err)
		return
	}


	if got != want {
		t.Errorf("update - failed to reset last time read got %v wanted %v", got, want)
	}

}



