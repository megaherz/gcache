package gcache

const NoExpiration int =  -1

type Cache struct {
	keyValue map[string]interface{}
}

// Set key to hold the string value.
// If key already holds a value, it is overwritten, regardless of its type.
func (c* Cache) Set(key string, value interface{}, expire int) {
	c.keyValue[key] = value
}

func (c* Cache) Update(key string, value interface{}) {

}

// Get the value of key.
// If the key does not exist the special value nil is returned
func (c* Cache) Get(key string) (interface{}, bool) {
	if value, ok := c.keyValue[key]; ok {
		return value, true
	}
	return nil, false
}

func (c* Cache) Del (key string) (exists bool) {
	delete(c.keyValue, key)
	return true
}

func NewCache() (*Cache) {
	return &Cache{
		keyValue:make(map[string]interface{}),
	}
}