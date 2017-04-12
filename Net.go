package main

import (
	"fmt"
	"net"
	"os"

	"gopkg.in/vmihailenco/msgpack.v2"
)

//Channel that handles input
var InputPacketChan = make(chan InputPacket, 100)

//ListenAddress The listen address for the server
var ListenAddress = "127.0.0.1:7778"

//BufSize The buffer size for receiving data
var BufSize = 2048

var (
	conn *net.UDPConn
)

//obtains the packet id from the packet
type InputPacket struct {
	addr *net.UDPAddr
	Id   int
	X    int
	Y    int
}

type StatePacket struct {
	Entities Entity
}

type TestPacket struct {
	TestFloat float32
}

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
	println("Listening for players...")

	buf := make([]byte, BufSize)

	for {
		n, addr, err := conn.ReadFromUDP(buf)
		CheckError(err)

		var inputPacket InputPacket
		err = msgpack.Unmarshal(buf[:n], &inputPacket)
		inputPacket.addr = addr
		InputPacketChan <- inputPacket
	}
	println("break")
}

//Sends data to players
func Send() {
	for _, p := range players {
		// var i = 0
		// entities := make([]Entity, len(players))
		// statePacket := StatePacket{}
		// statePacket.Entities = entities
		// for _, p2 := range players {
		// 	entities[i] = *p2.e
		// 	i++
		// }
		// b, err := msgpack.Marshal(statePacket)
		// CheckError(err)

		t := TestPacket{TestFloat: 1.4}
		b, err := msgpack.Marshal(t)
		CheckError(err)

		sendMessage(b, p.addr)
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
