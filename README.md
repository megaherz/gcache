# Gcache
Rest-similar Cache with REST protocol in Golang

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
* To run the server execute the server.sh script in the ./run directory
* Keys, List and Hashed share the same keys space. Therefore it's forbidden to create the same key for e.g. Keys and Lists 

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


### Set key
Http method: POST <br/>
Url: /keys?key={key}&value={value}&ttl={ttl} <br/>
#### Request
**key** - key to get - string - required <br/>
**value** - value to set - string - required <br/>
**ttl** - time to live in _seconds_ - int - required <br/>
#### Response
| Status Code  |    Meaning     |          Notes       |
|--------------|----------------|----------------------|
|      200     |  Ok            |                      |
|      400     |  Bad Request   |                      |
|      500     |  Server error  |                      |



### Update key
Http method: PATCH <br/>
Url: /keys?key={key}&value={value}&ttl={ttl} <br/>
#### Request
**key** - key to get - string - required <br/>
**value** - value to set - string - required <br/>
**ttl** - time to live in _seconds_ - int - optional <br/>
#### Response
| Status Code  |    Meaning     |          Notes       |
|--------------|----------------|----------------------|
|      200     |  Ok            |                      |
|      400     |  Bad Request   |                      |
|      404     |  Not Found     |                      |
|      500     |  Server error  |                      |


### Delete key
Http method: DELETE <br/>
Url: /keys?key={key} <br/>
#### Request
**key** - key to get - string - required <br/>
#### Response
| Status Code  |    Meaning     |          Notes       |
|--------------|----------------|----------------------|
|      200     |  Ok            |                      |
|      400     |  Bad Request   |                      |
|      404     |  Not Found     |                      |
|      500     |  Server error  |                      |


### List of keys
Http method: GET <br/>
Url: /keys <br/>
Internally implemented as a linked list
#### Response
| Status Code  |    Meaning     |          Notes                |
|--------------|----------------|-------------------------------|
|      200     |  Ok            | Body contains a list of keys  |
|      500     |  Server error  |                               |


### Left Push data to list (LPUSH)
Http method: POST <br/>
Url: /lists/lpush?list={list}&value={value} <br/>
#### Request
**list** - list name - string - required <br/>
**value** - value to push into the list - string - required <br/>

### Right Push data to list (RPUSH)
Http method: POST <br/>
Url: /lists/rpush?list={list}&value={value} <br/>
#### Request
**list** - list name - string - required <br/>
**value** - value to push into the list - string - required <br/>

### Left Pop data from list (LPOP)
Http method: POST <br/>
Url: /lists/lpop?list={list} <br/>
Pop removes element from the list therefore Http POST is used
#### Request
**list** - list name - string - required <br/>


### Right Pop data from list (RPOP)
Http method: POST <br/>
Url: /lists/rpop?list={list} <br/>
Pop removes element from the list therefore Http POST is used
#### Request
**list** - list name - string - required <br/>

### Range data from list (LRANGE)
Http method: GET <br/>
Url: /lists/range?list={list}&from={from}&to={to} <br/>

Out of range indexes will not produce an error. If start is larger than the end of the list, an empty list is returned. 
If stop is larger than the actual end of the list, Redis will treat it like the last element of the list. <br/>

Note that if you have a list of numbers from 0 to 100, LRANGE list 0 10 will return 11 elements, that is, the rightmost item is included. 

#### Request
**list** - list name - string - required <br/>
**from** - from index in range - int - required <br/>
**to** - to index in range - int - required <br/>

### Get field of hash (HGET)
Http method: GET <br/>
Url: /hashes/key={key}<br/>

### Set field of hash (HSET)
Http method: POST <br/>
Url: /hashes/key={key}&value={value} <br/>


##Performance
```go
func BenchmarkCache_SetGet(b *testing.B) {

	cache := NewCache()

	for n := 0; n < b.N; n++ {

		key := "key" + strconv.Itoa(n)

		cache.Set(key, "val", 1 * time.Second)
		_, err := cache.Get(key)

		if (err != nil){
			b.Error("Failed to GET", key)
		}
	}
}
```

Result: <br />
BenchmarkCache_SetGet-4   |   1000000   |     2032 ns/op

##TODO
* Validation that a key does not contain whitespaces
* Multiple set and get support on hashes
* Arrays for LRANGE and KEYS are returned as \n separated values. 
It's not safe if a value in List contains \n. Need to implement something as Array-Reply in Redis
http://redis.io/topics/protocol#array-reply 



