package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

var channelMap = make(map[uint64]net.Conn)

func send(conn_sock net.Conn, conn_client net.Conn, streamID uint64) {
	for {
		buffer := make([]byte, 1024)
		n, err := conn_sock.Read(buffer)
		if err != nil {
			return
		}
		id_bytes := make([]byte, 8)
		binary.BigEndian.PutUint64(id_bytes, streamID)

		len := make([]byte, 8)
		binary.BigEndian.PutUint64(len, uint64(n))

		aux := append(id_bytes, len...)
		buffer = append(aux, buffer...)

		conn_client.Write(buffer)
	}
}

func recv(conn_client net.Conn) {
	for {
		buffer := make([]byte, 1040)
		_, err := conn_client.Read(buffer)
		if err != nil {
			return
		}
		streamID := binary.BigEndian.Uint64(buffer[:8])
		len := binary.BigEndian.Uint64(buffer[8:16])

		conn_sock, prs := channelMap[streamID]

		if !prs {
			conn_sock, err = net.Dial("tcp", fmt.Sprintf("%s:%d", "localhost", 5555))
			if err != nil {
				return
			}
			channelMap[streamID] = conn_sock

			go send(conn_sock, conn_client, streamID)
		}

		if len > 0 {
			conn_sock.Write(buffer[16 : 16+len])
		}
	}
}

func main() {
	port := 7777

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return
	}

	conn_client, err := listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection:", err)
		return
	}

	recv(conn_client)
}
