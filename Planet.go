package main

import (
	"math/rand"
	"time"
)

var planets = make(map[int]*Entity)

var planetResources = [7]string{"egipt", "forest", "havay", "ice", "ice_gray", "orange_planet", "pine"}

func NewPlanet(pos Vect3, size Vect3) *Entity {
	p := NewEntity(pos, size)
	p.entityType = "planet"
	p.resourceId = planetResources[rand.Intn(len(planetResources))]
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
		}
	}
}

func (e *Entity) planetOnCollide(other *Entity) {

}
