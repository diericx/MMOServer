package main

import (
	"net"
)

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
	p.resourceId = "default_player"
	//funcs
	p.onUpdate = p.playerUpdateFunc
	p.onCollide = p.playerOnCollide
	//expire
	p.expireCounter = PLAYER_EXPIRE_TIME

	players[addr.String()] = p

	return p
}

func (e *Entity) playerUpdateFunc() {
	e.detectCollisions()
}

func (e *Entity) playerOnCollide(other *Entity) {

}
