
package main

import (
	"gcache/client"
	"gcache"
	"fmt"
	"time"
)

func main() {
	client := client.NewClient("localhost:8080")
	client.Set("key", "value", 25)
	fmt.Println(client.Get("key"))
}
