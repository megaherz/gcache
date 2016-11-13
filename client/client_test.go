package client

import (
	"testing"
)

const connectionString string = "http://localhost:8080"

func TestClient_SetGet(t *testing.T) {

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
}
