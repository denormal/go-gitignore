package gitignore

import(
    "sync"
)


type Cache interface {
    Set( string , Ignore )
    Get( string )           Ignore
} // Cache{}


type cache struct {
    _i      map[string]Ignore
    _lock   sync.Mutex
} // cache{}


func NewCache() Cache {
    return &cache{}
} // Cache()


func ( c *cache ) Set( path string , ig Ignore ) {
    if ig != nil {
        c._lock.Lock()
        c._i[ path ]    = ig
        c._lock.Unlock()
    }
} // Set()


func ( c *cache ) Get( path string ) Ignore {
    c._lock.Lock()
    _ig , _ok   := c._i[ path ]
    c._lock.Unlock()
    if _ok {
        return _ig
    } else {
        return nil
    }
} // Get()


// ensure cache supports the Cache interface
var global  Cache   = &cache{}
