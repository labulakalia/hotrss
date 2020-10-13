package storage

// FeedStorager storage and read feed data
type FeedStorager interface {
	GetFeedData(string) (string, error)
	SaveFeedData(string, string) error
}

// DefaultStorage init mem storage
var DefaultStorage = NewMemStorage()
