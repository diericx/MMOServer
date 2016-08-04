package main

import (
	"math"
	"math/rand"
	"strconv"
	"time"
	// "fmt"
)

type Npc struct {
	entity *Entity
}

var npcCounts = make(map[int]int)

func spawnNPC(npc_type int, damage int, shotTime float64, shotCooldown float64, health float64, bullet_range int, min_dist float64, max_dist float64) {
	//create npc var
	var newNpc Npc
	newNpc.entity = instantiateEntity()
	//create npc's entity object
	var npcEntity = *newNpc.entity
	//fill variables
	npcEntity.id = strconv.Itoa(rand.Intn(10000))
	npcEntity.entityType = npc_type
	npcEntity.weaponCooldownCap = shotTime
	npcEntity.weaponCooldown = shotCooldown
	npcEntity.health = health
	npcEntity.rect.rotation = 0
	npcEntity.fireRange = bullet_range
	var randAngle = rand.Float64() * 360
	var randRad = randAngle * math.Pi / 180
	var randRadius = randFloatInRange(min_dist, max_dist)
	var x = math.Cos(randRad) * randRadius
	var y = math.Sin(randRad) * randRadius //( rand.Float64() * ( METEOR_MAX_DIST - METEOR_MIN_DIST ) ) + METEOR_MIN_DIST
	npcEntity.rect = createRect(x, y, 3, 3)
	npcEntity.key = -1
	npcEntity.origin = Vector2{x: npcEntity.rect.pos.x, y: npcEntity.rect.pos.y}
	//npcEntity.move = moveNPC3

	//Keep track
	npcCounts[npc_type] += 1

	//update it
	npcEntity.update()

	*newNpc.entity = npcEntity

	npcs = append(npcs, newNpc)
}

func spawnNPCs() {
	//instantiate Test NPC
	var newNpc Npc
	newNpc.entity = instantiateEntity()

	var npcEntity = *newNpc.entity

	npcEntity.id = "-1"
	npcEntity.entityType = 3
	npcEntity.damage = 2
	npcEntity.fireRateCooldown = 1000
	npcEntity.fireRate = 1000
	npcEntity.health = 50
	npcEntity.rect.rotation = 0
	npcEntity.fireRange = 15
	var x float64 = 0
	var y float64 = 0
	npcEntity.rect = createRect(x, y, 3, 3)

	npcEntity.key = -1

	*newNpc.entity = npcEntity

	npcs = append(npcs, newNpc)

	//Spawn the rest
	for {
		//keep track of already spawned npcs

		//update npc count
		//npcCounts[npc.entity.entityType] += 1

		// playerNear, _ := isAPlayerNear(npc.origin, 50)
		// //findClosestPlayer(npc)
		// //fmt.Println(closestPlayer)
		// if playerNear == true {

		// 	//MOVE NPCs
		// 	if npc.npcType == 1 {
		// 	} else if npc.npcType == 2 {
		// 		moveEnemyRandomly(&npc.rect, npc.origin)
		// 	} else if npc.npcType == 3 {
		// 		// var newOriginX = smoothTranslate(npc.origin.x, closestPlayer.rect.x, NPC_3_MOVE_SPEED)
		// 		// var newOriginY = smoothTranslate(npc.origin.y, closestPlayer.rect.y, NPC_3_MOVE_SPEED)
		// 		// npc.origin = Vector2{x:newOriginX, y: newOriginY}
		// 		moveEnemyRadially(&npc.rect, npc.origin, 8, 0.5, step)
		// 	}

		// if npc.entity.shotTime != -1 {
		// 	if npc.entity.fireRateCooldown > 0 {
		// 		npc.entity.fireRateCooldown -= 1
		// 	} else if npc.entity.fireRateCooldown <= 0 {
		// 		npc.entity.fireRateCooldown = npc.entity.fireRate

		// 		//if a player is near this npc, shoot and update data
		// 		if npc.npcType == 2 {
		// 			handleShot("radialShotgunShot", npc, 0, npc.rect, &bullets)
		// 		} else if npc.npcType == 3 {
		// 			handleShot("radialShotgunShot", npc, 0, npc.rect, &bullets)
		// 		}
		// 	}
		// }
		// }

		if npcCounts[1] < METEOR_MAX_AMMOUNT {
			for i := 0; i < (METEOR_MAX_AMMOUNT - npcCounts[1]); i++ {
				spawnNPC(1, 0, -1, -1, 50, 15, METEOR_MIN_DIST, METEOR_MAX_DIST)
			}
		}
		if npcCounts[2] < NPC_2_MAX_AMMOUNT {
			for i := 0; i < (NPC_2_MAX_AMMOUNT - npcCounts[2]); i++ {
				spawnNPC(2, 2, 1000, 1000, 50, 15, NPC_2_MIN_DIST, NPC_2_MAX_DIST)
			}
		}
		if npcCounts[3] < NPC_3_MAX_AMMOUNT {
			for i := 0; i < (NPC_3_MAX_AMMOUNT - npcCounts[3]); i++ {
				spawnNPC(3, 2, 1000, 1000, 50, 15, NPC_3_MIN_DIST, NPC_3_MAX_DIST)
			}
		}

		time.Sleep((time.Second / time.Duration(1000)))
	}
}

func updateNPCs() {
	for {
		//get all the cells
		var cells []int
		var entities []*Entity
		for _, player := range players { 
			getSurroundingKeys(player.entity.rect.pos, &cells)
		}

		for _, cell := range cells {
			getEntitiesInKey(cell, &entities)
		}

		// fmt.Println("Entiteis: ", len(entities))

		// for _, e := range entities {
			//e.update()
		// }

		//wait
		time.Sleep((time.Second / time.Duration(1000)))
	}
}
