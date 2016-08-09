package main

import (
	"math"
	"net"
)

var players = make(map[string]*Entity)

var PLAYER_SIZE float64 = 1

func NewPlayer(addr *net.UDPAddr, pos Vect2, size Vect2) *Entity {
	p := NewEntity(pos, size)
	p.addr = addr
	p.entityType = "player"
	//funcs
	p.onUpdate = p.playerUpdateFunc
	p.onCollide = p.playerOnCollide

	players[addr.String()] = p

	return p
}

func (e *Entity) playerShoot() {
	var b = NewBullet(e.body.pos, e.body.size, e)
	b.body.angle = e.body.angle + (math.Pi / 2)
	b.body.vel = Vect2{x: math.Cos(b.body.angle) * e.stats.bulletSpeed, y: math.Sin(b.body.angle) * e.stats.bulletSpeed}
}

func (e *Entity) playerUpdateFunc() {
	e.detectCollisions()

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

func (e *Entity) playerOnCollide(other *Entity) {
	if other.entityType == "bullet" {
		println("COLLIDE WITH BULLET", other.value, e.Health())
		e.stats.health -= other.value
	} else if other.entityType == "item" {
		e.stats = e.stats.combine(other.stats)
	}
}
