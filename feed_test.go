package feeder

import (
	"testing"
	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
)

func TestNewFeed(t *testing.T) {

	tt := []struct {
		name    string
		size    int
		perPage int
	}{
		{name: "notifications", perPage: 2, size: 5 },
		{name: "news", perPage: 4, size: 6 },
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

			got := NewFeed(tc.name, tc.size, tc.perPage, client)

			if got.Name != tc.name {
				t.Errorf("got feed name %s but wanted %s", got.Name, tc.name)
			}

			if got.Size != tc.size && got.PerPage != tc.perPage {
				t.Errorf("configuration of activity size got %d wanted %d", got.Size, tc.size)
				t.Errorf("configuration of page size got %d wanted %d", got.PerPage, tc.perPage)
				return
			}

			if got.Activities == nil {
				t.Errorf("feed created with a nil activities slice - should not happen we want empty slice")
				return
			}

			if got.P == nil {
				t.Errorf("feed created with a nil Provider - should not be so send provider for feed")
				return
			}
		})
	}
}
