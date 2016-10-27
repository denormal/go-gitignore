package gitignore

import (
	"sync"
)

// Cache is the interface for the GitIgnore cache
type Cache interface {
	Set(string, GitIgnore)
	Get(string) GitIgnore
} // Cache{}

// cache is the default thread-safe cache implementation
type cache struct {
	_i    map[string]GitIgnore
	_lock sync.Mutex
} // cache{}

// NewCache returns a Cache instance. This is a thread-safe, in-memory cache
// for GitIgnore instances.
func NewCache() Cache {
	return &cache{}
} // Cache()

// Set stores the GitIgnore ignore against its path.
func (c *cache) Set(path string, ignore GitIgnore) {
	if ignore != nil {
		c._lock.Lock()
		c._i[path] = ignore
		c._lock.Unlock()
	}
} // Set()

// Get attempts to retrieve an GitIgnore instance associated with the given
// path. If the path is not known nil is returned.
func (c *cache) Get(path string) GitIgnore {
	c._lock.Lock()
	_ignore, _ok := c._i[path]
	c._lock.Unlock()
	if _ok {
		return _ignore
	} else {
		return nil
	}
} // Get()

// ensure cache supports the Cache interface
var global Cache = &cache{}
