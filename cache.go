package gitignore

import (
	"sync"
)

type Cache interface {
	Set(string, Ignore)
	Get(string) Ignore
} // Cache{}

type cache struct {
	_i    map[string]Ignore
	_lock sync.Mutex
} // cache{}

// NewCache returns a Cache instance. This is a thread-safe, in-memory cache
// for Ignore instances.
func NewCache() Cache {
	return &cache{}
} // Cache()

// Set stores the Ignore ig against its path.
func (c *cache) Set(path string, ig Ignore) {
	if ig != nil {
		c._lock.Lock()
		c._i[path] = ig
		c._lock.Unlock()
	}
} // Set()

// Get attempts to retrieve an Ignore instance associated with the given path.
// If the path is not known nil is returned.
func (c *cache) Get(path string) Ignore {
	c._lock.Lock()
	_ig, _ok := c._i[path]
	c._lock.Unlock()
	if _ok {
		return _ig
	} else {
		return nil
	}
} // Get()

// ensure cache supports the Cache interface
var global Cache = &cache{}
