package cache


type Cache interface {
	// Adds a value to the cache, returns true if addition was successful.
	Add(key, value interface{}) bool

	// Returns key's value from cache.
	Get(key interface{}) (value interface{})

	// Checks if a key exists in cache.
	Contains(key interface{}) (ok bool)

	// Removes a key from the cache.
	Remove(key interface{}) bool

	// Returns a slice of the keys in the cache.
	Keys() []interface{}

	// Clears all cache entries
	Clear()
}