package main

import (
	"fmt"
	"net"
)

func sendMessage(msg []byte, conn *net.UDPConn, addr *net.UDPAddr) {
	_, err := conn.WriteToUDP(msg, addr)
	if err != nil {
		fmt.Printf("Couldn't send response %v", err)
	}
}
