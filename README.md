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

##Notes
To run the server execute the server.sh script in the ./run directory

##Protocol

### Get key
Http method: GET <br/>
Url: /keys?key={key} <br/>
#### Request
**key** - key to get - string - required
#### Response
| Status Code  |    Meaning     |          Notes       |
|--------------|----------------|----------------------|
|      200     |  Ok            | Body contains value  |
|      400     |  Bad Request   |                      |
|      404     |  Not Found     |                      |
|      500     |  Server error  |                      |

```curl
curl -XGET http://localhost:8080/keys?key={key}

```

### Set key
Http method: POST <br/>
Url: /keys?key={key}&value={value}&ttl={ttl} <br/>
#### Request
**key** - key to get - string - required <br/>
**value** - value to set - string - required <br/>
**ttl** - time to live in _seconds_ - int - required <br/>


### Update key
Http method: PATCH <br/>
Url: /keys?key={key}&value={value}&ttl={ttl} <br/>
#### Request
**key** - key to get - string - required <br/>
**value** - value to set - string - required <br/>
**ttl** - time to live in _seconds_ - int - optional <br/>

### Delete key
Http method: DELETE <br/>
Url: /keys?key={key} <br/>
#### Request
**key** - key to get - string - required <br/>

### List of keys
Http method: GET <br/>
Url: /keys <br/>