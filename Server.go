package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"ugorji/go/codec"
)

type ReceivePacket struct {
	Action   string
	Token    string
	Uid      string
	Value    string
	Angle    float64
	X        int
	Y        int
	Shooting bool
}

type SendPacket struct {
	CurrentPlayerId string
	Objects         []EntityData
}

type EntityData struct {
	Id        string
	Username  string
	Parent    string //id
	Child     string //id
	Height    int
	Type      string
	Tag       string
	Health    float64
	HealthCap int
	Power     int
	MaxPower  int
	X         float64
	Y         float64
	Angle     float64
}

type ServerActionObj struct {
	receivePacketObj ReceivePacket
	sendPacketBytes  []byte
	addr             *net.UDPAddr
	entity           *Entity
}

var LISTEN_ADDRESS = "localhost:7777"
var BUF_SIZE = 2048

//variables for decoding
var (
	mh codec.MsgpackHandle
	r  io.Reader
	w  io.Writer
)

var serverInput = make(chan ServerActionObj, 1000)
var serverOutput = make(chan ServerActionObj, 1000)

var serverConn *net.UDPConn

func listenForPackets() {
	/* Lets prepare an address at any address at the listen address*/
	ServerAddr, err := net.ResolveUDPAddr("udp", LISTEN_ADDRESS)
	CheckError(err)

	/* Now listen at selected port */
	serverConn, err = net.ListenUDP("udp", ServerAddr)
	CheckError(err)

	defer serverConn.Close()

	buf := make([]byte, BUF_SIZE)

	for {
		n, addr, err := serverConn.ReadFromUDP(buf)
		CheckError(err)

		p := players[addr.String()]

		if p == nil {
			//player hasn't been instantiated yet
			println("creating player...")
			NewPlayer(addr, Vect2{x: 0, y: 0}, Vect2{x: 10, y: 10})
		} else {
			//player has been instantiated
			var msg = ReceivePacket{}
			//decode data
			dec := codec.NewDecoder(r, &mh)
			dec = codec.NewDecoderBytes(buf[0:n], &mh)
			err = dec.Decode(&msg)
			CheckError(err)
			//send the data to the stream to be processed later
			serverInput <- ServerActionObj{entity: p, receivePacketObj: msg, addr: addr}
		}
	}
}

//---Process input data---//
func processServerInput() {

	for len(serverInput) > 0 {
		var serverInputObj = <-serverInput
		var packet = serverInputObj.receivePacketObj
		var player = serverInputObj.entity

		if packet.Action == "update" {
			player.body.pos.x += float64(packet.X) * 5
			player.body.pos.y += float64(packet.Y) * 5
		}
	}

}

//---Process output data---//
func processServerOutput() {
	for _, p := range players {

		var objects = []EntityData{}

		for _, e := range entities {
			var ed EntityData
			ed.Id = e.id.String()
			ed.Type = e.entityType
			ed.X = e.body.pos.x
			ed.Y = e.body.pos.y
			ed.Angle = e.body.angle
			objects = append(objects, ed)
		}

		packetObj := SendPacket{
			CurrentPlayerId: p.id.String(),
			Objects:         objects,
		}

		//encode send packet
		var packet []byte
		enc := codec.NewEncoder(w, &mh)
		enc = codec.NewEncoderBytes(&packet, &mh)
		enc.Encode(packetObj)

		//add packet to queue of things to send
		var actionObj ServerActionObj
		actionObj.entity = p
		actionObj.sendPacketBytes = packet
		actionObj.addr = p.addr
		serverOutput <- actionObj

	}
}

func sendServerOutput() {
	//if there is nothing to send, return
	if len(serverOutput) == 0 {
		return
	}
	//if there is something to send, loop through and send all
	for len(serverOutput) > 0 {
		var outputObj = <-serverOutput
		sendMessage(outputObj.sendPacketBytes, serverConn, outputObj.addr)
	}
}

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}

func sendMessage(msg []byte, conn *net.UDPConn, addr *net.UDPAddr) {
	_, err := conn.WriteToUDP(msg, addr)
	if err != nil {
		fmt.Printf("Couldn't send response; %v\n", err)
	}
}
