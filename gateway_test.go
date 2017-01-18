package main


import (
	"testing"
	"net"
	"encoding/binary"
	"io"
	"strings"
	"fmt"
	"time"
)

const (
	TEST_SERVER_ADDR = "localhost:8888"
)


func TestTCP(t *testing.T) {
	conn, err := net.Dial("tcp", TEST_SERVER_ADDR)
	if err != nil {
		fmt.Printf("connect server error, %v\n", err)
		return
	}

	req_data := "test_req"
	req_size := len(req_data)
	req_buf := make([]byte, req_size + 2)
	binary.BigEndian.PutUint16(req_buf, uint16(req_size))
	copy(req_buf[2:], req_data)

	n, err := conn.Write(req_buf[:req_size+2])
	if err != nil {
		fmt.Printf("send data error, size=%d, %v", n, err)
		return
	}

	rsp_head := make([]byte, 2)
	n, err = io.ReadFull(conn, rsp_head)
	if err != nil {
		fmt.Printf("read rsp head error, %v", err)
		return
	}
	rsp_size := binary.BigEndian.Uint16(rsp_head)
	rsp_data := make([]byte, rsp_size)
	n, err = io.ReadFull(conn, rsp_data)
	if err != nil {
		fmt.Printf("read rsp data error, %v", err)
		return
	}

	// 如何拼接字符串 效率高：http://studygolang.com/articles/2507
	eq := strings.Compare(string(req_data[:]) + "_ack", string(rsp_data[:]))
	if eq != 0 {
		fmt.Printf("req data = %s, rsp data = %s", req_data[:], rsp_data[:])
		return
	}

	fmt.Println("ok!")

	//time.Sleep(time.Second * 100)
	conn.Close()
}