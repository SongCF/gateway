package main

import (
	"time"
	"fmt"
)


func agent(sess *Session) {
	// wait group
	wg.Add(1)
	defer wg.Done()
	// minute timer
	min_timer := time.After(time.Minute)

	// >> the main message loop <<
	// handles 4 types of message:
	//  1. from client
	//  2. from game service
	//  3. timer
	//  4. server shutdown signal
	for {
		select {
		case msg, ok := <- sess.in:
			if !ok {
				return
			}
			sess.packet_count++
			//TODO delete
			fmt.Printf("msg:%v, pack_count:%v\n", msg, sess.packet_count)

			// TODO game rpc
			// ret
			ret := string(msg[:]) + "_ack"
			send(sess, []byte(ret))
		case msg, ok := <- sess.ntf:
			if !ok {
				return
			}
			// TODO
			send(sess, msg)
		case <- min_timer:
			timer_work()
			min_timer = time.After(time.Minute)
		case <- sess.die:
			return
		}
	}
}


func timer_work() {
	fmt.Println("on timer")
}