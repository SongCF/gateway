package main


import (
	"testing"
	"net"
	"encoding/binary"
	"io"
	"strings"
)


func TestTCP(t *testing.T) {
	conn, err := net.Dial("tcp", TCP_PORT)
	if err != nil {
		t.Fatalf("connect server error, %v", err)
	}

	req_data := "test_req"
	req_size := len(req_data)
	req_buf := make([]byte, req_size + 2)
	binary.BigEndian.PutUint16(req_buf, uint16(req_size))
	copy(req_buf[2:], req_data)

	n, err := conn.Write(req_buf[:req_size+2])
	if err != nil {
		t.Fatalf("send data error, size=%d, %v", n, err)
	}

	rsp_head := make([]byte, 2)
	n, err = io.ReadFull(conn, rsp_head)
	if err != nil {
		t.Fatalf("read rsp head error, %v", err)
	}
	rsp_size := binary.BigEndian.Uint16(rsp_head)
	rsp_data := make([]byte, rsp_size)
	n, err = io.ReadFull(conn, rsp_data)
	if err != nil {
		t.Fatalf("read rsp data error, %v", err)
	}

	// 如何拼接字符串 效率高：http://studygolang.com/articles/2507
	eq := strings.Compare("ack" + string(req_data[:]), string(rsp_data[:]))
	if eq != 0 {
		t.Fatalf("req data = %v, rsp data = %v", req_data, rsp_data)
	}
}