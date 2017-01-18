package main

import (
	"net"
	"encoding/binary"
	"fmt"
)


func sender(sess *Session) {
	//TODO delete
	defer func() {
		fmt.Println("------ sender end.")
	}()

	defer global_wg.Done()
	for {
		select {
		case data, ok := <- sess.out:
			if !ok {
				return
			}
			raw_send(sess.conn, data)
		case <- sess.die:
			close(sess.out)
			close(sess.ntf)
			sess.conn.Close()
			return
		}
	}
}


func send(sess *Session, data []byte) {
	// in case of empty packet
	if data == nil {
		return
	}

	//TODO delete
	fmt.Printf("put in out buf:%v\n", data)
	sess.out <- data
}


func raw_send(conn net.Conn, data []byte) bool {
	//TODO delete
	fmt.Println("to raw_send...")

	size := len(data)
	cache := make([]byte, size + 2) //TODO cache 优化，不用每次都创建
	binary.BigEndian.PutUint16(cache, uint16(size))
	copy(cache[2:], data)

	n, err := conn.Write(cache[:size+2])
	if err != nil {
		fmt.Printf("Error send reply data, bytes: %v reason: %v", n, err)
		return false
	}
	//TODO delete
	fmt.Println("raw_send ok!")
	return true
}
