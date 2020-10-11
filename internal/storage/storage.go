package storage

// FeedStorager storage and read feed data
type FeedStorager interface {
	GetFeedData(string) ([]byte, error)
	SaveFeedData(string, []byte) error
}

// DefaultStorage init mem storage
var DefaultStorage = NewMemStorage()