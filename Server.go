package main

import (
	"fmt"
	"io"
	"math"
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
	X        int
	Y        int
	Shooting bool
}

type SendPacket struct {
	CurrentPlayerId string
	Objects         []EntityData
}

type EntityData struct {
	Id         string
	Username   string
	Type       string
	ResourceId string
	Energy     int
	Health     float64
	X          float64
	Y          float64
	Angle      float64
	//current player data only
	Inventory    []Item
	StatsObj     Stats
	StatUpgrades int
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
			var movX float64 = 0
			var movY float64 = 0

			//update expire
			player.expireCounter = PLAYER_EXPIRE_TIME

			//calculate angles
			angleInRad := -((packet.Angle * math.Pi) / 180)
			angleInRadForward := angleInRad + (math.Pi / 2)
			angleInRadRight := angleInRad

			//calculate velocity
			movX += math.Cos(angleInRadForward) * float64(packet.Y)
			movY += math.Sin(angleInRadForward) * float64(packet.Y)
			movX += math.Cos(angleInRadRight) * float64(packet.X)
			movY += math.Sin(angleInRadRight) * float64(packet.X)

			//edit body velocity
			player.body.vel.x = movX * player.stats.Speed
			player.body.vel.y = movY * player.stats.Speed

			//set angle and shooting variable
			player.body.angle = angleInRad
			player.shooting = serverInputObj.receivePacketObj.Shooting
		} else if packet.Action == "attempt-equip" {
			player.attemptToEquip(int(packet.X))
		} else if packet.Action == "upgradeHealth" {
			if player.getAvailableUpgrades() > 0 {
				s := Stats{}
				s.MaxHealth = HEALTH_MOD
				player.stats = player.stats.combine(s)
				player.statsUpgrades = append(player.statsUpgrades, s)
			}
		} else if packet.Action == "upgradeSpeed" {
			if player.getAvailableUpgrades() > 0 {
				s := Stats{}
				s.Speed = SPEED_MOD
				player.stats = player.stats.combine(s)
				player.statsUpgrades = append(player.statsUpgrades, s)
			}
		} else if packet.Action == "upgradeFireRate" {
			if player.getAvailableUpgrades() > 0 {
				println("UP FIRE RATE")
				s := Stats{}
				s.FireRate = FIRERATE_MOD
				player.stats = player.stats.combine(s)
				player.statsUpgrades = append(player.statsUpgrades, s)
			}
		}
	}

}

//---Process output data---//
func processServerOutput() {
	for _, p := range players {

		var objects = []EntityData{}

		var keys = p.getNearbyKeys(2)

		for _, key := range keys {
			//println("Key: ", key)
			for _, e := range m[key] {
				var ed EntityData
				ed.Id = e.id.String()
				ed.Energy = int(e.stats.Energy)
				ed.ResourceId = e.resourceId
				ed.Type = e.entityType
				ed.Health = e.Health()
				ed.X = e.body.pos.x
				ed.Y = e.body.pos.y
				ed.Angle = e.body.angle

				if e == p {
					//edit stats
					s := e.stats
					s.NextEnergyCheckpoint = energyCheckpoints[s.NextEnergyCheckpoint]
					//send player data
					ed.StatsObj = s
					ed.Inventory = e.inventory
					ed.StatUpgrades = e.getAvailableUpgrades()
				}
				objects = append(objects, ed)
			}
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
