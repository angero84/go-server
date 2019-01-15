package main

import (
	"fmt"
	"log"
	"net"
	"time"
	"kprotocol"
)

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:8989")
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	echoProtocol := &kprotocol.KProtocolEcho{}

	// ping <--> pong
	for i := 0; i < 3; i++ {
		// write
		conn.Write(kprotocol.NewKPacketEcho([]byte("hello"), false).Serialize())

		// read
		p, err := echoProtocol.ReadKPacket(conn)
		if err == nil {
			echoPacket := p.(*kprotocol.KPacketEcho)
			fmt.Printf("Server reply:[%v] [%v]\n", echoPacket.Length(), string(echoPacket.Body()))
		}

		time.Sleep(2 * time.Second)
	}

	conn.Close()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}