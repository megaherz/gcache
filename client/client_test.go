package client

import (
	"log"
	"strconv"
	"testing"
)

// NOTE, before running the test execute the server.sh script in the run folder.
// The scripts runs two cache servers on ports 8080 and 8081
// The 8081 server is run with authentication psw=123

const connectionString string = "http://localhost:8080"
const connectionStringAuth string = "http://localhost:8081"

func TestClient_SetGetDel(t *testing.T) {
	setDelGet(NewClient(connectionString), t)
}

func TestClient_SetGetDel_WithAuth(t *testing.T) {
	setDelGet(NewClientWithAuth(connectionStringAuth, "123"), t)
}

func setDelGet(client *Client, t *testing.T) {

	const key = "key"
	const value = "value"

	err := client.Set(key, value, 5)

	if err != nil {
		t.Error("Failed to set the key", err)
	}

	returnedValue, err := client.Get(key)

	if err != nil {
		t.Error("Failed to get the key", err)
	}

	if returnedValue != value {
		t.Error("Set value", value, "is not equal to returned value", returnedValue)
	}

	err = client.Del(key)

	if err != nil {
		t.Error("Failed to delete the key", err)
	}
}

func TestClient_Keys(t *testing.T) {

	client := NewClient(connectionString)
	keys, err := client.Keys()

	const key1 = "key1"
	const key2 = "key2"

	if err != nil {
		t.Error("Failed to get keys", err)
	}

	if len(keys) != 0 {
		t.Error("There should be no keys")
	}

	// Insert key1
	err = client.Set(key1, "value", 5)

	if err != nil {
		t.Error("Failed to set the key1", err)
	}

	// Insert key2
	err = client.Set(key2, "value", 5)

	if err != nil {
		t.Error("Failed to set the key2", err)
	}

	keys, err = client.Keys()

	if err != nil {
		t.Error("Failed to get keys", err)
	}

	if len(keys) != 2 {
		t.Error("There should be only one keys")
	}

	if !contains(keys, key1) || !contains(keys, key2) {
		t.Errorf("Keys contains unexpected key. Keys=%s", keys)
	}

}

func TestClient_HSet_HGET(t *testing.T) {
	const hashKey = "hashKey"
	const key = "key"
	const value = "value"

	client := NewClient(connectionString)
	err := client.HSet(hashKey, key, value)

	if err != nil {
		t.Errorf("Failed to hset '%s' with key '%s' and value '%s'. Err = %s", hashKey, key, value, err)
	}

	returnedValue, err := client.HGet(hashKey, key)

	if err != nil {
		t.Error("Failed to get the key", err)
	}

	if returnedValue != value {
		t.Error("Value", value, "is not equal to returned value", returnedValue)
	}
}

func TestClient_LRange(t *testing.T) {
	const listKey = "rangelistKey"

	client := NewClient(connectionString)

	// LPUSH 10 items
	for i := 0; i < 10; i++ {
		err := client.LPush(listKey, strconv.Itoa(i))
		if err != nil {
			t.Errorf("Failed to lpush. ListKey = '%s'. Error = %s", listKey, err)
		}
	}

	values, err := client.LRange(listKey, 2, 4)

	if err != nil {
		t.Fatalf("LRange failed. Err = %s", err)
	}

	log.Println("LRANGE", values)

	for _, value := range values {
		if value != "2" && value != "3" && value != "4" {
			t.Errorf("Values contain unexpected value '%s", value)
		}
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
