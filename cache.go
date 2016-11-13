package gcache

import (
	"time"
	"errors"
	"sync"
	"container/heap"
	"runtime"
)

const DefaultEvictionInterval time.Duration = 1 * time.Second

var ErrKeyNotFound = errors.New("Key Not Found")

type item struct {
	key      string
	value    interface{}
	ttl      time.Duration
	expireAt time.Time
}

type priorityQueue []*item

func (pq priorityQueue) Len() int {
	return len(pq)
}

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].expireAt.Before(pq[j].expireAt)
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *priorityQueue) Push(x interface{}) {
	item := x.(*item)
	q := *pq
	q = append(q, item)
	*pq = q
}

func (pq *priorityQueue) Pop() interface{} {
	a := *pq
	n := len(a)
	item := a[n - 1]
	*pq = a[0 : n - 1]
	return item
}

type Cache struct {
	items map[string]*item
	pq    *priorityQueue
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

// Evict expired items from the cache
func (c*Cache) evict() {

	now := time.Now()

	for c.pq.Len() != 0 {

		//TODO: peek

		item := heap.Pop(c.pq).(*item)
		if (item.expireAt.Before(now)) {
			delete(c.items, item.key)
		} else {
			heap.Push(c.pq, item)
			break
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

// Get the number of items in the cache
func (c *Cache) Count() int {
	c.mutex.Lock()
	c.evict()
	count := len(c.items)
	c.mutex.Unlock()
	return count
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

	c.evict()

	expireAt := time.Now().Add(ttl)

	item := &item{
		key: key,
		value: value,
		ttl: ttl,
		expireAt: expireAt,
	}

	c.items[key] = item

	heap.Push(c.pq, item)
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

// Delete the value of the key
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

// Return all keys in the cache
func (c*Cache) Keys() ([]string) {

	c.mutex.Lock()
	c.evict()
	keys := make([]string, len(c.items))
	i := 0
	for k := range c.items {
		keys[i] = k
		i++
	}

	c.mutex.Unlock()
	return keys
}

// Schedule execution of the given function within a specified delay
func schedule(what func(), delay time.Duration) chan bool {
	stop := make(chan bool)

	go func() {
		for {
			what()
			select {
			case <-time.After(delay):
			case <-stop:
				return
			}
		}
	}()

	return stop
}

// Create a new cache
func NewCache() *Cache {

	pq := priorityQueue{}
	heap.Init(&pq)

	cache := &Cache{
		items:make(map[string]*item),
		pq: &pq,
	}

	// Schedule eviction execution on interval
	stop := schedule(func() {
		cache.mutex.Lock()
		cache.evict()
		cache.mutex.Unlock()
	}, DefaultEvictionInterval)

	// Stop scheduling on finalization
	runtime.SetFinalizer(cache, func(cache *Cache) {
		stop <- true
	})

	return cache
}