# gcache
Cache with REST protocol in Golang

##Features
* Expiration support (Ttl - time to live)
* Pure Go implementation
* Thread safe
* REST protocol
* Go client library

##Example Usage
```go
    // Create new cache
	cache := gcache.NewCache()

	// Set value with ttl (time to live)
	cache.Set("key", "value", 1 * time.Second)

	// Get value from the cache
	value, err := cache.Get("key")

	// Delete key from the cache
	err := cache.Del("key")

	// Update key with a new value
	cache.Update("key", "value")

	// Update key with a new value and tll
	cache.UpdateWithTll("key", "value", 3 * time.Second)
	
	// Get all keys
	keys := cache.Keys()
	
	// Get number of items in the cache
	count := cache.Count()
```

###Protocol

#### Get key
Http method: GET
Url: /keys?key={key}
##### Request
**key** - key to get - string - required

#### Set key
Http method: POST
Url: /keys?key={key}&value={value}&ttl={ttl}
##### Request
**key** - key to get - string - required
**value** - value to set - string - required
**ttl** - time to live in _seconds_ - int - required


#### Update key
Http method: PATCH
Url: /keys?key={key}&value={value}&ttl={ttl}
##### Request
**key** - key to get - string - required
**value** - value to set - string - required
**ttl** - time to live in _seconds_ - int - optional

#### Delete key
Http method: DELETE
Url: /keys?key={key}
##### Request
**key** - key to get - string - required

#### List of keys
Http method: GET
Url: /keys