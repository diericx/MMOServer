package main

import (
	"fmt"
	"net"
	"ugorji/go/codec"
	"io"
	"time"
	"math"
)

const listenAddr = ":7777"

const baseAddr = "http://192.168.1.18:3000/api/v1/"

var serverConn *net.UDPConn

var (
	mh codec.MsgpackHandle
	r  io.Reader
	w  io.Writer
)

func sendMessage(msg []byte, conn *net.UDPConn, addr *net.UDPAddr) {
	_, err := conn.WriteToUDP(msg, addr)
	if err != nil {
		fmt.Printf("Couldn't send response %v", err)
	}
}

func getDataFromPlayer(player *Player, buf []byte, n int) {
	var res = &Update{}
	//decode data
	dec := codec.NewDecoder(r, &mh)
	dec = codec.NewDecoderBytes(buf[0:n], &mh)
	err := dec.Decode(res)

	//fmt.Println(res.Action)

	if err == nil {
		
		if res.Action == "update" {
			//if action is UPDATE
			player.lastUpdate = time.Now()

			player.id = res.ID

			player.shooting = res.Shooting

			player.xMovement = res.X
			player.yMovement = res.Y

			player.rect.rotation = res.Rotation

		
		} else if res.Action == "equip" {
			//if action is EQUIP

			//put item to equip data in variables
			var itemToEquip = player.inventory[res.Value] //H1
			//put item currently wearing into variabls
			var currentItem = ""

			if string(itemToEquip[0]) == "W" {
				currentItem = player.gear[2]
				//replace equiped item
				player.gear[2] = itemToEquip
			} else if string(itemToEquip[0]) == "H" {
				currentItem = player.gear[0] //H2
				//replace equiped item
				player.gear[0] = itemToEquip
			}

			//remove equipped item from inventory
			removeItemFromInventoryViaIndex(&player.inventory, res.Value)
			//add item that was replaced
			addItemToInventory(&player.inventory, res.Value, currentItem)

		} else if res.Action == "drop" {
			println(string(res.Value))
			removeItemFromInventoryViaIndex(&player.inventory, res.Value)
		} else if res.Action == "shoot" {
			player.health = player.health
		} else if res.Action == "upgradeHealthCap" {
			if player.skillPoints >= 1 && player.healthCap < HEALTHCAP_CAP {
				player.healthCap += 1
				player.skillPoints -= 1
			}
		} else if res.Action == "upgradeHealthRegen" {
			if player.skillPoints >= 1 && player.healthRegen < HEALTH_REGEN_CAP {
				player.healthRegen += 1
				player.skillPoints -= 1
			}
		} else if res.Action == "upgradeSpeed" {
			if player.skillPoints >= 1 && player.speed < MOVE_SPEED_CAP {
				player.speed += 1
				player.skillPoints -= 1
			}
		} else if res.Action == "upgradeShieldCap" {
			if player.skillPoints >= 1 && player.shieldCap < SHIELD_CAP {
				player.shieldCap += 1
				player.skillPoints -= 1
			}
		} else if res.Action == "upgradeShieldRegen" {
			if player.skillPoints >= 1 && player.shieldRegen < SHIELD_REGEN_CAP {
				player.shieldRegen += 1
				player.skillPoints -= 1
			}
		} else if res.Action == "upgradeEnergyCap" {
			if player.scraps >= 100 && player.energyCap < ENERGY_CAP_CAP {
				player.energyCap += 1
				player.skillPoints -= 1
			}
		} else if res.Action == "upgradeEnergyRegen" {
			if player.skillPoints >= 1 && player.energyRegen < ENERGY_REGEN_CAP {
				player.energyRegen += 1
				player.skillPoints -= 1
			}
		} else if res.Action == "upgradeShieldRegen" {
			if player.skillPoints >= 1 && player.shieldRegen < SHIELD_REGEN_CAP {
				player.shieldRegen += 1
				player.skillPoints -= 1
			}
		} else if res.Action == "upgradeDamage" {
			if player.skillPoints >= 1 && player.damage < DAMAGE_CAP {
				player.damage += 1
				player.skillPoints -= 1
			}
		} else if res.Action == "upgradeFireRate" {
			if player.skillPoints >= 1 && player.fireRate < FIRE_RATE_CAP {
				player.fireRate += 1
				player.skillPoints -= 1
			}
		} else if res.Action == "jump" {
			if player.scraps >= 200 {
				player.scraps -= 200
				player.rect.x = res.X
				player.rect.y = res.Y
			}
		}
	}
	//}

}

func sendData() {

	var speedMod = 30
	var sentWelcomeMsg = false
	for {

		//var bulletPackets []*BulletUpdate

		for _, player := range players {

			bulletIDs := make([]int, 0)
			bulletXs := make([]float64, 0)
			bulletYs := make([]float64, 0)
			bulletRots := make([]int, 0)

			otherPlayerIDs := make([]string, 0)
			otherPlayerXs := make([]float64, 0)
			otherPlayerYs := make([]float64, 0)
			otherPlayerRots := make([]int, 0)
			otherPlayerHlths := make([]float64, 0)
			otherPlayerGearSets := make([][]string, 0)

			npcIDs := make([]int, 0)
			npcTypes := make([]int, 0)
			npcXs := make([]float64, 0)
			npcYs := make([]float64, 0)
			npcHlths := make([]float64, 0)

			targetPlayers_rotations := make([]float64, 0)

			if player.id != "" {

				//put all bullets into one array that are CLOSE TO THE PLAYER
				//WARNING: MAY CAUSE LAG
				for _, bullet := range bullets {

					var dist = math.Sqrt(math.Pow(bullet.rect.x-player.rect.x, 2) + math.Pow(bullet.rect.y-player.rect.y, 2))
					if dist <= PLAYER_LOAD_DIST {
						bulletIDs = append(bulletIDs, bullet.ID)
						bulletXs = append(bulletXs, bullet.rect.x)
						bulletYs = append(bulletYs, bullet.rect.y)
						bulletRots = append(bulletRots, bullet.rect.rotation)
					}

				}

				//get all data from other players
				//WARNING: MAY CAUSE LAG
				for _, otherPlayer := range players {
					if player.id == "1" && player != otherPlayer {
						//add other player's rotations to list
						var rotation = getAngleBetween2Vectors(Vector2{x: player.rect.x, y: player.rect.y}, Vector2{x: otherPlayer.rect.x, y: otherPlayer.rect.y})
						targetPlayers_rotations = append(targetPlayers_rotations, rotation)
					}

					var dist = math.Sqrt(math.Pow(otherPlayer.rect.x-player.rect.x, 2) + math.Pow(otherPlayer.rect.y-player.rect.y, 2))
					if dist <= PLAYER_LOAD_DIST && player != otherPlayer {
						otherPlayerIDs = append(otherPlayerIDs, otherPlayer.id)
						otherPlayerXs = append(otherPlayerXs, otherPlayer.rect.x)
						otherPlayerYs = append(otherPlayerYs, otherPlayer.rect.y)
						otherPlayerRots = append(otherPlayerRots, otherPlayer.rect.rotation)
						otherPlayerHlths = append(otherPlayerHlths, otherPlayer.health)

						gearSet := []string{otherPlayer.gear[0], otherPlayer.gear[1], otherPlayer.gear[2], otherPlayer.gear[3]}
						otherPlayerGearSets = append(otherPlayerGearSets, gearSet)
					}

				}

				//get all data from NPCs
				//WARNING: MAY CAUSE LAG
				for _, npc := range npcs {
					var dist = math.Sqrt(math.Pow(npc.rect.x-player.rect.x, 2) + math.Pow(npc.rect.y-player.rect.y, 2))
					if dist <= PLAYER_LOAD_DIST {
						npcIDs = append(npcIDs, npc.id)
						npcTypes = append(npcTypes, npc.npcType)
						npcXs = append(npcXs, npc.rect.x)
						npcYs = append(npcYs, npc.rect.y)
						npcHlths = append(npcHlths, npc.health)
					}
				}

				//create new gear obj for the other players current gear set
				// otherPlayersGear := gearSet{
				//     cockpit: otherPlayer.gear.cockpit,
				//     lasers: otherPlayer.gear.lasers,
				//     wings: otherPlayer.gear.wings,
				//     jets: otherPlayer.gear.jets,
				// }

				//new table that has multiple updates

				//create update packet
				res1D := &Update{
					Action:                  "playerUpdate",
					ID:                      player.id,
					Level:                   player.level,
					XP:                      player.xp,
					SkillPoints:             player.skillPoints,
					Infamy:                  player.infamy,
					FireRate:                player.fireRate,
					Health:                  player.health,
					HealthCap:               player.healthCap,
					HealthRegen:             player.healthRegen,
					Energy:                  player.energy,
					EnergyCap:               player.energyCap,
					EnergyRegen:             player.energyRegen,
					Shield:                  player.shield,
					ShieldCap:               player.shieldCap,
					ShieldRegen:             player.shieldRegen,
					Speed:                   player.speed,
					Damage:                  player.damage,
					Scraps:                  player.scraps,
					X:                       player.rect.x,
					Y:                       player.rect.y,
					Rotation:                player.rect.rotation,
					Gear:                    player.gear,
					Inventory:               player.inventory,
					IsNPC:                   false,
					TargetNPC_rotation:      player.targetNPC_rotation,
					TargetPlayers_rotations: targetPlayers_rotations,
					OtherPlayerIDs:          otherPlayerIDs,
					OtherPlayerXs:           otherPlayerXs,
					OtherPlayerYs:           otherPlayerYs,
					OtherPlayerRots:         otherPlayerRots,
					OtherPlayerHlths:        otherPlayerHlths,
					OtherPlayerGearSets:     otherPlayerGearSets,
					BulletIDs:               bulletIDs,
					BulletXs:                bulletXs,
					BulletYs:                bulletYs,
					BulletRots:              bulletRots,
					NpcIDs:                  npcIDs,
					NpcTypes:                npcTypes,
					NpcXs:                   npcXs,
					NpcYs:                   npcYs,
					NpcHlths:                npcHlths,
				}

				var newByteArray []byte
				enc := codec.NewEncoder(w, &mh)
				enc = codec.NewEncoderBytes(&newByteArray, &mh)
				enc.Encode(res1D)

				var stringMessage = string(newByteArray)
				//fmt.Println(stringMessage)
				//create header
				//var header = intToBinaryString(len(stringMessage))
				//fmt.Println(len(stringMessage))
				//format message message
				stringMessage = stringMessage
				//print message
				if player.id == "1" {
					//fmt.Printf(stringMessage + "\n")
				}
				// fmt.Printf("Header: %v\n", header);
				// fmt.Printf("Packet Length: %v\n", len(newByteArray));

				//send message
				//go fmt.Fprint(player.rwc, stringMessage)
				sendMessage([]byte(stringMessage), serverConn, &player.addr)
			} else {
				if sentWelcomeMsg == false {
					sentWelcomeMsg = true
					res1F := &Message{
						Action: "message",
						Data:   "connected!",
					}

					var newByteArray []byte
					enc := codec.NewEncoder(w, &mh)
					enc = codec.NewEncoderBytes(&newByteArray, &mh)
					enc.Encode(res1F)

					var stringMessage = string(newByteArray)
					//---create header---
					//var header = intToBinaryString(len(stringMessage))
					//---add header---
					stringMessage = stringMessage

					stringMessage = "<?xml version=\"1.0\"?>\n<cross-domain-policy>\n<allow-access-from domain=\"*\" to-ports=\"7770-7780\"/>\n</cross-domain-policy>\n"
					//send message
					sendMessage([]byte(stringMessage), serverConn, &player.addr)
					//go fmt.Fprint(player.rwc, stringMessage)
				}
			}
		}

		//remove bullets AFTER sending data
		clearBulletRemoveList(&bullets)

		time.Sleep(1 * (time.Second / time.Duration(speedMod)))
	}

}
