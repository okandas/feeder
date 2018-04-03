package feeder

// Backend is the backend of our feeds
type Backend interface {
	Store(user, value string, at int64) (int64, error)
	Delete(user, value string, at int64) (int64, error)
	HealthCheck() (string, error)
	Wipe(user string) (int64, error)
	All(user string) ([]string, error)
	Paginate(user string, page, perPage int) ([]string, error)
	Count(user string) (int, error)
	UnRead(user string, at int64) (int64, error)
	ResetLastRead(user string, at int64) (bool, error)
	LastRead(user string) (int64, error)
	RecalculateCount(user string) (int64, error)
}
