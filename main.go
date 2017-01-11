package main

import (
	"os"
	"net"
	"fmt"
	"net/http"
)


func main() {
	//start pprof
	// http://localhost:6060/debug/pprof/
	go func() {
		err := http.ListenAndServe("0.0.0.0:6060", nil)
		fmt.Println("start pprof ", err)
	}()

	//handle signal
	go signalHandler()

	//start tcp server
	go tcpServer()

	//init service discover

	//wait forever
	select {}
}


func tcpServer() {
	addr, err := net.ResolveTCPAddr("tcp4", TCP_PORT)
	if err != nil {
		fmt.Println("Error: resolve tcp addr error, ", err)
		os.Exit(-1)
	}
	listener, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		fmt.Println("Error: listen tcp error, ", err)
		os.Exit(-1)
	}
	fmt.Println("listen on: ", listener.Addr())

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("Waring: accept tcp error, ", err)
			continue
		}

		// start a session
		go handleClient(conn)

		// check server close signal
		select {
		case <- die:
			listener.Close()
			return
		default:
		}
	}
}
