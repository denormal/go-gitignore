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

package gitignore_test

import (
	"testing"

	"github.com/ianlewis/go-gitignore"
)

func TestCache(t *testing.T) {
	// populate the cache with the defined tests
	_cache := gitignore.NewCache()
	for _k, _v := range _CACHETEST {
		_cache.Set(_k, _v)
	}

	// attempt to retrieve the values from the cache
	//		- if a GitIgnore instance is returned, ensure it is the correct
	//		  instance, and not some other instance
	for _k, _v := range _CACHETEST {
		_found := _cache.Get(_k)
		if _found != _v {
			t.Errorf("cache Get() mismatch; expected %v, got %v",
				_v, _found,
			)
		}
	}

	// ensure unknown cache keys return nil
	for _, _k := range _CACHEUNKNOWN {
		_found := _cache.Get(_k)
		if _found != nil {
			t.Errorf("cache.Get() unexpected return for key %q; "+
				"expected nil, got %v",
				_k, _found,
			)
		}
	}

	// ensure we can update the cache
	_ignore := null()
	for _k := range _CACHETEST {
		_cache.Set(_k, _ignore)
	}
	for _k := range _CACHETEST {
		_found := _cache.Get(_k)
		if _found != _ignore {
			t.Errorf("cache Get() mismatch; expected %v, got %v",
				_ignore, _found,
			)
		}
	}
} // TestCache()
