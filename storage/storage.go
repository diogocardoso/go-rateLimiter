package storage

type Storage interface {
	Increment(key string) (int64, error)
	Block(key string, seconds int64) error
	IsBlocked(key string) (bool, error)
	Reset(key string) error
}
