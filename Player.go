package main

import (
	"math"
	"net"
)

var players = make(map[string]*Entity)

func NewPlayer(addr *net.UDPAddr, pos Vect2, size Vect2) *Entity {
	p := NewEntity(pos, size)
	p.addr = addr
	p.entityType = "player"
	p.onUpdate = p.playerUpdateFunc

	players[addr.String()] = p

	return p
}

func (e *Entity) playerShoot() {
	var b = NewBullet(e.body.pos, e.body.size)
	b.body.angle = e.body.angle + (math.Pi / 2)
	b.body.vel = Vect2{x: math.Cos(b.body.angle) * 100, y: math.Sin(b.body.angle) * 100}
}

func (e *Entity) playerUpdateFunc() {
	if e.stats.shootCoolDown >= 0 {
		e.stats.shootCoolDown -= 1
	}

	if e.shooting {
		if e.stats.shootCoolDown <= 0 {
			//if able to shoot, call shoot function
			e.playerShoot()
			e.stats.shootCoolDown = e.stats.shootTime
		}
	}
}
