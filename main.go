package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"turbo_snail/broker"
	"turbo_snail/config"
	"turbo_snail/http"
	"turbo_snail/restore"
	"turbo_snail/tcp"
)

const (
	TCP_PORT  = "7777"
	HTTP_PORT = "7000"
)

func main() {
	config.Init()
	// recieve msgs from tcp connection
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	turboSnailBroker := broker.Get()
	if err := restore.Init(turboSnailBroker); err != nil {
		log.Fatalf("Error occured while trying to restore. \n %s", err.Error())
	}
	var wg sync.WaitGroup
	// start http server
	wg.Add(1)
	go http.StartServer(config.HTTP_PORT, &wg)
	// on each tcp message , look for the track
	wg.Add(1)
	go tcp.Listen(config.TCP_PORT, turboSnailBroker, &wg)

	<-sigChan

}
