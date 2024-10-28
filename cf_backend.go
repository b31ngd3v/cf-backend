package main

import (
	"os"
	"os/signal"

	"github.com/b31ngd3v/cf-backend/internal/server"
)

func main() {
	go server.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop
}
