
package main

import "gcache/server"

func main() {
	server := server.NewServer()
	server.Run(":8080")
}
