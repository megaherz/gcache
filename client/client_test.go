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

	const key = "key"

	if (err != nil) {
		t.Error("Failed to get keys", err)
	}

	if (len(keys) != 0){
		t.Error("There should be no keys")
	}

	err = client.Set(key, "value", 5)

	if (err != nil) {
		t.Error("Failed to set the key", err)
	}

	keys, err = client.Keys()

	if (err != nil) {
		t.Error("Failed to get keys", err)
	}

	if (len(keys) != 1){
		t.Error("There should be only one keys")
	}

	if (keys[0] !=  key) {
		t.Errorf("Unexpected key %s", keys[0])
	}

}
