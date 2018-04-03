package feeder

// Feeds are multiple feeds within the application
type Feeds struct {
	Feeds []Feed `json:"feeds"`
}

// Feed represents a feed
type Feed struct {
	Name          string   `json:"name"`
	Size          int      `json:"max_size"` // items can be in the feed
	PerPage       int      `json:"per_page"`
	Activities    []*Event `json:"activities"`
	P	  		  Backend `json:"-"`
}

// NewFeed instantiates a new feed struct
func NewFeed(name string, size, perPage int,  b Backend) *Feed {
	return &Feed{
		Name:          name,
		Size:          size,
		PerPage:       perPage,
		Activities:    []*Event{},
		P: 	   b,
	}

}

//TODO: should also register our feed in the register and then have a way of accessing that feed with a method
