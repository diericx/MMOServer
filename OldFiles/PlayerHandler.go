package main

import (
	// "fmt"
	"math"
	"net"
	"time"
)

type Player struct {
	entity                *Entity
	addr                  net.UDPAddr
	level                 int
	xp                    int
	skillPoints           int
	movement              Vector2
	scraps                int32
	targetEntity          *Entity
	targetEntity_rotation float64
	gear                  []string
	inventory             []string
	infamy                int
}

func instantiatePlayer(addr *net.UDPAddr) *Player {
	// setup new player and its stats
	var newPlayer Player
	newPlayer.addr = *addr
	newPlayer.entity = instantiateEntity()

	var playerEntity = *newPlayer.entity
	//Fill in all values
	playerEntity.id = ""
	playerEntity.health = 100
	playerEntity.energy = 50
	playerEntity.shield = 10
	playerEntity.weaponCooldownCap = 0.5
	playerEntity.weaponBulletCount = 1

	//playerEntity.targetEntity = npcs[0].entity

	playerEntity.rect = createRect(0, 0, 3, 3)

	playerEntity.key = -1 //getEntityKey(Vector2{x: newPlayer.rect.x, y: newPlayer.rect.y})

	newPlayer.gear = []string{
		"H1",
		"L1",
		"W1",
		"T1"}

	newPlayer.inventory = make([]string, 8)

	*newPlayer.entity = playerEntity

	players = append(players, newPlayer)
	println("Player Joined!")

	return &newPlayer
	//go getDataFromPlayer(newPlayer)

}

func updatePlayers() {
	for {
		// fmt.Println(len(m))
		for _, player := range players {
			var entity = player.entity
			//move the player
			movePlayer(&player)
			//update player's hash data
			updateEntityCellData(player.entity)

			//check if player is not respinding
			if time.Since(player.entity.lastUpdate).Seconds() >= 0.5 {
				//remove player
				removePlayer(&player)
				break
			} else {
				//------update player data-------
				//update XP and Level
				var currentXPCap float64 = float64(BASE_XP) * (math.Pow(float64(player.level), float64(LEVEL_XP_FACTOR)))
				var currentXPCapRounded = int(currentXPCap) //round number
				if player.xp >= currentXPCapRounded {
					var diff = player.xp - currentXPCapRounded
					player.level += 1
					player.xp = diff
					player.skillPoints += 1
				}

				//update health stat
				hullHealthCapAttr := getItemAttribute(player.gear[0], "healthCap")
				var healthCap = float64(BASE_HEALTH_CAP_VALUE) + (10 * float64(player.entity.healthCap)) + hullHealthCapAttr
				var healthRegen = BASE_HEALTH_REGEN_VALUE + (0.002 * float64(player.entity.healthRegen))
				if entity.health < healthCap {
					entity.health += healthRegen
					if entity.health > healthCap {
						entity.health = healthCap
					}
				}

				//update shield stat
				var shieldCap = BASE_SHIELD_CAP_VALUE + (10 * float64(entity.shieldCap))
				var shieldRegen = BASE_SHIELD_REGEN_VALUE + (0.01 * float64(entity.shieldRegen))
				if entity.shield < shieldCap {
					entity.shield += shieldRegen
					if entity.shield > shieldCap {
						entity.shield = shieldCap
					}
				}

				//update energy stat
				var energyCap = BASE_ENERGY_CAP_VALUE + (10 * float64(entity.energyCap))
				var energyRegen = BASE_ENERGY_REGEN_VALUE + (0.01 * float64(entity.energyRegen))
				if entity.energy < energyCap {
					entity.energy += energyRegen
					if entity.energy > energyCap {
						entity.energy = energyCap
					}
				}

				//shoot
				var fireRate = BASE_FIRE_RATE_VALUE - (10 * entity.fireRate)

				if entity.fireRateCooldown < 0 {
					entity.fireRateCooldown = 0
					//player.fireRateCooldown = fireRate
				} else {
					entity.fireRateCooldown -= 1
				}

				if entity.shooting {

					if entity.fireRateCooldown == 0 {
						entity.fireRateCooldown = fireRate

						// spawn new bullet
						//handleShot("singleShot", player, entity.damage, entity.rect, &bullets)
					}
				}

				//update player target rotation
				if player.targetEntity != nil {
					if (distance(Vector2{x: entity.rect.pos.x, y: entity.rect.pos.y}, player.targetEntity.origin) >= 20) {
						var rotation = getAngleBetween2Vectors(Vector2{x: entity.rect.pos.x, y: entity.rect.pos.y}, player.targetEntity.origin)
						player.targetEntity_rotation = rotation
					} else {
						player.targetEntity_rotation = -1
					}
					//if the npc is still alive
					if player.targetEntity.alive == false {
						//if the npc is far enough away
						player.targetEntity = npcs[0].entity
					}
				}

				//Update NPCs near the player
				// var npcsNear []*Npc
				// getNpcsInAllKeysNearPos(Vector2{x: entity.rect.pos.x, y: entity.rect.pos.y}, &npcsNear)
				// for _, npc := range npcsNear {
				// 	if npc.entity.entityType == 3 {
				// 		//moveNPC3(npc)
				// 	}
				// }
				//fmt.Println(len(npcsNear))
			}
		}
		time.Sleep((time.Second / time.Duration(1000)))
	}
}

func removePlayer(player *Player) {
	for i, otherPlayer := range players {
		if otherPlayer.addr.String() == player.addr.String() {
			players = append(players[:i], players[i+1:]...)
			return
		}
	}
}

func movePlayer(player *Player) {
	//-----------------
	// fmt.Println(player.movement)
	var playerCp = *player
	var entity = playerCp.entity
	//Move Player
	wingSpeedAttr := 1 // getItemAttribute(player.gear[2], "speed")

	var speed = BASE_SPEED_VALUE + (0.5 * float64(entity.speed+int(wingSpeedAttr)))
	entity.rect.pos.x = entity.rect.pos.x + (player.movement.x*(speed/100))/PLAYER_SPEED_DAMPEN
	entity.rect.pos.y = entity.rect.pos.y + (player.movement.y*(speed/100))/PLAYER_SPEED_DAMPEN

	if entity.rect.pos.x >= ARENA_SIZE {
		entity.rect.pos.x = ARENA_SIZE
	} else if entity.rect.pos.x <= -ARENA_SIZE {
		entity.rect.pos.x = -ARENA_SIZE
	}

	if entity.rect.pos.y >= ARENA_SIZE {
		entity.rect.pos.y = ARENA_SIZE
	} else if entity.rect.pos.y <= -ARENA_SIZE {
		entity.rect.pos.y = -ARENA_SIZE
	}
	//----------------
	*player = playerCp
}

func playerAlreadyExists(addr *net.UDPAddr) *Player {
	foundPlayer := -1
	i := 0
	for _, player := range players {
		var playerAddress = player.addr
		var inputAddress = addr
		if playerAddress.String() == inputAddress.String() {
			foundPlayer = i
		}
		i++
	}
	if (foundPlayer == -1) {
		return nil 
	} else {
		return &players[foundPlayer]
	}
}
