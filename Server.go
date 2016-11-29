package main

import (
	"fmt"
	"io"
	"net"
	"os"

	"github.com/ugorji/go/codec"
)

type ReceivePacket struct {
	Action   string
	Token    string
	Uid      string
	Value    string
	Angle    float64
	IDs      []int
	X        int
	Y        int
	Shooting bool
}

type UpdatePacket struct {
	Action          string
	CurrentPlayerId int
	Objects         []EntityData
	ObjectsMin      []EntityDataMin
}

type EntityExtendedDataPacket struct {
	Action string
	Id     string
	EED    EntityExtendedData
}

type EntityExtendedData struct {
	ExtendedDataHash uint32
	StatsObj         Stats
	StatUpgrades     int
}

type EntityData struct {
	Id         int
	OriginId   int
	Type       string
	ResourceId string
	Count      int
	X          float32
	Y          float32
}

type EntityDataMin struct {
	Id int
}

type ServerActionObj struct {
	receivePacketObj ReceivePacket
	sendPacketBytes  []byte
	addr             *net.UDPAddr
	entity           *Entity
}

var LISTEN_ADDRESS = ":7777"
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
			NewPlayer(addr, Vect2{x: 0, y: 0}, Vect2{x: PLAYER_SIZE, y: PLAYER_SIZE})
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
			//update expire
			player.expireCounter = PLAYER_EXPIRE_TIME
			player.SetPosition(float64(packet.X), float64(packet.Y))
			//data requests
			for _, id := range packet.IDs {
				player.dataRequests[id] = true
			}

		} else if packet.Action == "select" {
			player.selectedEntities = []int{}
			for _, id := range packet.IDs {
				if entities[id].origin == player {
					println("Selected: ", id)
					player.selectedEntities = append(player.selectedEntities, id)
				}
			}
		} else if packet.Action == "attack" {
			println("Attack: ", packet.IDs[0])
			var planetToAttack = packet.IDs[0]
			player.attackPlanet(planetToAttack)
		}
	}

}

//---Process output data---//
func processServerOutput() {
	for _, p := range players {

		var objects = []EntityData{}
		var objectsMin = []EntityDataMin{}

		var keys = p.getNearbyKeys(2)

		for _, key := range keys {
			for _, e := range m[key] {
				//if it hasnt changed, add the min data to packet and cont.
				if (changedEntities[e.id] == false && p.dataRequests[e.id] == false && e.entityType != "player") || len(objects) > 10 {
					var ed EntityDataMin
					ed.Id = e.id
					objectsMin = append(objectsMin, ed)
					continue
				}
				//if its changed, add all its data to packet
				var ed EntityData
				ed.Id = e.id
				if e.origin != nil {
					ed.OriginId = e.origin.id
				} else {
					ed.OriginId = -1
				}
				ed.Type = e.entityType
				ed.ResourceId = e.resourceId
				ed.Count = e.stats.Count
				if e == p {
					//Move the player to its target position if there is one
					if e.body.targetPos.x != 0 {
						ed.X = float32(e.body.targetPos.x)
					}
					if e.body.targetPos.y != 0 {
						ed.Y = float32(e.body.targetPos.y)
					}
					if int(e.body.pos.x) == int(e.body.targetPos.x) {
						e.body.targetPos.x = 0
					}
					if int(e.body.pos.y) == int(e.body.targetPos.y) {
						e.body.targetPos.y = 0
					}
				} else {
					ed.Y = float32(e.Position().y)
					ed.X = float32(e.Position().x)
				}

				//update data requests
				p.dataRequests[e.id] = false

				objects = append(objects, ed)
			}
		}

		packetObj := UpdatePacket{
			Action:          "update",
			CurrentPlayerId: p.id,
			Objects:         objects,
			ObjectsMin:      objectsMin,
		}

		//encode send packet
		var packet []byte
		enc := codec.NewEncoder(w, &mh)
		enc = codec.NewEncoderBytes(&packet, &mh)
		enc.Encode(packetObj)
		//println("Packet length: ", len(packet))

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
