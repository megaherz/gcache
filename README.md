# Gcache
Redis-similar Cache with REST protocol in Golang

##Features
* Expiration support (Ttl - time to live)
* Pure Go implementation
* Thread safe
* REST protocol
* Auth support
* Client scalability (multiple servers share the key space) 
* Go client library

##Example of usage
####Keys
```go
    // Create new cache
	cache := gcache.NewCache()

	// Set value with ttl (time to live)
	cache.Set("key", "value", 1 * time.Second)

	// Get value from the cache
	value, err := cache.Get("key")
	
	// Get ttl
    ttl, err := cache.Ttl("key")
    
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
####Lists	
```go	
    // Create new cache
	cache := gcache.NewCache()
	
	// Left push
	cache.LPush("key", "value")
	
    // Right push                       
	cache.RPush("key", "value")    
	
	// Left pop
	value, err := cache.LPop("key")
	
    // Right pop                           
    value, err := cache.RPop("key")   
    
    // Range
    values, err := cache.LRange("key", 2, 10)    
```  	
####Hashes
```go	
    // Create new cache
	cache := gcache.NewCache()

	// Add/Update hash valuy
	cache.HSet("key", "hashKey", "some hash value")  
	
	// Get hash value
	value, err := cache.HGet("key", "hashKey") 
	
```

##Server
The server might be run with or without authentication
```go
 // Run the server without authentication
 server := server.NewServer()
 
 // Run the server with authentication support. 
 // Password should be passed into the NewServerWithAuth function
 server := server.NewServerWithAuth("pass")
```

##Notes
* Keys, List and Hashed share the same keys space. Therefore it's forbidden to create the same key for e.g. Keys and Lists 
* Arrays for LRANGE and KEYS are returned as csv (encoding/csv package). Json is not used to make the protocol simple
* Before running the client_test.go execute the server.sh script in the ./cmd/gcache folder. 
The script runs two cache servers on ports 8080 and 8081. The 8081 server is run with authentication psw=123

##Protocol
The server exposes a REST protocol. Authentication is optional and is implemented as Authorization header with password as a value.
If the authentication is failed a standard 401 status is returned in response.

### Get key
Http method: GET <br/>
Url: /keys?key={key} <br/>
#### Request
**key** - key to get - string - required
#### Response
| Status Code  |    Meaning     |       Notes          |
|--------------|----------------|----------------------|
|      200     |  Ok            | Body is the value  |
|      400     |  Bad Request   |                      |
|      401     |  Auth failed   |                      |  
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
|      401     |  Auth failed   |                      |  
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
|      401     |  Auth failed   |                      |  
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
|      401     |  Auth failed   |                      |  
|      404     |  Not Found     |                      |
|      500     |  Server error  |                      |


### List of keys
Http method: GET <br/>
Url: /keys <br/>
Internally implemented as a linked list
#### Response
| Status Code  |    Meaning     |          Notes                    |
|--------------|----------------|-----------------------------------|
|      200     |  Ok            | Body is the csv list of keys  |
|      401     |  Auth failed   |                                   |  
|      500     |  Server error  |                                   |


### Left Push data to list (LPUSH)
Http method: POST <br/>
Url: /lists?op=lpush&key={key}&value={value} <br/>
#### Request
**key** - key list name - string - required <br/>
**value** - value to push into the list - string - required <br/>

#### Response                                              
| Status Code  |    Meaning     |          Notes       |   
|--------------|----------------|----------------------|   
|      200     |  Ok            |                      |   
|      400     |  Bad Request   |                      |   
|      401     |  Auth failed   |                      |     
|      500     |  Server error  |                      |   

### Right Push data to list (RPUSH)
Http method: POST <br/>
Url: /lists?op=rpush&key={key}&value={value} <br/>
#### Request
**key** - key list name - string - required <br/>
**value** - value to push into the list - string - required <br/>

#### Response                                            
| Status Code  |    Meaning     |          Notes       | 
|--------------|----------------|----------------------| 
|      200     |  Ok            |                      | 
|      400     |  Bad Request   |                      | 
|      401     |  Auth failed   |                      | 
|      500     |  Server error  |                      | 


### Left Pop data from list (LPOP)
Http method: POST <br/>
Url: /lists?op=lpop&key={key} <br/>
Pop removes element from the list therefore Http POST is used
#### Request
**key** - key list name - string - required <br/>

#### Response                                            
| Status Code  |    Meaning     |          Notes       | 
|--------------|----------------|----------------------| 
|      200     |  Ok            | Body is the  poped value | 
|      400     |  Bad Request   |                      | 
|      401     |  Auth failed   |                      | 
|      404     |  Not Found     |                      |  
|      500     |  Server error  |                      | 


### Right Pop data from list (RPOP)
Http method: POST <br/>
Url: /lists?op=rpop&key={key} <br/>
Pop removes element from the list therefore Http POST is used
#### Request
**key** - key list name - string - required <br/>

#### Response                                             
| Status Code  |    Meaning     |          Notes       |  
|--------------|----------------|----------------------|  
|      200     |  Ok            | Body is the  poped value |
|      400     |  Bad Request   |                      |  
|      401     |  Auth failed   |                      |  
|      404     |  Not Found     |                      |  
|      500     |  Server error  |                      |  


### Range data from list (LRANGE)
Http method: GET <br/>
Url: /lists?op=range&key={key}&from={from}&to={to} <br/>

Out of range indexes will not produce an error. If start is larger than the end of the list, an empty list is returned. 
If stop is larger than the actual end of the list, Redis will treat it like the last element of the list. <br/>

Note that if you have a list of numbers from 0 to 100, LRANGE list 0 10 will return 11 elements, that is, the rightmost item is included. 


#### Request
**key** - key list name - string - required <br/>
**from** - from index in range - int - required <br/>
**to** - to index in range - int - required <br/>

#### Response                                                         
| Status Code  |    Meaning     |          Notes       |              
|--------------|----------------|----------------------|              
|      200     |  Ok            | Body is the cvs list of values |    
|      400     |  Bad Request   |                      |              
|      401     |  Auth failed   |                      |              
|      404     |  Not Found     |                      |              
|      500     |  Server error  |                      |              


### Get field of hash (HGET)
Http method: GET <br/>
Url: /hashes/key={key}&hashKey={hasKey}<br/>

#### Request
**key** - key hash name - string - required <br/>
**hashKey** - hash key field - string - required <br/>

#### Response                                                         
| Status Code  |    Meaning     |          Notes       |              
|--------------|----------------|----------------------|              
|      200     |  Ok            |                      |    
|      400     |  Bad Request   |Including when the  key in hash is not found  |              
|      401     |  Auth failed   |                      |              
|      404     |  Not Found     |                      |              
|      500     |  Server error  |                      |              


### Set field of hash (HSET)
Http method: POST <br/>
Url: /hashes/key={key}&hashKey={hasKey}&value={value} <br/>

#### Request
**key** - key list name - string - required <br/>
**hashKey** - hash key field - string - required <br/>
**value** - hash value field - string - required <br/>

#### Response                                              
| Status Code  |    Meaning     |          Notes       |   
|--------------|----------------|----------------------|   
|      200     |  Ok            |                      |   
|      400     |  Bad Request   |                      |   
|      401     |  Auth failed   |                      |    
|      500     |  Server error  |                      |   

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
* Multiple set and get support on hashes
* Optimize number of locks since as for now when a new key is inserted the whole cache is blocked
* Return in response a reason what exactly is not valid on BadRequest(400) 
* Implement the Ttl method in client. The method should return ttl by the given key
* More client unit tests




