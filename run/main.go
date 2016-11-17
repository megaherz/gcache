
package main

import (
	"gcache/server"
	"os/signal"
	"os"
	"syscall"
	"flag"
)

func main() {

	psw := flag.String("pws", "", "authentication password")
	flag.Parse()

	// Exit on Ctrl+C
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(1)
	}()

	server := server.NewServerWithAuth(*psw)
	server.SetUrlLogging(true)
	server.Run(":8080")
}
