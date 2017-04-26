//go:generate msgp
package main

import (
	"fmt"
	"net"
	"os"

	msgpack "gopkg.in/vmihailenco/msgpack.v2"
)

//ListenAddress The listen address for the server
var ListenAddress = "127.0.0.1:7778"

//BufSize The buffer size for receiving data
var bufSize = 2048

var (
	conn *net.UDPConn
)

//GetPacketID Get's the packet id from a packet
type GetPacketID struct {
	ID int
}

//InputPacket Get's input data from packet
type InputPacket struct {
	ID     int
	X      int
	Y      int
	entity *Entity
}

//StatePacket sends whole state to a player
type StatePacket struct {
	Entities []Entity
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

	buf := make([]byte, bufSize)

	for {
		n, addr, err := conn.ReadFromUDP(buf)
		CheckError(err)

		var pID GetPacketID
		err = msgpack.Unmarshal(buf[:n], &pID)
		CheckError(err)

		p := GetPlayer(addr)
		if p == nil {
			p = NewPlayer(addr)
		}

		if pID.ID == 1 { //Movement input
			//Unpack data
			var packet InputPacket
			err = msgpack.Unmarshal(buf[:n], &packet)
			//Send it to the channel
			packet.entity = p.e
			inputChan <- packet
			println(len(inputChan))
		}

		// p.e.X = data.X
		// p.e.Y = data.Y
		// p.e.Z = data.Z
	}
}

//Send Sends data to players
func Send() {

	for _, p := range players {

		var entitiesToSend = make([]Entity, len(players), len(players))
		var statePacket StatePacket
		i := 0

		for _, p2 := range players {
			entitiesToSend[i] = *p2.e
			i++
		}
		statePacket.Entities = entitiesToSend
		//---vmihailenco---
		b, err := msgpack.Marshal(statePacket)
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
