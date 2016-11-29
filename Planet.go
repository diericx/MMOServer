package main

import "time"

var planets = make(map[string]*Entity)

func NewPlanet(pos Vect2, size Vect2) *Entity {
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
			e.stats.Count += 1
			e.stats.LastCountUpdate = time.Now()
			println(e.stats.Count)
		}
	}
}

func (e *Entity) planetOnCollide(other *Entity) {

}
