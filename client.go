package main

import (
	"net"
	"time"
	"fmt"
	"io"
	"encoding/binary"
)


type Session struct {
	ip           net.IP
	conn         net.Conn

	user_id      int32
	app_id       string

	close        chan struct{} //会话关闭信号
	in           chan []byte
	out          chan []byte
	ntf          chan []byte

	packet_count int32         //对包进行计数
	connect_time time.Time
}

func (sess *Session) init(conn net.Conn) bool {
	host, port, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		fmt.Println("cannot get remote address:", err)
		return false
	}
	fmt.Printf("new connection from:%v port:%v", host, port)

	sess.ip = net.ParseIP(host)
	sess.conn = conn

	sess.close = make(chan struct{})
	sess.in = make(chan []byte)
	sess.out = make(chan []byte)
	sess.ntf = make(chan []byte)

	sess.packet_count = 0
	sess.connect_time = time.Now()

	return true
}

func (sess *Session) clean() {
	close(sess.close)
	close(sess.in)
	close(sess.out)
	close(sess.ntf)
	// close the connection
	sess.conn.Close()
}



// the go routine is used for reading incoming PACKETS
// each packet is defined as :
// | 2B size |     DATA       |
func handleClient(conn net.Conn) {
	// init session
	var sess Session
	ok := sess.init(conn)
	if !ok {
		fmt.Println("Error: init seesion error.")
		return
	}
	defer sess.clean()

	//go out buf
	go sender(&sess)
	//go agent
	go agent(&sess)

	//read loop
	header := make([]byte, 2)
	for {
		// solve dead link problem:
		// physical disconnection without any communication between client and server
		// will cause the read to block FOREVER, so a timeout is a rescue.
		conn.SetReadDeadline(time.Now().Add(TCP_READ_DEADLINE * time.Second))

		// read head, 如果客户端关闭socket，在这里会return，然后执行 defer sess.clean
		n, err := io.ReadFull(conn, header)
		if err != nil {
			fmt.Printf("Error: read header failed, ip:%v reason:%v size:%v\n", sess.ip, err, n)
			return
		}
		size := binary.BigEndian.Uint16(header)

		// read data
		payload := make([]byte, size)  //TODO 优化 使用固定分配好的buf，不用每次都从新分配
		n, err = io.ReadFull(conn, payload)
		if err != nil {
			fmt.Printf("Error: read payload failed, ip:%v reason:%v size:%v\n", sess.ip, err, n)
			return
		}

		//TODO delete
		fmt.Printf("recv size:%v, data:%v\n", size, payload)

		select {
		case sess.in <- payload:
		case <- die:
			// close sess.close in defer
			fmt.Printf("connection closed by logic, ip:%v\n", sess.ip)
			return
		}
	}
}


