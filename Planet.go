package main

import (
	"fmt"
	"time"
)

var planets = make(map[int]*Entity)

func NewPlanet(pos Vect3, size Vect3) *Entity {
	p := NewEntity(pos, size)
	p.entityType = "planet"
	p.resourceId = "default_planet"
	//funcs
	p.onUpdate = p.planetUpdateFunc
	p.onCollide = p.planetOnCollide
	planets[p.id] = p
	return p
}

func (e *Entity) planetUpdateFunc() {
	var now = time.Now()
	var dLastCountUpdate time.Duration = now.Sub(e.stats.LastCountUpdate)
	if e.origin != nil {
		if dLastCountUpdate >= e.stats.CountUpdateCooldown {
			e.SetCount(e.stats.Count + 1)
			e.stats.LastCountUpdate = time.Now()
			fmt.Printf("My PLante: %v, %v, %v, %v \n", e.id, e.Position().x, e.Position().y, e.Position().z)
		}
	}
}

func (e *Entity) planetOnCollide(other *Entity) {

}
