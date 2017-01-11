package main

import (
	"net"
	"encoding/binary"
	"fmt"
)


func sender(sess *Session) {
	for {
		select {
		case data := <- sess.out:
			raw_send(sess.conn, data)
		case <- die:
			close(sess.die)
			return
		}
	}
}


func send(sess *Session, data []byte) {
	// in case of empty packet
	if data == nil {
		return
	}
	sess.out <- data
}


func raw_send(conn net.Conn, data []byte) bool {
	size := len(data)
	cache := make([]byte, size + 2) //TODO cache 优化，不用每次都创建
	binary.BigEndian.PutUint16(cache, uint16(size))
	copy(cache[2:], data)

	n, err := conn.Write(cache[:size+2])
	if err != nil {
		fmt.Printf("Error send reply data, bytes: %v reason: %v", n, err)
		return false
	}
	return true
}
