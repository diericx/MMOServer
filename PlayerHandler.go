package main

import (
	"net"
	"time"
	"math"
)

var players []*Player

func instantiatePlayer(addr *net.UDPAddr) *Player {
	// setup new player and its stats
	newPlayer := new(Player)
	newPlayer.addr = *addr
	newPlayer.id = ""
	newPlayer.infamy = 0
	newPlayer.health = 100
	newPlayer.healthCap = 0
	newPlayer.healthRegen = 0
	newPlayer.energy = 50
	newPlayer.energyCap = 0
	newPlayer.energyRegen = 0
	newPlayer.shield = 10
	newPlayer.shieldCap = 0
	newPlayer.shieldRegen = 0 //per tenth of a second
	newPlayer.fireRate = 0
	newPlayer.fireRateCooldown = 0
	newPlayer.damage = 0
	newPlayer.speed = 0
	newPlayer.scraps = 0
	newPlayer.weaponCooldownCap = 0.5
	newPlayer.weaponCooldown = 0
	newPlayer.weaponBulletCount = 1

	newPlayer.targetNPC = *npcs[0]
	newPlayer.targetNPC_rotation = 0

	newPlayer.rect = createRect(0, 0, 3, 3)

	newPlayer.gear = []string{
		"H1",
		"L1",
		"W1",
		"T1"}

	newPlayer.inventory = make([]string, 8)

	players = append(players, newPlayer)
	println("Player Joined!")

	return newPlayer
	//go getDataFromPlayer(newPlayer)

}


func updatePlayerStats() {
	for {
		for _, player := range players {

			//check if player is not respinding
			if time.Since(player.lastUpdate).Seconds() >= 0.5 {
				//remove player
				removePlayer(player)
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
				var healthCap = float64(BASE_HEALTH_CAP_VALUE) + (10 * float64(player.healthCap)) + hullHealthCapAttr
				var healthRegen = BASE_HEALTH_REGEN_VALUE + (0.002 * float64(player.healthRegen))
				if player.health < healthCap {
					player.health += healthRegen
					if player.health > healthCap {
						player.health = healthCap
					}
				}

				//update shield stat
				var shieldCap = BASE_SHIELD_CAP_VALUE + (10 * float64(player.shieldCap))
				var shieldRegen = BASE_SHIELD_REGEN_VALUE + (0.01 * float64(player.shieldRegen))
				if player.shield < shieldCap {
					player.shield += shieldRegen
					if player.shield > shieldCap {
						player.shield = shieldCap
					}
				}

				//update energy stat
				var energyCap = BASE_ENERGY_CAP_VALUE + (10 * float64(player.energyCap))
				var energyRegen = BASE_ENERGY_REGEN_VALUE + (0.01 * float64(player.energyRegen))
				if player.energy < energyCap {
					player.energy += energyRegen
					if player.energy > energyCap {
						player.energy = energyCap
					}
				}

				//shoot
				var fireRate = BASE_FIRE_RATE_VALUE - (10 * player.fireRate)
				if player.shooting {
					player.fireRateCooldown -= 1

					if player.fireRateCooldown <= 0 {
						player.fireRateCooldown = fireRate

						// spawn new bullet
						handleShot("singleShot", player, player.damage, player.rect, &bullets)
					}
				} else {
					player.fireRateCooldown = 0
				}

				//update player target rotation
				//if player.targetNPC != nil {
				if (distance(Vector2{x: player.rect.x, y: player.rect.y}, player.targetNPC.origin) >= 20) {
					var rotation = getAngleBetween2Vectors(Vector2{x: player.rect.x, y: player.rect.y}, player.targetNPC.origin)
					player.targetNPC_rotation = rotation
				} else {
					player.targetNPC_rotation = -1
				}
				//}

				//if the npc is still alive
				if player.targetNPC.alive == false {
					//if the npc is far enough away
					player.targetNPC = *npcs[0]
				}
			}
		}
		time.Sleep((time.Second / time.Duration(1000)))
	}
}

func movePlayers() {
	for {
		for _, player := range players {

			wingSpeedAttr := getItemAttribute(player.gear[2], "speed")

			var speed = BASE_SPEED_VALUE + (0.5 * float64(player.speed+int(wingSpeedAttr)))
			player.rect.x = player.rect.x + (player.xMovement * (speed / 100))
			player.rect.y = player.rect.y + (player.yMovement * (speed / 100))

			if player.rect.x >= ARENA_SIZE {
				player.rect.x = ARENA_SIZE
			} else if player.rect.x <= -ARENA_SIZE {
				player.rect.x = -ARENA_SIZE
			}

			if player.rect.y >= ARENA_SIZE {
				player.rect.y = ARENA_SIZE
			} else if player.rect.y <= -ARENA_SIZE {
				player.rect.y = -ARENA_SIZE
			}
			//player.rect.rotation = player.Rotation
		}
		time.Sleep((time.Second / time.Duration(300)))
	}

}

func playerAlreadyExists(addr *net.UDPAddr) *Player {
	var foundPlayer *Player
	for _, player := range players {
		var playerAddress = &player.addr
		var inputAddress = addr
		if playerAddress.String() == inputAddress.String() {
			foundPlayer = player
		}
	}
	return foundPlayer
}

func findPlayerIndex(p *Player) int {
	var i = 0
	var foundIndex = -1
	for _, player := range players {
		if p == player {
			foundIndex = i
		}
		i++
	}
	return foundIndex
}

func findPlayerIndexByID(pID string) int {
	var i = 0
	var foundIndex = -1
	for _, player := range players {
		if player.id == pID {
			foundIndex = i
		}
		i++
	}
	return foundIndex
}

func removePlayerFromList(p *Player) {
	var foundIndex = findPlayerIndex(p)

	if foundIndex != -1 {
		players = append(players[:foundIndex], players[foundIndex+1:]...)
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

//Check if a player is near a location//
//!!VERY EXPENSIVE!!//
func isAPlayerNear(pos Vector2, sightRange float64) bool {
	var value = false
	for _, player := range players {
		var dist = math.Sqrt(math.Pow(player.rect.x-pos.x, 2) + math.Pow(player.rect.y-pos.y, 2))
		if dist <= sightRange {
			value = true
			break
		}
	}
	return value
}