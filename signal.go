package main

import (
	"os"
	"syscall"
	"os/signal"
	"sync"
)

var (
	global_wg sync.WaitGroup
	// server close signal
	global_die = make(chan struct{})
	sigChan = make(chan os.Signal, 1)
)

func signalHandler() {
	signal.Notify(sigChan, syscall.SIGTERM)
	for {
		msg := <- sigChan
		switch msg {
		case syscall.SIGTERM:
			println("sigterm received")
			close(global_die)
			println("waiting for agents close, please wait...")
			// wait close agent
			global_wg.Wait()
			println("shutdown.")
			os.Exit(0)
		default:
		}
	}
}