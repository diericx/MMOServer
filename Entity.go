package main

import (
	"math"
	"net"
	"strconv"
)

type Vect2 struct {
	x float64
	y float64
}

type Stats struct {
	Health               float64
	Defense              int
	MaxHealth            int
	FireRate             int
	FireCoolDown         int
	Speed                float64
	BulletSpeed          float64
	Damage               float64
	Energy               int
	NextEnergyCheckpoint int
}

type Entity struct {
	//server stuff
	id         string
	addr       *net.UDPAddr
	hasChanged bool
	//
	key           int
	body          Body
	entityType    string
	origin        *Entity
	active        bool
	resourceId    string
	expireCounter int
	//
	stats Stats
	//user defined functions
	onUpdate  func()
	onCollide func(other *Entity)
	onShoot   func()
	onRemove  func()
}

//holding array
var entityIdIncrement int = 0
var entities = make(map[string]*Entity)

//hash map array
var m = make(map[int]map[string]*Entity)

//channel for entities to remove
var entitiesToRemove = make(chan ServerActionObj, 1000)

var STAT_CALC_MOD = Stats{
	MaxHealth:   25,
	Defense:     1,
	FireRate:    5,
	Speed:       0.2,
	BulletSpeed: 1,
	Damage:      1,
}

//hash cell size
var CELL_SIZE = 15
var INVENTORY_MAX = 10
var energyCheckpoints = []int{}

func NewEntity(pos Vect2, size Vect2) *Entity {
	newEntity := Entity{}
	newEntity.active = true

	newEntity.id = strconv.Itoa(entityIdIncrement)
	newEntity.resourceId = "default_entity"
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
		Health:               100,
		MaxHealth:            4,
		Defense:              1,
		FireRate:             3,
		FireCoolDown:         3,
		Speed:                1,
		BulletSpeed:          1,
		Damage:               1,
		NextEnergyCheckpoint: 0,
	}
	return stats
}

func updateEntities() {
	for _, p := range players {
		if p.onUpdate != nil {
			p.onUpdate()
			p._onUpdate()
		}
	}

}

func removeEntities() {
	if len(entitiesToRemove) <= 0 {
		return
	}

	var e = <-entitiesToRemove
	e.entity.RemoveSelf()
}

func (e *Entity) _onUpdate() {
	e.updateEntityCellData()

	//check for expiration
	if e.expireCounter > 0 {
		e.expireCounter -= 1
	} else if e.expireCounter == 0 {
		e.RemoveSelf()
	}

	//p.moveEntity(Vect2{x: movX * 15, y: movY * 15})
	e.body.pos.x += e.body.vel.x
	e.body.pos.y += e.body.vel.y

	if e.Health() <= 0 {
		e.Die()
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
	e.stats.Health = 100
	e.stats.Energy = e.stats.Energy / 2
}

func (e *Entity) RemoveSelf() {
	if e.onRemove != nil {
		e.onRemove()
	}

	delete(entities, e.id)
	delete(players, e.addr.String())
	removeFromMap(e.key, e.id)
}

//------Helper functions with body--------

//body
func (e *Entity) SetPosition(x float64, y float64) {
	e.body.pos.x = x
	e.body.pos.y = y
}

func (e *Entity) Position() Vect2 {
	return e.body.pos
}

func (e *Entity) Health() float64 {
	return e.stats.Health
}

//stats
func (s Stats) add(s2 Stats) Stats {
	s.BulletSpeed += s2.BulletSpeed
	s.Energy += s2.Energy
	s.Health += s2.Health
	s.Defense += s2.Defense
	s.MaxHealth += s2.MaxHealth
	s.FireRate += s2.FireRate
	s.Speed += s2.Speed
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
		m[c] = make(map[string]*Entity)
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
func removeFromMap(key int, id string) {
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
