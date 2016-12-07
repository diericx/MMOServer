package main

import (
	"math"
	"net"
	"time"
)

type Stats struct {
	Count               int
	LastCountUpdate     time.Time
	CountUpdateCooldown time.Duration
}

type Entity struct {
	//server stuff
	id           int
	addr         *net.UDPAddr
	hasChanged   bool
	dataRequests map[int]bool
	//
	key           int
	body          Body
	entityType    string
	origin        *Entity
	target        *Entity
	active        bool
	resourceId    string
	expireCounter int
	//
	stats             Stats
	selectedEntities  []int
	possessedEntities []*Entity
	//user defined functions
	onUpdate  func()
	onCollide func(other *Entity)
	onShoot   func()
	onRemove  func()
}

//holding array
var entityIdIncrement int = 0
var entities = make(map[int]*Entity)
var changedEntities = make(map[int]bool)

//hash map array
var m = make(map[int]map[int]*Entity)

//channel for entities to remove
var entitiesToRemove = make(chan ServerActionObj, 1000)

//hash cell size
var CELL_SIZE = 75
var INVENTORY_MAX = 10
var energyCheckpoints = []int{}

func NewEntity(pos Vect2, size Vect2) *Entity {
	newEntity := Entity{}
	newEntity.active = true
	newEntity.dataRequests = make(map[int]bool)

	newEntity.id = entityIdIncrement
	newEntity.resourceId = "default_planet"
	newEntity.body.pos = pos
	newEntity.body.size = size
	newEntity.stats = NewDefaultBaseStats()
	newEntity.body.points = make([]Vect2, 4)
	newEntity.expireCounter = -1
	//gen equipped hash

	newEntity.addToCell(newEntity.calcKey())

	entities[newEntity.id] = &newEntity

	entityIdIncrement += 1

	return &newEntity
}

//default stat values
func NewDefaultBaseStats() Stats {
	stats := Stats{
		Count:               0,
		CountUpdateCooldown: 1500 * time.Millisecond,
	}
	return stats
}

func updateEntities() {

	for _, p := range players {
		//p.hasChanged = false
		if p.onUpdate != nil {
			p.onUpdate()
			p._onUpdate()
		}
	}

	for _, p := range planets {
		//p.hasChanged = false
		if p.onUpdate != nil {
			p.onUpdate()
			p._onUpdate()
		}
	}

}

func (e *Entity) attackPlanet(pID int) {
	var p2a = entities[pID]
	for _, p := range e.selectedEntities {
		entities[p].stats.Count = entities[p].stats.Count / 2
		p2a.stats.Count += entities[p].stats.Count
		p2a.SetOrigin(e)
	}
}

func (e *Entity) SetOrigin(o *Entity) {
	var prevO = e.origin
	e.origin = o
	if prevO != e.origin {
		changedEntities[e.id] = true
	}
}

func (e *Entity) SetCount(c int) {
	var prevC = e.stats.Count
	e.stats.Count = c
	if prevC != e.stats.Count {
		changedEntities[e.id] = true
	}
}

func (e *Entity) _onUpdate() {
	e.updateEntityCellData()

	//check for expiration
	if e.expireCounter > 0 {
		e.expireCounter -= 1
	} else if e.expireCounter == 0 {
		e.RemoveSelf()
	}

	e.SetPosition(e.body.pos.x+e.body.vel.x, e.body.pos.y+e.body.vel.y)

	if e.stats.Count <= 0 {
		e.SetOrigin(nil)
	}
}

func (e *Entity) distanceTo(e2 *Entity) float64 {
	var loc1 = e.Position()
	var loc2 = e2.Position()
	var deltaX = float64(loc2.x - loc1.x)
	var deltaY = float64(loc2.y - loc1.y)
	return math.Sqrt((deltaX * deltaX) + (deltaY * deltaY))
}

func (e *Entity) Die() {
	e.SetPosition(0, 0)
}

func (e *Entity) RemoveSelf() {
	if e.onRemove != nil {
		e.onRemove()
	}

	delete(entities, e.id)
	delete(players, e.addr.String())
	removeFromMap(e.key, e.id)
}

func removeEntities() {
	if len(entitiesToRemove) <= 0 {
		return
	}

	var e = <-entitiesToRemove
	e.entity.RemoveSelf()
}

func resetVariables() {
	changedEntities = make(map[int]bool)
}

//------Helper functions with body--------

//body
func (e *Entity) SetPosition(x float64, y float64) {
	var prevPos = e.body.pos
	e.body.pos.x = x
	e.body.pos.y = y
	if prevPos != e.body.pos {
		changedEntities[e.id] = true
	}
}

func (e *Entity) Position() Vect2 {
	return e.body.pos
}

//stats
func (s Stats) add(s2 Stats) Stats {

	return s
}

//------HASH MAP--------

//CELL DATA
func (e *Entity) calcKey() int {
	var xCell = math.Floor(e.Position().x / float64(CELL_SIZE))
	var yCell = math.Floor(e.Position().y / float64(CELL_SIZE))
	var key int = int(xCell)*1000 + int(yCell)
	return key
}

//add to map
func (e *Entity) addToCell(c int) {
	//if new map in new Position doesnt exist, create it
	if m[c] == nil {
		m[c] = make(map[int]*Entity)
	}

	m[c][e.id] = e
	e.key = c
}

//if the entity is in a new cell, update it's cell data
func (e *Entity) updateEntityCellData() {
	var oldKey = e.key
	var freshKey = e.calcKey()

	if oldKey == freshKey {
		return
	}

	//remove entity from old array
	if m[oldKey][e.id] != nil {
		removeFromMap(oldKey, e.id)
	}

	//add to new array (cell)
	e.addToCell(freshKey)

	e.key = freshKey
}

//remove entity from map
func removeFromMap(key int, id int) {
	//remove entity from old key array
	delete(m[key], id)
	//if old key array is empty, remove it
	if len(m[key]) == 0 {
		delete(m, key)
	}
}

//find nearest player in cell
func (e *Entity) findNearestPlayer(maxDist float64) *Entity {
	var nearestPlayer *Entity
	var minDistance float64 = 9999999
	var keys = e.getNearbyKeys(2)
	for _, key := range keys {
		for _, other := range m[key] {
			if e != other && other != nil && other.entityType == "Player" {
				var dist = e.distanceTo(other)
				if dist < minDistance && dist < maxDist {
					nearestPlayer = other
					minDistance = dist
				}
			}
		}
	}
	return nearestPlayer
}

//get keys of cells around current cell
func (e *Entity) getNearbyKeys(length int) []int {

	var keys = []int{
		e.key,
	}

	//top left cell
	var startCell = e.key - (1000 * length) + (length)

	//loop through horizontally
	for i := 0; i < (length*2)+1; i++ {
		//lop vertically and get all cells below this cell
		for j := 0; j < (length*2)+1; j++ {
			var newKey = startCell - j
			if newKey == e.key {
				continue
			}
			keys = append(keys, newKey)
		}
		startCell += 1000
	}

	return keys
}
