package main

import (
	// "fmt"
	"math"
)

var m = make(map[int][]*Entity)

//!!! Make these functions like player.updateCellData()
// func updatePlayerCellData(player *Player) {
// 	var key = getEntityKey(Vector2{x: player.entity.rect.pos.x, y: player.entity.rect.pos.y})
// 	fmt.Println("PlayerCell: ", key)
// 	if player.entity.key != key {
// 		// fmt.Println("Changing cells...")
// 		// fmt.Println("Length of last cell:", len(m[player.key]), "\n")
// 		// fmt.Println("Length of new cell :", len(m[key]), "\n")
// 		//move player to the correct key list
// 		//m[player.key] = append(m[:player.key], m[player.key+1:]...)
// 		m[key] = append(m[key], player.entity)
// 		//print("\n", len(m), "\n")
// 		//remove player from old list
// 		var indexInOldList = findEntityInCellByID(player.entity.key, player.entity.id)
// 		//fmt.Println("inded in old list  :", indexInOldList, "\n")
// 		if indexInOldList >= 0 {
// 			//var oldList = m[player.key]
// 			//print("\nLength of old list:", len(oldList))
// 			m[player.entity.key] = append(m[player.entity.key][:indexInOldList], m[player.entity.key][indexInOldList+1:]...)
// 			//print("\nLength of old list P2:", len(oldList))
// 			if len(m[player.entity.key]) == 0 {
// 				delete(m, player.entity.key)
// 			}
// 		}
// 		player.entity.key = key
// 	} else {
// 		//do not move payer, do nothing
// 	}
// }

// //!!! Make these functions like npc.updateCellData()
// func updateNPCCellData(npc *Npc) {
// 	//var key = getEntityKey(Vector2{x: npc.rect.x, y: npc.rect.y})
// 	var key = getEntityKey(npc.entity.origin)
// 	if npc.entity.id == "-1" {
// 		// fmt.Println("NPC -1's key is           : ", npc.key)
// 		// fmt.Println("NPC -1's calculated key is: ", key)
// 		// fmt.Println("Are key and Calculated sam: ", npc.key != key)
// 		// fmt.Println("Length of NPC -1's cell is: ", len(m[key]))
// 	}
// 	if npc.entity.key != key {
// 		//move npc to the correct key list
// 		//m[player.key] = append(m[:player.key], m[player.key+1:]...)
// 		m[key] = append(m[key], npc.entity)
// 		var indexInOldList = findEntityInCellByID(npc.entity.key, npc.entity.id)
// 		if npc.entity.id == "-1" {
// 			// fmt.Println("NPC -1's key is          : ", npc.key)
// 			// fmt.Println("NPC -1's calculated key is: ", key)
// 			// fmt.Println("Are key and Calculated sam: ", npc.key == key)
// 			fmt.Println("Length of NPC -1's cell is: ", len(m[key]))
// 			fmt.Println(indexInOldList)
// 		}
// 		//fmt.Println(len(m[key]))
// 		//remove player from old list

// 		if indexInOldList >= 0 {
// 			//var oldList = m[npc.key]
// 			//print("\nLength of old list:", len(oldList))
// 			m[npc.entity.key] = append(m[npc.entity.key][:indexInOldList], m[npc.entity.key][indexInOldList+1:]...)
// 			//print("\nLength of old list P2:", len(oldList))
// 			if len(m[npc.entity.key]) == 0 {
// 				delete(m, npc.entity.key)
// 			}
// 		}
// 		npc.entity.key = key
// 	} else {
// 		//do not move npc, do nothing
// 	}
// }

func updateEntityCellData(e *Entity) {
	var key = getEntityKey(Vector2{x: e.rect.pos.x, y: e.rect.pos.y})
	if e.key != key {
		// fmt.Println("Changing cells...")
		// fmt.Println("Length of last cell:", len(m[player.key]), "\n")
		// fmt.Println("Length of new cell :", len(m[key]), "\n")
		//move player to the correct key list
		//m[player.key] = append(m[:player.key], m[player.key+1:]...)
		m[key] = append(m[key], e)
		//print("\n", len(m), "\n")
		//remove player from old list
		var indexInOldList = findEntityInCellByID(e.key, e.id)
		//fmt.Println("inded in old list  :", indexInOldList, "\n")
		if indexInOldList >= 0 {
			//var oldList = m[player.key]
			//print("\nLength of old list:", len(oldList))
			m[e.key] = append(m[e.key][:indexInOldList], m[e.key][indexInOldList+1:]...)
			//print("\nLength of old list P2:", len(oldList))
			if len(m[e.key]) == 0 {
				delete(m, e.key)
			}
		}
		e.key = key
	} else {
		//do not move payer, do nothing
	}
}

func findEntityInCellByID(key int, id string) int {
	var value = -1
	var i = 0
	for _, e := range m[key] {
		if e.id == id {
			value = i
		}
		i++
	}
	return value
}

func getEntityObject(key int, index int) *Entity {
	var entity = m[key][index]
	return entity
}

func removeEntityFromCell(key int, id string) {
	var keyInCell = findEntityInCellByID(key, id)
	//fmt.Println("cell: ", key, "Key in cell: ", keyInCell, "id: ", id)
	m[key] = append(m[key][:keyInCell], m[key][keyInCell+1:]...)
}

func getEntityKey(pos Vector2) int {
	var xCell = math.Floor(pos.x / CELL_SIZE) //Cell size = 30
	var yCell = math.Floor(pos.y / CELL_SIZE)
	var key int = int(xCell)*1000 + int(yCell)
	return key
}

func getSurroundingKeys(pos Vector2, cells *[]int) {
	for i := 0; i <= 8; i++ {
		var posCopy = pos
		if i == 0 {
			//default
		} else if i == 1 {
			//top
			posCopy.y -= CELL_SIZE
		} else if i == 2 {
			//right
			posCopy.x += CELL_SIZE
		} else if i == 3 {
			//bottom
			posCopy.y += CELL_SIZE
		} else if i == 4 {
			//left
			posCopy.x -= CELL_SIZE
		} else if i == 5 {
			//topright
			posCopy.x += CELL_SIZE
			posCopy.y += CELL_SIZE
		} else if i == 6 {
			//bottomright
			posCopy.x += CELL_SIZE
			posCopy.y -= CELL_SIZE
		} else if i == 7 {
			//bottomleft
			posCopy.x -= CELL_SIZE
			posCopy.y -= CELL_SIZE
		} else if i == 8 {
			//topleft
			posCopy.x -= CELL_SIZE
			posCopy.y += CELL_SIZE
		}
		*cells = append(*cells, getEntityKey(posCopy))
	}
}

func getEntitiesInKey(key int, entities *[]*Entity) {
	for _, e := range m[key] {
		*entities = append(*entities, e)
		e.rect.pos.x = 0
		e.rect.pos.y = 0
	}
}

// func getPlayersInKey(key int, players *[]*Player) bool {
// 	var foundAPlayer = false
// 	if len(m[key]) > 1 {
// 		for i := 0; i < len(m[key]); i++ {
// 			player, _ := getEntityObject(key, i)
// 			if player != nil {
// 				*players = append(*players, player)
// 				foundAPlayer = true
// 			}
// 		}
// 	}
// 	return foundAPlayer
// }

// func getNpcsInKey(key int, npcs *[]*Npc) bool {
// 	var foundAnNpc = false
// 	//if len(m[key]) > 1 {
// 	for i := 0; i < len(m[key]); i++ {
// 		_, npc := getEntityObject(key, i)
// 		if npc != nil {
// 			*npcs = append(*npcs, npc)
// 			foundAnNpc = true
// 		}
// 	}
// 	//}
// 	return foundAnNpc
// }

// func getEntitiesInAllKeysNear(pos Vector2, entities *[]*Entity) {
// 	for i := 0; i <= 8; i++ {
// 		var posCopy = pos
// 		if i == 0 {
// 			//default
// 		} else if i == 1 {
// 			//top
// 			posCopy.y -= CELL_SIZE
// 		} else if i == 2 {
// 			//right
// 			posCopy.x += CELL_SIZE
// 		} else if i == 3 {
// 			//bottom
// 			posCopy.y += CELL_SIZE
// 		} else if i == 4 {
// 			//left
// 			posCopy.x -= CELL_SIZE
// 		} else if i == 5 {
// 			//topright
// 			posCopy.x += CELL_SIZE
// 			posCopy.y += CELL_SIZE
// 		} else if i == 6 {
// 			//bottomright
// 			posCopy.x += CELL_SIZE
// 			posCopy.y -= CELL_SIZE
// 		} else if i == 7 {
// 			//bottomleft
// 			posCopy.x -= CELL_SIZE
// 			posCopy.y -= CELL_SIZE
// 		} else if i == 8 {
// 			//topleft
// 			posCopy.x -= CELL_SIZE
// 			posCopy.y += CELL_SIZE
// 		}
// 		var key = getEntityKey(posCopy)
// 		getEntitiesInKey()
// 	}
// }

// func getPlayersInKey(key int, players *[]*Player) bool {
// 	var foundAPlayer = false
// 	if len(m[key]) > 1 {
// 		for i := 0; i < len(m[key]); i++ {
// 			player, _ := getEntityObject(key, i)
// 			if player != nil {
// 				*players = append(*players, player)
// 				foundAPlayer = true
// 			}
// 		}
// 	}
// 	return foundAPlayer
// }

// func getPlayersInAllKeysNearPos(pos Vector2, players *[]*Player) {
// 	for i := 0; i <= 8; i++ {
// 		var posCopy = pos
// 		if i == 0 {
// 			//default
// 		} else if i == 1 {
// 			//top
// 			posCopy.y -= CELL_SIZE
// 		} else if i == 2 {
// 			//right
// 			posCopy.x += CELL_SIZE
// 		} else if i == 3 {
// 			//bottom
// 			posCopy.y += CELL_SIZE
// 		} else if i == 4 {
// 			//left
// 			posCopy.x -= CELL_SIZE
// 		} else if i == 5 {
// 			//topright
// 			posCopy.x += CELL_SIZE
// 			posCopy.y += CELL_SIZE
// 		} else if i == 6 {
// 			//bottomright
// 			posCopy.x += CELL_SIZE
// 			posCopy.y -= CELL_SIZE
// 		} else if i == 7 {
// 			//bottomleft
// 			posCopy.x -= CELL_SIZE
// 			posCopy.y -= CELL_SIZE
// 		} else if i == 8 {
// 			//topleft
// 			posCopy.x -= CELL_SIZE
// 			posCopy.y += CELL_SIZE
// 		}
// 		var key = getEntityKey(posCopy)
// 		getPlayersInKey(key, players)

// 		if pos.x == 0 {
// 			//fmt.Println("Key: ", key, "Length of Key: ", len(m[key]))
// 			//fmt.Println(found)
// 		}

// 	}
// }

// func getNpcsInAllKeysNearPos(pos Vector2, npcs *[]*Npc) {
// 	for i := 0; i <= 8; i++ {
// 		var posCopy = pos
// 		if i == 0 {
// 			//default
// 		} else if i == 1 {
// 			//top
// 			posCopy.y -= CELL_SIZE
// 		} else if i == 2 {
// 			//right
// 			posCopy.x += CELL_SIZE
// 		} else if i == 3 {
// 			//bottom
// 			posCopy.y += CELL_SIZE
// 		} else if i == 4 {
// 			//left
// 			posCopy.x -= CELL_SIZE
// 		} else if i == 5 {
// 			//topright
// 			posCopy.x += CELL_SIZE
// 			posCopy.y += CELL_SIZE
// 		} else if i == 6 {
// 			//bottomright
// 			posCopy.x += CELL_SIZE
// 			posCopy.y -= CELL_SIZE
// 		} else if i == 7 {
// 			//bottomleft
// 			posCopy.x -= CELL_SIZE
// 			posCopy.y -= CELL_SIZE
// 		} else if i == 8 {
// 			//topleft
// 			posCopy.x -= CELL_SIZE
// 			posCopy.y += CELL_SIZE
// 		}
// 		var key = getEntityKey(posCopy)
// 		getNpcsInKey(key, npcs)

// 		//if pos.x == 0 {
// 		fmt.Println("Key: ", key, "Length of Key: ", len(m[key]))
// 		//fmt.Println(found)
// 		//}

// 	}
// }

// func isAPlayerNear(pos Vector2, sightRange float64) (bool, *Player) {
// 	//create return variables
// 	var closestPlayer *Player
// 	var isAPlayerNear = false
// 	//other vars
// 	var players []*Player
// 	//var key = getEntityKey(pos)
// 	var closestDistance float64 = -1

// 	//make sure the cell holds more than just this npc
// 	// if len(m[key]) > 1 {
// 	// 	for i := 0; i < len(m[key]); i++ {
// 	// 		player, _ := getEntityObject(key, i)
// 	// 		if player != nil {
// 	// 			players = append(players, player)
// 	// 		}
// 	// 	}
// 	// }
// 	//getPlayersInKey(key, &players)
// 	getPlayersInAllKeysNearPos(pos, &players)

// 	//find closest player in array
// 	for _, p := range players {
// 		var d = distance(Vector2{x: p.entity.rect.pos.x, y: p.entity.rect.pos.y}, pos)
// 		if d <= sightRange {
// 			isAPlayerNear = true
// 		}
// 		if closestDistance == -1 || d < closestDistance {
// 			closestPlayer = p
// 			closestDistance = d
// 		}
// 	}

// 	//loop through players to see if
// 	return isAPlayerNear, closestPlayer
// }
