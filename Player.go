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
	b.body.vel = Vect2{x: math.Cos(b.body.angle) * e.stats_calc1.BulletSpeed, y: math.Sin(b.body.angle) * e.stats_calc1.BulletSpeed}
}

func (e *Entity) playerUpdateFunc() {
	e.detectCollisions()

	if e.stats_calc1.FireCoolDown >= 0 {
		e.stats_calc1.FireCoolDown -= 1
	}

	if e.shooting {
		if e.stats_calc1.FireCoolDown <= 0 {
			//if able to shoot, call shoot function
			e.playerShoot()
			e.stats_calc1.FireCoolDown = e.stats_calc1.FireRate
		}
	}
}

func (e *Entity) playerOnCollide(other *Entity) {
	if other.entityType == "bullet" {
		println("collided with bullet!")
		var damageToTake = other.value - float64(e.stats_calc1.Defense)
		if damageToTake <= 0 {
			damageToTake = 1
		}
		e.stats_calc1.Health -= damageToTake
	} else if other.entityType == "item-stat-alter" {
		e.stats_calc1 = e.stats_calc1.add(other.stats_calc1)
		e.extendedDataHash = e.generateExtendedDataHash()
	} else if other.entityType == "item-pickup" {
		e.addItemToInventory(other.inventory[0])
	}
}
