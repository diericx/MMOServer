package main

//4262452001132560
//08/18
//683

import "math"

//map data
var m = make(map[int]map[string]*Entity)
var CELL_SIZE = 1000

//CELL DATA
func (e *Entity) calcKey() int {
	var xCell = math.Floor(e.body.Position().x / float64(CELL_SIZE))
	var yCell = math.Floor(e.body.Position().y / float64(CELL_SIZE))
	var key int = int(xCell)*1000 + int(yCell)
	return key
}

func (e *Entity) updateEntityCellData() {
	var oldKey = e.key
	var freshKey = e.calcKey()

	if oldKey == freshKey {
		return
	}

	//println("Old Key: ", oldKey, ", New Key: ", freshKey)

	//remove entity from old array
	if m[oldKey][e.id.String()] != nil {
		removeFromMap(oldKey, e.id.String())
	}

	//if new map in new Position doesnt exist, create it
	if m[freshKey] == nil {
		m[freshKey] = make(map[string]*Entity)
	}

	m[freshKey][e.id.String()] = e

	e.key = freshKey
}

func removeFromMap(key int, id string) {
	//remove entity from old key array
	delete(m[key], id)
	//if old key array is empty, remove it
	if len(m[key]) == 0 {
		delete(m, key)
	}
}

func (e *Entity) findNearestPlayer(maxDist float64) *Entity {
	var nearestPlayer *Entity
	var minDistance float64 = 9999999
	var keys = e.getNearbyKeys()
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

func (e *Entity) getNearbyKeys() []int {
	var length = 2

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
