package gcache

import (
	"strconv"
	"testing"
	"time"
)

func TestCache_Keys(t *testing.T) {
	cache := NewCache()

	cache.Set("key1", "value", time.Second)
	cache.Set("key2", "value", time.Second)
	cache.Set("key3", "value", time.Second)

	keys := cache.Keys()

	if len(keys) != 3 {
		t.Error("Expected is", 3, "but actual is", len(keys))
	}
}

func TestCache_Get(t *testing.T) {

	const value = 24
	const key = "key1"

	cache := NewCache()
	cache.Set(key, value, 5*time.Second)
	returnedValue, err := cache.Get(key)

	if err == ErrKeyNotFound {
		t.Error("Value with key", key, "does not exist")
	}

	if value != returnedValue {
		t.Error("Expected value", value, "is not equal to", returnedValue)
	}
}

func TestCache_Del(t *testing.T) {
	const key = "key1"

	cache := NewCache()
	cache.Set(key, "value", time.Second)
	err := cache.Del(key)

	if len(cache.Keys()) != 0 || err != nil {
		t.Error("The key has not been removed")
	}

	err = cache.Del(key)

	if err != ErrKeyNotFound {
		t.Error("The key still exists in Cache")
	}

}

func TestCache_Update(t *testing.T) {

	const key = "key1"
	const expectedValue = "expected"

	cache := NewCache()
	cache.Set(key, "value", 10*time.Microsecond)

	err := cache.Update(key, expectedValue)

	if err != nil {
		t.Error("Failed to update")
	}

	returnedValue, _ := cache.Get(key)

	if returnedValue != expectedValue {
		t.Error("Failed to update the value. The returned value is", returnedValue)
	}
}

func TestCache_UpdateWithTll(t *testing.T) {

	const key = "key1"
	const expectedTtl = time.Second * 5
	const expectedValue = "expected"

	cache := NewCache()
	cache.Set(key, "value", 10*time.Microsecond)

	err := cache.UpdateWithTll(key, expectedValue, expectedTtl)

	if err != nil {
		t.Error("Failed to update")
	}

	returnedValue, _ := cache.Get(key)

	if returnedValue != expectedValue {
		t.Error("Failed to update the value. The returned value is", returnedValue)
	}

	returnedTtl, _ := cache.Ttl(key)

	if returnedTtl != expectedTtl {
		t.Error("Failed to update the tll. The returned tll is", returnedTtl)
	}

}

func TestCache_Expiration(t *testing.T) {

	const key = "key1"

	cache := NewCache()
	cache.Set(key, "value", 10*time.Microsecond)

	time.Sleep(time.Second)

	_, err := cache.Get(key)

	if err != ErrKeyNotFound {
		t.Error("The key has not been expired")
	}

}

func TestCache_ParallelGetUpdate(t *testing.T) {

	const key = "key1"

	cache := NewCache()

	cache.Set(key, 0, time.Second*10)

	done := make(chan bool)

	// Get
	go func() {
		for i := 0; i < 100; i++ {
			value, _ := cache.Get("key")
			t.Log("Value", value)
		}

		done <- true
	}()

	// Set
	go func() {
		for i := 0; i < 100; i++ {
			cache.Update(key, i*100)
		}

		done <- true
	}()

	for i := 0; i < 2; i++ {
		<-done
	}
}

func TestCache_Eviction(t *testing.T) {

	cache := NewCache()

	for i := 0; i < 100; i++ {
		cache.Set("key"+string(i), "val", 1*time.Second)
	}

	count := cache.Count()

	t.Log("Count", count)

	time.Sleep(2 * time.Second)

	count = cache.Count()

	if count != 0 {
		t.Error("Count should be 0. There are ", count, "items which have not been evicted")
	}

}

func TestCache_LPush_LPop(t *testing.T) {

	const key = "list"
	const value = "value"

	cache := NewCache()

	err := cache.LPush(key, value)

	if err != nil {
		t.Fatal("Failed to push value into the list")
	}

	returnedValue, err := cache.LPop(key)

	if err != nil {
		t.Fatal("Failed to pop value from the list")
	}

	if returnedValue != value {
		t.Errorf("Retuned value '%s' does not equal to '%s", returnedValue, value)
	}
}

func TestCache_RPush_RPop(t *testing.T) {

	const key = "list"
	const value = "value"

	cache := NewCache()

	err := cache.RPush(key, value)

	if err != nil {
		t.Fatal("Failed to push value into the list")
	}

	returnedValue, err := cache.RPop(key)

	if err != nil {
		t.Fatal("Failed to pop value from the list")
	}

	if returnedValue != value {
		t.Errorf("Retuned value '%s' does not equal to '%s", returnedValue, value)
	}
}

func TestCache_LRange(t *testing.T) {
	cache := NewCache()

	const key = "list"
	const from  = 50

	for i := 0; i < 100; i++ {
		cache.LPush(key, i)
	}

	// Act 'to' less 'from'
	values, err := cache.LRange(key, 50, 40)

	if err != nil {
		t.Fatal("Failed to get a range of values from the list")
	}

	if len(values) != 0 {
		t.Errorf("Expected %d but actual %d", len(values), 0)
	}

	// Act equal to and from
	values, err = cache.LRange(key, from, from)

	if err != nil {
		t.Fatal("Failed to get exact on value from the list")
	}

	if len(values) != 1 {
		t.Errorf("Expected %d but actual %d", len(values), 0)
	}


	// Act
	values, err = cache.LRange(key, from, 70)

	if err != nil {
		t.Fatal("Failed to get a range of values from the list")
	}

	for i, value := range values {
		if (i + from) != value {
			t.Errorf("Expected %d but received %d", i, value)
		}
	}
}

func TestCache_HSet_HGet(t *testing.T) {

	const key = "list"
	const listKey = "list"
	const value = "value"

	cache := NewCache()

	err := cache.HSet(key, listKey, value)

	if err != nil {
		t.Fatal("Failed to push value into the hash")
	}

	returnedValue, err := cache.HGet(key, listKey)

	if err != nil {
		t.Fatal("Failed to get value from the hash")
	}

	if returnedValue != value {
		t.Errorf("Retuned value '%s' does not equal to '%s", returnedValue, value)
	}
}

func BenchmarkCache_SetGet(b *testing.B) {

	cache := NewCache()

	for n := 0; n < b.N; n++ {

		key := "key" + strconv.Itoa(n)

		cache.Set(key, "val", 1*time.Second)
		_, err := cache.Get(key)

		if err != nil {
			b.Error("Failed to GET", key)
		}
	}
}

func BenchmarkCache_SetGet_Parallel(b *testing.B) {

	cache := NewCache()

	b.RunParallel(func(pb *testing.PB) {

		for pb.Next() {

			key := strconv.Itoa(time.Now().Nanosecond())

			cache.Set(key, "val", 10*time.Second)
			_, err := cache.Get(key)

			if err != nil {
				b.Error("Failed to GET", key)
			}
		}

	})
}
