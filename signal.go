package main

import (
	"os"
	"syscall"
	"os/signal"
	"sync"
)

var (
	wg sync.WaitGroup
	// server close signal
	die = make(chan struct{})
)

func signalHandler() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM)
	for {
		msg := <- sigChan
		switch msg {
		case syscall.SIGTERM:
			println("sigterm received")
			close(die)
			println("waiting for agents close, please wait...")
			// wait: 1.client(read sock), 2.agent(in), 3.out -> close conn
			wg.Wait()
			println("shutdown.")
			os.Exit(0)
		default:
		}
	}
}