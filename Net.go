package main

import (
	"fmt"
	"net"
	"os"

	"gopkg.in/vmihailenco/msgpack.v2"
)

//Channel that handles input
var MovePacketChan = make(chan MovePacket, 100)

//ListenAddress The listen address for the server
var ListenAddress = "127.0.0.1:7778"

//BufSize The buffer size for receiving data
var BufSize = 2048

var (
	conn *net.UDPConn
)

//obtains the packet id from the packet
type PacketID struct {
	ID int
}

//obtains input from user
type MovePacket struct {
	e *Entity
	X int
	Y int
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

		var parsedPacketID PacketID
		err = msgpack.Unmarshal(buf[:n], &parsedPacketID)

		//TODO - Put player creation on the main thread
		p := GetPlayer(addr)
		if p == nil {
			p = NewPlayer(addr)
			println("New Player Connected!")
		}
		//println(parsedPacketID.ID)
		if parsedPacketID.ID == 1 {
			var movePacket MovePacket
			err = msgpack.Unmarshal(buf[:n], &movePacket)
			movePacket.e = p.e
			MovePacketChan <- movePacket
			// err = msgpack.Unmarshal(buf[:n], &input)
			// p.e.X += float32(input.X)
			// p.e.Y += float32(input.Y)
		}

		// p.e.X = data.X
		// p.e.Y = data.Y
		// p.e.Z = data.Z
	}
	println("break")
}

//Sends data to players
func Send() {
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
