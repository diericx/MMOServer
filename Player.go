package main

import (
	"math"
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

func (e *Entity) playerShoot() {
	var b = NewBullet(e.body.pos, Vect2{x: 0.5, y: 1}, e)
	b.body.angle = e.body.angle + (math.Pi / 2)
	b.body.vel = Vect2{x: math.Cos(b.body.angle) * e.stats_calc.BulletSpeed, y: math.Sin(b.body.angle) * e.stats_calc.BulletSpeed}
}

func (e *Entity) playerUpdateFunc() {
	e.detectCollisions()

	if e.stats.FireCoolDown >= 0 {
		e.stats.FireCoolDown -= 1
	}

	if e.shooting {
		if e.stats.FireCoolDown <= 0 {
			//if able to shoot, call shoot function
			e.playerShoot()
			e.stats.FireCoolDown = e.stats.FireRate
		}
	}
}

func (e *Entity) playerOnCollide(other *Entity) {
	if other.entityType == "bullet" {
		e.stats.Health -= other.value
	} else if other.entityType == "item-stat-alter" {
		e.stats = e.stats.combine(other.stats)
	} else if other.entityType == "item-pickup" {
		e.addItemToInventory(e.inventory[0])
	}
}
