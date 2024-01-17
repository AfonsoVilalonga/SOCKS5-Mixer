package main

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/knadh/chunkedreader"
)

var channelMap = make(map[uint64]net.Conn)

func send(conn_server net.Conn, conn_client net.Conn, streamID uint64) {
	for {
		buffer := make([]byte, 1024)
		n, err := conn_client.Read(buffer)
		if err != nil {
			return
		}

		id_bytes := make([]byte, 8)
		binary.BigEndian.PutUint64(id_bytes, streamID)

		len := make([]byte, 8)
		binary.BigEndian.PutUint64(len, uint64(n))

		aux := append(id_bytes, len...)
		buffer = append(aux, buffer...)

		conn_server.Write(buffer)

	}
}

func readFromServer(conn_server net.Conn) {
	ch := chunkedreader.New(conn_server, 1040)
	for ch.Read() {
		buffer := ch.Bytes()

		streamID := binary.BigEndian.Uint64(buffer[:8])
		len := binary.BigEndian.Uint64(buffer[8:16])

		conn_sock := channelMap[streamID]
		if len > 0 {
			conn_sock.Write(buffer[16 : 16+len])
		}
	}
}

func main() {
	var streamID uint64 = 0
	port := 8080

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return
	}
	defer listener.Close()

	serverAddr := "127.0.0.1"
	serverPort := 7777

	conn_server, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverAddr, serverPort))
	if err != nil {
		return
	}

	go readFromServer(conn_server)

	for {
		conn_client, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		channelMap[streamID] = conn_client

		go send(conn_server, conn_client, streamID)
		streamID++
	}
}
