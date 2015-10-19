package main

import(
	"math/rand"
	"math"
	"time"
)

var npcs []*Npc

func spawnNPCs() {
	//create test NPC type 2
	newNPC := new(Npc)
	newNPC.id = rand.Intn(10000)
	newNPC.npcType = 2
	newNPC.damage = 2
	newNPC.shotTime = 1000
	newNPC.shotCooldown = 1000
	newNPC.health = 50
	newNPC.rect.rotation = 0
	newNPC.bulletRange = 15
	var x float64 = 0
	var y float64 = 0
	newNPC.rect = createRect(x, y, 3, 3)

	npcs = append(npcs, newNPC)
}

func spawnNPC(npc_type int, damage int, shotTime int, shotCooldown int, health float64, bullet_range int, min_dist float64, max_dist float64) {
	newNPC := new(Npc)
	newNPC.id = rand.Intn(10000)
	newNPC.npcType = npc_type
	newNPC.shotTime = shotTime
	newNPC.shotCooldown = shotCooldown
	newNPC.health = health
	newNPC.rect.rotation = 0
	newNPC.bulletRange = bullet_range
	var randAngle = rand.Float64() * 360
	var randRad = randAngle * math.Pi / 180
	var randRadius = randFloatInRange(min_dist, max_dist)
	var x = math.Cos(randRad) * randRadius
	var y = math.Sin(randRad) * randRadius //( rand.Float64() * ( METEOR_MAX_DIST - METEOR_MIN_DIST ) ) + METEOR_MIN_DIST
	newNPC.rect = createRect(x, y, 3, 3)
	newNPC.origin = Vector2{x: newNPC.rect.x, y: newNPC.rect.y}
	npcs = append(npcs, newNPC)
}

func removeNpcFromList(n *Npc) {
	var i = 0
	var foundIndex = -1
	for _, npc := range npcs {
		if n == npc {
			foundIndex = i
		}
		i++
	}
	if foundIndex != -1 {
		npcs = append(npcs[:foundIndex], npcs[foundIndex+1:]...)
	}
	var copy = *n
	copy.alive = false
	*n = copy
}

func moveEnemyRandomly(rect *Rectangle, targetPos Vector2) {
	var rectCopy = *rect
	if rectCopy.x == targetPos.x && rectCopy.y == targetPos.y {
		//move randomly if already on position
		var randMultiplierX = rand.Float64() * 8
		var randMultiplierY = rand.Float64() * 8
		rectCopy.x = targetPos.x + NPC_2_MAX_MOVE_DIST*randMultiplierX
		rectCopy.y = targetPos.y + NPC_2_MAX_MOVE_DIST*randMultiplierY
	} else {
		rectCopy.x = targetPos.x
		rectCopy.y = targetPos.y
	}
	*rect = rectCopy
}


func updateNPCs() {
	for {
		//count of meteors
		var meteor_count = 0
		var npc_2_count = 0

		for _, npc := range npcs {

			//MOVE NPCs
			if npc.npcType == 1 {
				//count amount of meteors
				meteor_count += 1
			} else if npc.npcType == 2 {
				npc_2_count += 1
				//if (npc_2_count == 1) {
				moveEnemyRandomly(&npc.rect, npc.origin)
				//}
				// npc.rect.x = npc.rect.x + rand.Float64() - 0.5
				// npc.rect.y = npc.rect.y + rand.Float64() - 0.5
			}

			if npc.shotTime != -1 {
				if npc.shotCooldown > 0 {
					npc.shotCooldown -= 1
				} else if npc.shotCooldown <= 0 {
					npc.shotCooldown = npc.shotTime
					if (isAPlayerNear(Vector2{x: npc.rect.x, y: npc.rect.y}, NPC_2_SIGHT)) {
						handleShot("radialShotgunShot", npc, 0, npc.rect, &bullets)
					}
				}
			}
		}

		if meteor_count < METEOR_MAX_AMMOUNT {
			for i := 0; i < (METEOR_MAX_AMMOUNT - meteor_count); i++ {
				spawnNPC(1, 0, -1, -1, 50, 15, METEOR_MIN_DIST, METEOR_MAX_DIST)
				// newNPC := new(Npc)
				// newNPC.id = rand.Intn(10000)
				// newNPC.npcType = 1
				// newNPC.shotTime = -1
				// newNPC.shotCooldown = -1
				// newNPC.health = 50
				// newNPC.rect.rotation = 0
				// newNPC.bulletRange = 10k
				// var randAngle = rand.Float64() * 360
				// var randRad = randAngle * math.Pi / 180
				// var randRadius = randFloatInRange(METEOR_MIN_DIST, METEOR_MAX_DIST)
				// var x = math.Cos(randRad) * randRadius
				// var y = math.Sin(randRad) * randRadius//( rand.Float64() * ( METEOR_MAX_DIST - METEOR_MIN_DIST ) ) + METEOR_MIN_DIST
				// newNPC.rect = createRect(x, y, 3, 3)
				// npcs = append(npcs, newNPC)
			}
		}
		if npc_2_count < NPC_2_MAX_AMMOUNT {
			for i := 0; i < (NPC_2_MAX_AMMOUNT - npc_2_count); i++ {
				//spawnNPC(2, 2, 1000, 1000, 50, 15, METEOR_MIN_DIST, METEOR_MAX_DIST)
				spawnNPC(2, 2, 1000, 1000, 50, 15, NPC_2_MIN_DIST, NPC_2_MAX_DIST)
			}
		}

		time.Sleep((time.Second / time.Duration(1000)))
	}
}