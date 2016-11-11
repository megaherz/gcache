package gcache

import "testing"

func TestKeyValue(t *testing.T) {

	value := 24
	key := "key1"

	cache := NewCache()
 	cache.Set(key, value, 5)
	returnedValue, exist := cache.Get(key)

	if (!exist) {
		t.Error("Value with key", key, "does not exist" )
	}

	if (value != returnedValue) {
		t.Error("Expected value", value, "is not equal to", returnedValue)
	}
}

func TestList(t *testing.T) {

}

func TestDict(t *testing.T) {

}