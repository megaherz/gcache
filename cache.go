package gcache

import (
 	"time"
	"errors"
)

var ErrKeyNotFound = errors.New("Key Not Found")

type item struct {
	value interface{}
	ttl time.Duration
	expireAt time.Time
}

type Cache struct {
	items map[string]*item
}


func (c* Cache) getItem (key string) (*item, bool) {
	if item, ok := c.items[key]; ok {
		expired := item.expireAt.Before(time.Now())

		if (expired) {
			return nil, false
		}

		return item, true
	}
	return nil, false
}

// Get the value of key.
// If the key does not exist the special value nil is returned
func (c* Cache) Get(key string) (interface{}, error) {
	if item, ok := c.getItem(key); ok {
		return item.value, nil
	}

	return nil, ErrKeyNotFound
}

// Get the Tll (Time to live) of the key
func (c* Cache) Ttl(key string) (time.Duration, error) {
	item, exists := c.getItem(key)
	if !exists {
		return -1, ErrKeyNotFound
	}

	return item.ttl, nil
}

// Set key to hold the value.
// If key already holds a value, it is overwritten, regardless of its type.
func (c* Cache) Set(key string, value interface{}, ttl time.Duration) {

	expireAt := time.Now().Add(ttl)
	c.items[key] = &item {
		value: value,
		ttl: ttl,
		expireAt: expireAt,
	}
}

// Update the value of the key
func (c* Cache) Update(key string, value interface{}) error {
	item, exists := c.getItem(key)
	if !exists {
		return ErrKeyNotFound
	}

	c.Set(key, value, item.ttl)
	return nil
}

// Update the value of the key as well as TTL (Time to live)
func (c* Cache) UpdateWithTll(key string, value interface{}, ttl time.Duration) error{
	_, exists := c.getItem(key)
	if !exists {
		return ErrKeyNotFound
	}

	c.Set(key, value, ttl)
	return nil
}

func (c* Cache) Del (key string) (err error) {

	if _, ok := c.getItem(key); ok {
		delete(c.items, key)
	} else {
		err = ErrKeyNotFound
	}

	return err
}

func (c* Cache) Keys() ([]string) {

	keys := make([]string, len(c.items))
	i := 0
	for k := range c.items {
		keys[i] = k
		i++
	}
	return keys
}

func NewCache() (*Cache) {
	return &Cache{
		items:make(map[string] *item),
	}
}