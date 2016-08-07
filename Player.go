package main

import "net"

var players = make(map[string]*Entity)

func NewPlayer(addr *net.UDPAddr, pos Vect2, size Vect2) *Entity {
	p := NewEntity(pos, size)
	p.addr = addr
	p.entityType = "player"

	players[addr.String()] = p

	return p
}
