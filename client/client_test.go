package client

import (
	"testing"
)

const connectionString string = "http://localhost:8080"

func TestClient_SetGetDel(t *testing.T) {

	const key = "key"
	const value  = "value"

	client := NewClient(connectionString)
	err := client.Set(key, value, 5)

	if (err != nil) {
		t.Error("Failed to set the key", err)
	}

	returnedValue, err := client.Get(key)

	if (err != nil) {
		t.Error("Failed to get the key", err)
	}

	if (returnedValue != value) {
		t.Error("Set value", value, "is not equal to returned value", returnedValue)
	}

	err = client.Del(key)

	if (err != nil) {
		t.Error("Failed to delete the key", err)
	}
}

func TestClient_Keys(t *testing.T) {

	client := NewClient(connectionString)
	keys, err := client.Keys()

	const key1 = "key1"
	const key2 = "key2"

	if (err != nil) {
		t.Error("Failed to get keys", err)
	}

	if (len(keys) != 0){
		t.Error("There should be no keys")
	}

	// Insert key1
	err = client.Set(key1, "value", 5)

	if (err != nil) {
		t.Error("Failed to set the key1", err)
	}

	// Insert key2
	err = client.Set(key2, "value", 5)

	if (err != nil) {
		t.Error("Failed to set the key2", err)
	}

	keys, err = client.Keys()

	if (err != nil) {
		t.Error("Failed to get keys", err)
	}

	if (len(keys) != 2){
		t.Error("There should be only one keys")
	}

	if !contains(keys, key1) || !contains(keys, key2) {
		t.Errorf("Keys contains unexpected key. Keys=%s", keys)
	}

}

func TestClient_HSet_HGET(t *testing.T) {
	const hashKey = "hashKey"
	const key = "key"
	const value  = "value"

	client := NewClient(connectionString)
	err := client.HSet(hashKey, key, value)

	if (err != nil) {
		t.Errorf("Failed to hset '%s' with key '%s' and value '%s'. Err = %s", hashKey, key, value, err)
	}

	returnedValue, err := client.HGet(hashKey, key)

	if (err != nil) {
		t.Error("Failed to get the key", err)
	}

	if (returnedValue != value) {
		t.Error("Value", value, "is not equal to returned value", returnedValue)
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

