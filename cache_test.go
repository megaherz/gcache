package gcache

import (
	"testing"
	"time"
)

func TestCache_Keys(t *testing.T) {
	cache := NewCache()

	cache.Set("key1", "value", time.Second)
	cache.Set("key2", "value", time.Second)
	cache.Set("key3", "value", time.Second)

	keys := cache.Keys()

	if (len(keys) != 3) {
		t.Error("Expected is", 3, "but actual is", len(keys))
	}
}

func TestCache_Get(t *testing.T) {

	value := 24
	key := "key1"

	cache := NewCache()
 	cache.Set(key, value, 5 * time.Second)
	returnedValue, err := cache.Get(key)


	if (err == ErrKeyNotFound) {
		t.Error("Value with key", key, "does not exist" )
	}

	if (value != returnedValue) {
		t.Error("Expected value", value, "is not equal to", returnedValue)
	}
}

func TestCache_Del(t *testing.T) {
	key := "key1"

	cache := NewCache()
	cache.Set(key, "value", time.Second)
	err := cache.Del(key)

	if (len(cache.Keys()) != 0 || err != nil) {
		t.Error("The key has not been removed")
	}

	err = cache.Del(key)

	if (err != ErrKeyNotFound) {
		t.Error("The key still exists in Cache")
	}

}

func TestCache_Update(t *testing.T) {

	key := "key1"
	expectedValue := "expected"

	cache := NewCache()
	cache.Set(key, "value", 10 * time.Microsecond)

	err := cache.Update(key, expectedValue)

	if (err != nil) {
		t.Error("Failed to update")
	}

	returnedValue, _ := cache.Get(key)

	if returnedValue != expectedValue {
		t.Error("Failed to update the value. The returned value is", returnedValue)
	}
}

func TestCache_UpdateWithTll(t *testing.T) {
	key := "key1"
	expectedTtl := time.Second * 5
	expectedValue := "expected"

	cache := NewCache()
	cache.Set(key, "value", 10 * time.Microsecond)

	err := cache.UpdateWithTll(key, expectedValue, expectedTtl)

	if (err != nil) {
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

	key := "key1"

	cache := NewCache()
	cache.Set(key, "value", 10 * time.Microsecond)

	time.Sleep(time.Second)

	_, err :=  cache.Get(key)

	if (err != ErrKeyNotFound) {
		t.Error("The key has not been expired")
	}

}