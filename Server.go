package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"ugorji/go/codec"

	"golang.org/x/net/websocket"
)

type Server struct {
	addr string
}

type Socket struct {
	io.ReadWriter
	done   chan bool
	closed bool
}

type SendMessage struct {
	CurrentPlayer string
	Objects       []Data
	Leaderboard   [5]LeaderboardEntry
}

type LeaderboardEntry struct {
	Username string
	Score    int
}

type ReceiveMessage struct {
	Action   string
	Token    string
	Uid      string
	Value    string
	Angle    float64
	X        int
	Y        int
	Shooting bool
}

type TargetData struct {
	Distance int
	Angle    float64
}

type Data struct {
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
	Rot       float64
}

type CurrentPlayer struct {
	Id        string
	Username  string
	Parent    string
	Child     string
	Height    int
	Score     int
	Health    float64
	HealthCap int
	Power     int
	MaxPower  int
	X         float64
	Y         float64
	Rot       float64
	Angle     float64
}

type PlayerDataObject struct {
	action       string
	packetObj    ReceiveMessage
	packetString string
	socket       Socket
	player       *Entity
}

var (
	mh codec.MsgpackHandle
	r  io.Reader
	w  io.Writer
)
var firebaseUrl = "https://diericx.firebaseIO.com/"

//---
//SOCKET---
//---
func (s Socket) Close() error {
	s.done <- true

	s.closed = true

	return nil
}

func socketHandler(ws *websocket.Conn) {
	println("Got connection!")
	s := Socket{ws, make(chan bool), false}
	//add the new player to the data stream so it will be created
	serverInput <- PlayerDataObject{action: "newPlayer", socket: s}

	<-s.done
}

//---
//-----
//----

func GetData(p *Entity) {
	for !p.socket.closed {
		//serialize the packet
		buf := make([]byte, 500)
		n, err := p.socket.Read(buf)
		if err != nil {
			println("Player <", p.username, "> disconnected")
			p.socket.Close()
			p.value = 0
			//p.removeSelf()

			return
		}

		var msg = ReceiveMessage{}
		//decode data
		dec := codec.NewDecoder(r, &mh)
		dec = codec.NewDecoderBytes(buf[0:n], &mh)
		err = dec.Decode(&msg)

		//send the data to the stream to be used later
		serverInput <- PlayerDataObject{action: "updatePlayer", player: p, packetObj: msg}
	}
}

func NewServer(addr string) *Server {
	newServer := new(Server)
	newServer.addr = addr
	return newServer
}

func (s Server) listenForPlayers() {
	println("waiting for connection...")
	//http.Handle("/", websocket.Handler(socketHandler))
	//http.ListenAndServe(s.addr, nil) //192.168.2.36
	l, err := net.Listen("tcp", s.addr)

	if err != nil {
		log.Fatal(err)
	}

	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		println("Got connection!")
		s := Socket{c, make(chan bool), false}
		//add the new player to the data stream so it will be created
		serverInput <- PlayerDataObject{action: "newPlayer", socket: s}

		<-s.done
	}
}

//---Process Data---//
func processServerInput() {
	if len(serverInput) == 0 {
		return
	}

	for len(serverInput) > 0 {
		var serverInputObj = <-serverInput
		if serverInputObj.action == "newPlayer" {
			println("New Player...")
			var p = NewPlayer(serverInputObj.socket, Vect2{x: 0, y: 0}, Vect2{x: 40, y: 40})
			go GetData(p)
		} else if serverInputObj.action == "updatePlayer" {
			var p = serverInputObj.player
			var msg = serverInputObj.packetObj

			if msg.Action == "update" {
				//activate user
				if !p.active {
					p.authorize(firebaseUrl+"users/"+msg.Uid, msg.Token)
				}

				var movX float64 = 0
				var movY float64 = 0

				if msg.Shooting == true {
					if p.shoot != nil {
						p.shoot(nil)
					}
				}

				//Edit rotation
				angleInRad := ((msg.Angle * math.Pi) / 180)
				angleInRadForward := angleInRad + (math.Pi / 2)
				angleInRadRight := angleInRad

				movX += math.Cos(angleInRadForward) * float64(msg.Y)
				movY += math.Sin(angleInRadForward) * float64(msg.Y)

				movX += math.Cos(angleInRadRight) * float64(msg.X)
				movY += math.Sin(angleInRadRight) * float64(msg.X)

				//p.moveEntity(Vect2{x: movX * 15, y: movY * 15})
				p.body.vel.x = movX * 15
				p.body.vel.y = movY * 15

				p.body.rotation = angleInRad

			} else if msg.Action == "jumpRight" {
				movX := math.Cos(p.body.rotation)
				movY := math.Sin(p.body.rotation)

				//p.moveEntity(Vect2{x: movX * 5000, y: movY * 5000})
				p.body.vel.x = movX
				p.body.vel.y = movY
			} else if msg.Action == "jumpLeft" {
				movX := math.Cos(p.body.rotation)
				movY := math.Sin(p.body.rotation)
				//p.moveEntity(Vect2{x: -movX * 5000, y: -movY * 5000})
				p.body.vel.x = movX
				p.body.vel.y = movY
			} else if msg.Action == "setParent" {
				var e = findEntityById(msg.Value)
				if e != nil {
					p.setParent(e)
				}
			} else if msg.Action == "equipItem" {
			}
		}
	}
}

func processServerOutput() {
	for _, p := range players {

		if !p.active {
			continue
		}

		var playerDataObj PlayerDataObject
		var objects = []Data{}
		var lb = [5]LeaderboardEntry{}

		var keys = p.getNearbyKeys()

		for _, key := range keys {
			for _, e := range m[key] {
				//if p != e && e != nil {

				//don't send players without a username yet
				// if e.entityType == "Player" {
				// 	if e.username == "" {
				// 		continue
				// 	}
				// }

				var object Data
				//set parent and child if they are there
				if e.parent != nil {
					object.Parent = e.parent.id.String()
				}
				if e.child != nil {
					object.Child = e.child.id.String()
				}
				object.Height = e.height

				object.Id = e.id.String()
				object.Username = e.username
				object.Type = e.entityType
				object.Tag = e.tag
				object.Health = e.health
				object.HealthCap = e.healthCap
				object.Power = e.power
				object.X = float64(e.body.Position().x)
				object.Y = float64(e.body.Position().y)
				object.Rot = e.body.rotation
				objects = append(objects, object)
				//}
			}
		}

		//calculate angles
		// if target != nil {
		// 	var newTargetData = TargetData{Distance: int(p.distanceTo(target)), Angle: p.lookAt(target)}
		// 	targets = append(targets, newTargetData)
		// }

		//set parent and child if they are there
		// if p.parent != nil {
		// 	cp.Parent = p.parent.id.String()
		// }
		// if p.child != nil {
		// 	cp.Child = p.child.id.String()
		// }
		// cp.Height = p.height
		//
		// cp.Id = p.id.String()
		// cp.Username = p.username
		// cp.Score = int(p.score)
		// cp.Health = p.health
		// cp.HealthCap = p.healthCap
		// cp.Power = p.power
		// cp.MaxPower = p.powerMax
		// cp.X = float64(p.body.Position().X)
		// cp.Y = float64(p.body.Position().Y)
		// cp.Rot = 0

		//populate leaderboard
		for k, v := range leaderboard {
			if v != nil {
				lb[k] = LeaderboardEntry{Username: v.username, Score: int(v.score)}
			}
		}

		m := SendMessage{
			p.id.String(),
			objects,
			lb,
		}

		//b, _ := json.Marshal(m)
		var newByteArray []byte
		enc := codec.NewEncoder(w, &mh)
		enc = codec.NewEncoderBytes(&newByteArray, &mh)
		enc.Encode(m)

		//println(string(newByteArray))
		//println(len(m.Objects))

		playerDataObj.player = p
		playerDataObj.action = "sendData"
		playerDataObj.packetString = string(newByteArray)
		serverOutput <- playerDataObj
	}
}

func sendServerOutput() {
	//if there is nothing to send, return
	if len(serverOutput) == 0 {
		return
	}
	//if there is something to send, loop through and send all
	for len(serverOutput) > 0 {
		var serverOutputObj = <-serverOutput
		if serverOutputObj.action == "sendData" {
			fmt.Fprint(serverOutputObj.player.socket, string(serverOutputObj.packetString))
		}
	}
}
