package client

import (
	"log"
	"strconv"
	"testing"
)

// NOTE, before running the test execute the server.sh script in the ./cmd/gcache folder.
// The scripts runs two cache servers on ports 8080 and 8081
// The 8081 server is run with authentication psw=123

const (
	connectionString string = "http://localhost:8080"
	connectionStringAuth string = "http://localhost:8081"
	psw string = "123"
)

func TestClient_SetGetDel(t *testing.T) {

	const key = "key"
	const value = "value"

	conns := Connections{
		{connectionString, ""},
	}

	client := NewClient(conns)
	test_SetGetDel(client, key, value, t)
}


func TestClient_SetGetDel_WithAuth(t *testing.T) {

	const key = "key"
	const value = "value"

	conns := Connections{
		{connectionStringAuth, psw},
	}

	client := NewClient(conns)
	test_SetGetDel(client, key, value, t)
}

func TestClient_Keys_Sharded(t *testing.T) {

	const n = 10

	// Two shards
	conns := Connections{
		{connectionString, ""},
		{connectionStringAuth, psw},
	}

	client := NewClient(conns)

	for i := 0; i < n; i++ {
		err := client.Set(strconv.Itoa(i), "value", 5)

		if err != nil {
			t.Error("Failed to set the key", err)
		}
	}

	keys, err := client.Keys()

	if err != nil {
		t.Error("Failed to get keys", err)
	}

	if len(keys) != n {
		log.Println("Keys", keys)
		t.Errorf("There should be %d keys, but there are %d keys", n, len(keys))
	}

	// Make sure keys are distributed between shards
	client1 := NewClient(Connections{
		{connectionString, ""},
	})

	client2 := NewClient(Connections{
		{connectionStringAuth, psw},
	})


	keys1, _ := client1.Keys()
	keys2, _ := client2.Keys()


	if len(keys1) == 0 || len(keys2) == 0 {
		t.Error("Keys are not distributed")
	}

	//Tear down
	for i := 0; i < n; i++ {
		client.Del(strconv.Itoa(i))
	}
}


func test_SetGetDel(client *Client, key string, value string, t *testing.T) {

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

	conns := Connections{
		{connectionString, ""},
	}

	client := NewClient(conns)

	keys, err := client.Keys()

	const key1 = "key1"
	const key2 = "key2"

	if err != nil {
		t.Error("Failed to get keys", err)
	}

	if len(keys) != 0 {
		t.Errorf("There should be no keys, but there are %d keys", len(keys))
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
		log.Println("Keys", keys)
		t.Errorf("There should be '%d' keys, but there are %d keys", 2, len(keys))
	}

	if !contains(keys, key1) || !contains(keys, key2) {
		t.Errorf("Keys contains unexpected key. Keys=%s", keys)
	}

	// Tear down
	client.Del(key1)
	client.Del(key2)

}


func TestClient_Update(t *testing.T) {

	const key = "key"
	const updatedValue = "updated"

	conns := Connections{
		{connectionString, ""},
	}

	client := NewClient(conns)

	err := client.Update(key, "value")

	if err == nil {
		t.Error("Expected: key not found")
	}

	// Insert key
	err = client.Set(key, "value", 5)

	if err != nil {
		t.Error("Failed to set the key", err)
	}

	err = client.Update(key, updatedValue)

	if (err != nil) {
		t.Error("Failed to update. Err = ", err)
	}

	value, err := client.Get(key)

	if value != updatedValue {
		t.Errorf("Update value '%s' does not equal to returned value '%s'", updatedValue, value)
	}

	// Tear down
	client.Del(key)

}

func TestClient_UpdateWithTtl(t *testing.T) {

	const key = "key"
	const updatedValue = "updated"

	conns := Connections{
		{connectionString, ""},
	}

	client := NewClient(conns)

	err := client.UpdateWithTtl(key, "value", 5)

	if err == nil {
		t.Error("Expected: key not found")
	}

	// Insert key
	err = client.Set(key, "value", 5)

	if err != nil {
		t.Error("Failed to set the key", err)
	}

	// Update
	err = client.UpdateWithTtl(key, updatedValue, 25)

	// Assertions
	if (err != nil) {
		t.Error("Failed to update. Err = ", err)
	}

	value, err := client.Get(key)

	if value != updatedValue {
		t.Errorf("Update value '%s' does not equal to returned value '%s'", updatedValue, value)
	}

	// Tear down
	client.Del(key)

}

func TestClient_HSet_HGET(t *testing.T) {

	const key = "key"
	const hashKey = "hashKey"
	const value = "value"

	conns := Connections{
		{connectionString, ""},
	}

	client := NewClient(conns)

	err := client.HSet(key, hashKey, value)

	if err != nil {
		t.Errorf("Failed to hset '%s' with hash key '%s' and value '%s'. Err = %s", key, hashKey, value, err)
	}

	returnedValue, err := client.HGet(key, hashKey)

	if err != nil {
		t.Error("Failed to get the key", err)
	}

	if returnedValue != value {
		t.Error("Value", value, "is not equal to returned value", returnedValue)
	}

	// Tear down
	client.Del(key)
}

func TestClient_LRange_LPUSH_LPOP(t *testing.T) {
	const listKey = "rangelistKey"

	conns := Connections{
		{connectionString, ""},
	}

	client := NewClient(conns)

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

	// Tear down
	for i := 0; i < 10; i++ {
		_, err := client.LPop(listKey)
		if err != nil {
			t.Errorf("Failed to lpop. ListKey = '%s'. Error = %s", listKey, err)
		}
	}

	// Tear down
	client.Del(listKey)
}

func TestClient_RPush_RPop(t *testing.T) {

	const key = "rpushpopkey"
	const value = "value"

	conns := Connections{
		{connectionString, ""},
	}

	client := NewClient(conns)

	// RPush
	err := client.RPush(key, value)
	if err != nil {
		t.Errorf("Failed to RPush. ListKey = '%s'. Error = %s", key, err)
	}

	// RPop
	returnedValue, err := client.RPop(key)
	if err != nil {
		t.Errorf("Failed to RPop. ListKey = '%s'. Error = %s", key, err)
	}

	if returnedValue != value {
		t.Error("Value", value, "is not equal to returned value", returnedValue)
	}

	// Tear down
	client.Del(key)

}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
