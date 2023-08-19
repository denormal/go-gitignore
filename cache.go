// Copyright 2016 Denormal Limited
// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gitignore

import (
	"sync"
)

// Cache is the interface for the GitIgnore cache
type Cache interface {
	// Set stores the GitIgnore ignore against its path.
	Set(path string, ig GitIgnore)

	// Get attempts to retrieve an GitIgnore instance associated with the given
	// path. If the path is not known nil is returned.
	Get(path string) GitIgnore
}

// cache is the default thread-safe cache implementation
type cache struct {
	_i    map[string]GitIgnore
	_lock sync.Mutex
}

// NewCache returns a Cache instance. This is a thread-safe, in-memory cache
// for GitIgnore instances.
func NewCache() Cache {
	return &cache{}
} // Cache()

// Set stores the GitIgnore ignore against its path.
func (c *cache) Set(path string, ignore GitIgnore) {
	if ignore == nil {
		return
	}

	// ensure the map is defined
	if c._i == nil {
		c._i = make(map[string]GitIgnore)
	}

	// set the cache item
	c._lock.Lock()
	c._i[path] = ignore
	c._lock.Unlock()
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
var _ Cache = &cache{}
