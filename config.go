package feeder

// Config is the configuration for a database used by each Feed
type Config struct {
	Engine string `default:"redis" required:"true"`
	Host   string `default:"localhost" required:"true" desc:"database Host URL"`
	Port   string `default:":6379" required:"true" desc:"address where our server will listen to"`
	Db     int    `default:"1" desc:"number of redis instances"`
	Debug  bool   `default:"false" desc:"config start in debug mode"`
}
