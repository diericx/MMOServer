package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"gopkg.in/vmihailenco/msgpack.v2"
)

//ListenAddress The listen address for the server
var ListenAddress = ":7778"

//BufSize The buffer size for receiving data
var BufSize = 2048

//FrameWaitTime time between frames
var FrameWaitTime float64 = 33

var (
	conn *net.UDPConn
)

//InitConnection Initializes the server connection
func InitConnection() {
	/* Lets prepare an address at any address at the listen address*/
	ServerAddr, err := net.ResolveUDPAddr("udp", ListenAddress)
	CheckError(err)

	/* Now listen at selected port */
	conn, err = net.ListenUDP("udp", ServerAddr)
	CheckError(err)
}

//Listen Listens for incomming packets, also closes the connection
func Listen() {
	defer conn.Close()

	buf := make([]byte, BufSize)

	for {
		n, addr, err := conn.ReadFromUDP(buf)
		CheckError(err)

		var data Entity
		err = msgpack.Unmarshal(buf[:n], &data)

		p := GetPlayer(addr)
		if p == nil {
			p = NewPlayer(addr)
		}

		p.e.X = data.X
		p.e.Y = data.Y
		p.e.Z = data.Z
	}
}

//Send Sends data to players
func Send() {
	for {
		w := ForLoopWaiter{start: time.Now()}

		for _, p := range players {
			for _, p2 := range players {
				if p == p2 {
					continue
				}
				b, err := msgpack.Marshal(p2.e)
				CheckError(err)

				sendMessage(b, p.addr)
			}
		}

		w.waitForTime(FrameWaitTime)
	}
}

func sendMessage(msg []byte, addr *net.UDPAddr) {
	_, err := conn.WriteToUDP(msg, addr)
	if err != nil {
		fmt.Printf("Couldn't send response; %v\n", err)
	}
}

//CheckError Checks errors and prints the message if there is one
func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}
