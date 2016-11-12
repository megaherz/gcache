package gcache

import (
	"time"
	"errors"
	"sync"
)

var ErrKeyNotFound = errors.New("Key Not Found")

type item struct {
	value    interface{}
	ttl      time.Duration
	expireAt time.Time
}

type Cache struct {
	items map[string]*item
	mutex sync.RWMutex
}

func (c*Cache) getItem(key string) (*item, bool) {
	if item, ok := c.items[key]; ok {
		expired := item.expireAt.Before(time.Now())

		if (expired) {
			return nil, false
		}

		return item, true
	}
	return nil, false
}

func (c*Cache) evict() {

	//Not efficient since full scan: rewrite!

	now := time.Now()
	for key, item := range c.items {
		if (item.expireAt.Before(now)) {
			//It's safe. From the spec: If map entries that have not yet been reached are removed during iteration, the corresponding iteration values will not be produced
			delete(key, item)
		}
	}
}

// Get the value of key.
// If the key does not exist the special value nil is returned
func (c*Cache) Get(key string) (interface{}, error) {
	c.mutex.RLock()
	if item, ok := c.getItem(key); ok {
		c.mutex.RUnlock()
		return item.value, nil
	}

	c.mutex.RUnlock()

	return nil, ErrKeyNotFound
}

// Get the Tll (Time to live) of the key
func (c*Cache) Ttl(key string) (time.Duration, error) {
	c.mutex.RLock()
	item, exists := c.getItem(key)
	if !exists {
		c.mutex.RUnlock()
		return -1, ErrKeyNotFound
	}

	c.mutex.RUnlock()
	return item.ttl, nil
}

// Set key to hold the value.
// If key already holds a value, it is overwritten, regardless of its type.
func (c*Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mutex.Lock()
	c.set(key, value, ttl)
	c.mutex.Unlock()
}

func (c*Cache) set(key string, value interface{}, ttl time.Duration) {
	expireAt := time.Now().Add(ttl)
	c.items[key] = &item{
		value: value,
		ttl: ttl,
		expireAt: expireAt,
	}

	c.evict()
}

// Update the value of the key
func (c*Cache) Update(key string, value interface{}) error {
	c.mutex.Lock()
	item, exists := c.getItem(key)
	if !exists {
		c.mutex.Unlock()
		return ErrKeyNotFound
	}

	c.set(key, value, item.ttl)
	c.mutex.Unlock()

	return nil
}

// Update the value of the key as well as TTL (Time to live)
func (c*Cache) UpdateWithTll(key string, value interface{}, ttl time.Duration) error {

	c.mutex.Lock()
	_, exists := c.getItem(key)
	if !exists {
		c.mutex.Unlock()
		return ErrKeyNotFound
	}

	c.set(key, value, ttl)
	c.mutex.Unlock()

	return nil
}

func (c*Cache) Del(key string) (err error) {

	c.mutex.Lock()
	if _, ok := c.getItem(key); ok {
		delete(c.items, key)
		c.mutex.Unlock()
	} else {
		c.mutex.Unlock()
		err = ErrKeyNotFound
	}

	return err
}

func (c*Cache) Keys() ([]string) {

	c.mutex.Lock()
	keys := make([]string, len(c.items))
	i := 0
	for k := range c.items {
		keys[i] = k
		i++
	}

	c.mutex.Unlock()
	return keys
}

func NewCache() (*Cache) {
	return &Cache{
		items:make(map[string]*item),
	}
}