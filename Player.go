package main

import (
	"math/rand"
	"net"
)

type Player struct {
	//server stuff
	id   string
	addr *net.UDPAddr
}

var players = make(map[string]*Entity)

var PLAYER_EXPIRE_TIME int = 60

var PLAYER_SIZE float64 = 1
var HEALTH_MOD int = 20
var SPEED_MOD float64 = 0.5
var FIRERATE_MOD int = 2

func NewPlayer(addr *net.UDPAddr, pos Vect2, size Vect2) *Entity {
	p := NewEntity(pos, size)
	p.addr = addr
	p.entityType = "player"
	p.resourceId = "empty"
	//funcs
	p.onUpdate = p.playerUpdateFunc
	p.onCollide = p.playerOnCollide
	//expire
	p.expireCounter = PLAYER_EXPIRE_TIME
	//create new planet for the player
	newPlanet := NewPlanet(Vect2{rand.Float64() * 50, rand.Float64() * 50}, Vect2{1, 1})
	newPlanet.Origin(p)
	p.body.targetPos = Vect2{newPlanet.body.pos.x, newPlanet.body.pos.y}
	p.possessedEntities = append(p.possessedEntities, newPlanet)
	players[addr.String()] = p

	return p
}

func (e *Entity) playerUpdateFunc() {
	e.detectCollisions()
}

func (e *Entity) playerOnCollide(other *Entity) {

}
