package database

// Database interface
type Database interface {
	Init() error
	Get(key string) (string, error)
	Update(key, value string) error
	Delete(key string) error
	Contains(key string) bool
}

const (
	JsonDB = "kvstore.json" // json database file name

	MapType  = "map"  // to use in-memory database
	JsonType = "json" // to use json file as database

	MaxDbEntry = 1000 // maximum database entries
	MaxKeySz   = 16   // maximum length of key
	MaxValSz   = 32   // maximum length of value
)
