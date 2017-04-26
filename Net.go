//go:generate msgp
package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"

	msgpack "gopkg.in/vmihailenco/msgpack.v2"

	codec "github.com/ugorji/go/codec"
)

//ListenAddress The listen address for the server
var ListenAddress = "127.0.0.1:7778"

//BufSize The buffer size for receiving data
var BufSize = 2048

//FrameWaitTime time between frames
var FrameWaitTime float64 = 33

var (
	conn *net.UDPConn
)

var (
	v      interface{} // value to decode/encode into
	reader io.Reader
	writer io.Writer
	b      []byte
	mh     codec.MsgpackHandle
)

type InputPacket struct {
	Id int
	X  int
	Y  int
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

	buf := make([]byte, BufSize)

	for {
		n, addr, err := conn.ReadFromUDP(buf)
		CheckError(err)

		var packet InputPacket

		dec := codec.NewDecoder(reader, &mh)
		dec = codec.NewDecoderBytes(buf[:n], &mh)
		err = dec.Decode(&packet)
		println(packet.X)

		p := GetPlayer(addr)
		if p == nil {
			p = NewPlayer(addr)
		}

		println(p.e.Id)

		// p.e.X = data.X
		// p.e.Y = data.Y
		// p.e.Z = data.Z
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

				t := TestPacket{TestFloat: 69}

				//---UGORJI---
				// enc := codec.NewEncoder(writer, &mh)
				// enc = codec.NewEncoderBytes(&b, &mh)
				// err := enc.Encode(t)
				// println(err)

				//---vmihailenco---
				b, err := msgpack.Marshal(t)
				CheckError(err)

				// data, _ := p2.e.MarshalMsg(nil)
				// b, err := msgpack.Marshal(p2.e)
				// CheckError(err)

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
