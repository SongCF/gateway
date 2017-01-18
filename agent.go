package main

import (
	"time"
	"fmt"
)


func agent(sess *Session) {
	//TODO delete
	defer func() {
		fmt.Println("------ agent end.")
	}()

	defer close(sess.die)  //sess.in关闭和收到global_die都会结束agent

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
			if !ok {  //session 中如果连接断开，或收到global_die会关闭sess.in
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
		case <- global_die:
			return
		}
	}
}


func timer_work() {
	fmt.Println("on timer")
}