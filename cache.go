package gcache

import (
	"time"
	"errors"
	"sync"
	"container/heap"
	"container/list"
	"runtime"
	"fmt"
)

const DefaultEvictionInterval time.Duration = 1 * time.Second
const MaxDuration time.Duration = 1<<63 - 1

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

func (pq *priorityQueue) Peek() interface{} {
	a := *pq
	return a[0]
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

		item := c.pq.Peek().(*item)

		if (item.expireAt.Before(now)) {
			heap.Pop(c.pq)
			delete(c.items, item.key)
		} else {
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

// ====== LIST ======

// Left push value into the list
func (c*Cache) LPush(listKey string, value interface{}) error {

	return c.listPush(listKey, value, func(l *list.List)  {
		l.PushBack(value)
	})
}

func (c*Cache) listPush(listKey string, value interface{}, push func(l *list.List)) error {
	c.mutex.Lock()

	if item, ok := c.getItem(listKey); ok {
		l, ok := item.value.(*list.List)
		if (!ok){
			c.mutex.Unlock()
			return fmt.Errorf("Given %s is not a list key", listKey)
		}
		push(l)
	} else {
		//Create new list
		l := list.New()
		l = l.Init()
		l.PushFront(value)
		c.set(listKey, l, MaxDuration)
	}

	c.mutex.Unlock()

	return nil
}

func (c*Cache) listPop(listKey string, pop func(l *list.List) interface{}) (interface{}, error) {
	c.mutex.RLock()

	if item, ok := c.getItem(listKey); ok {
		l, ok := item.value.(*list.List)
		if (!ok){
			c.mutex.RUnlock()
			return nil, fmt.Errorf("Given %s is not a list key", listKey)
		}

		element := pop(l)
		c.mutex.RUnlock()

		return element, nil

	} else {
		c.mutex.RUnlock()
		return nil, ErrKeyNotFound
	}
}

// Right push value into the list
func (c*Cache) RPush(listKey string, value interface{}) error {

	return c.listPush(listKey, value, func(l *list.List)  {
		l.PushFront(value)
	})
}

// Left pop value from the list
func (c*Cache) LPop(listKey string) (interface{}, error) {

	return c.listPop(listKey, func(l *list.List) interface{} {
		elem := l.Back()
		l.Remove(elem)
		return elem.Value
	})
}

// Right pop vaues from the list
func (c*Cache) RPop(listKey string)  (interface{}, error)  {

	return c.listPop(listKey, func(l *list.List) interface{} {
		elem := l.Front()
		l.Remove(elem)
		return elem.Value
	})
}

// Returns a range of values from the list
func (c*Cache) LRange(listKey string, from int, to int)  ([]interface{}, error)  {

	c.mutex.RLock()

	if item, ok := c.getItem(listKey); ok {
		l, ok := item.value.(*list.List)
		if (!ok){
			c.mutex.RUnlock()
			return nil, fmt.Errorf("Given %s is not a list key", listKey)
		}

		index := 0
		result := make([]interface{}, 0)

		//TODO: validate from and to

		for e := l.Front(); e != nil; e = e.Next() {
			if (index >= from && index <= to) {
				result = append(result, e.Value)
			}

			index ++
		}

		c.mutex.RUnlock()
		return result, nil

	} else {
		c.mutex.RUnlock()
		return nil, ErrKeyNotFound
	}

}

// ====== HASH ======

func (c*Cache) HSet(hashKey string, key string, value interface{}) error {
	c.mutex.Lock()

	if item, ok := c.getItem(hashKey); ok {
		hash, ok := item.value.(map[string]interface{})
		if (!ok){
			c.mutex.Unlock()
			return fmt.Errorf("Given %s is not a hash key", hashKey)
		}
		hash[key] = value
	} else {
		//Create new hash
		hash := make(map[string]interface{})
		hash[key] = value
		c.set(hashKey, hash, MaxDuration)
	}

	c.mutex.Unlock()

	return nil
}

func (c*Cache) HGet(hashKey string, key string) (interface{}, error) {

	c.mutex.RLock()

	if item, ok := c.getItem(hashKey); ok {
		hash, ok := item.value.(map[string]interface{})
		if (!ok){
			c.mutex.RUnlock()
			return nil, fmt.Errorf("Given '%s' is not a hash key", hashKey)
		}

		value, ok := hash[key]
		if (!ok) {
			c.mutex.RUnlock()
			return nil, fmt.Errorf("Hash '%s' does not contain '%s' key", hashKey, key)
		}

		c.mutex.RUnlock()

		return value, nil

	} else {
		c.mutex.RUnlock()
		return nil, ErrKeyNotFound
	}

	c.mutex.RUnlock()

	return nil, nil
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