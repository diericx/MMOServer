package main

var ships = make(map[int]map[int]*Entity)

var shipSpeed float64 = 0.2

func NewShip(pos Vect3, size Vect3, t *Entity) *Entity {
	p := NewEntity(pos, size)
	p.addr = addr
	p.entityType = "ship"
	p.resourceId = "ship"
	p.target = t

	p.onUpdate = p.shipUpdateFunc
	p.onCollide = p.shipOnCollide
}

func (e *Entity) shipUpdateFunc() {
	var Vect3 trans = Vect3{0, 0, 0}
	//calc x translation
	if e.Position().x < e.target.Position().x {
		trans.x = shipSpeed
	} else if e.Position().x > e.target.Position().x {
		trans.x = -shipSpeed
	}
	//calc y translation
	if e.Position().y < e.target.Position().y {
		trans.y = shipSpeed
	} else if e.Position().y > e.target.Position().y {
		trans.y = -shipSpeed
	}
	//calc x translation
	if e.Position().z < e.target.Position().z {
		trans.z = shipSpeed
	} else if e.Position().x > e.target.Position().z {
		trans.z = -shipSpeed
	}

}

func (e *Entity) shipOnCollide(other *Entity) {

}
