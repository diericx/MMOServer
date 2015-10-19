package main

import (
	//"fmt"
	// "encoding/json"
	// "io/ioutil"
	"math/rand"
	"math"
	"time"
)

var bullets []*Bullet

func handleShot(shotType string, shooter interface{}, damage int, rect Rectangle, bullets *[]*Bullet) {

	if shotType == "singleShot" {
		fireSingleShot(shooter, bullets)
	} else if shotType == "radialShotgunShot" {
		fireRadialShotgunShot(shooter, bullets)
	}

}

func fireSingleShot(shooter interface{}, bullets *[]*Bullet) {
	var newBullet = fireBullet(shooter)
	*bullets = append(*bullets, newBullet)
}

func fireRadialShotgunShot(shooter interface{}, bullets *[]*Bullet) {
	for i := 0; i < 8; i++ {
		var newBullet = fireBullet(shooter)
		newBullet.rect.rotation = i*45
		*bullets = append(*bullets, newBullet)
	}
}

func fireBullet(shooter interface{}) *Bullet {
	//instantiate bullet
	newBullet := new(Bullet)

	//check if bullet.shooter is a player
	if p, ok := shooter.(*Player); ok {
		//player
		newBullet.shooter = p
		updateBulletAttributes(newBullet, p.damage, p.rect, p.rect.rotation, getItemAttribute(p.gear[1], "range"))
	} else if npc, ok := shooter.(*Npc); ok {
		//npc
		newBullet.shooter = npc
		updateBulletAttributes(newBullet, npc.damage, npc.rect, npc.rect.rotation, float64(npc.bulletRange) )
	}

	return newBullet
}

func updateBulletAttributes(newBullet *Bullet, damage int, rect Rectangle, rotation int, bulletRange float64) {
	newBullet.ID = rand.Intn(1000)
	newBullet.damage = damage
	newBullet.origin = Vector2{rect.x, rect.y}
	newBullet.rect = createRect(rect.x, rect.y, 0.17, 0.5)
	newBullet.rect.rotation = rotation
	newBullet.bulletRange = bulletRange
}


func moveBullets() {
	for {

		for _, bullet := range bullets {
			bulletPos := Vector2{x: bullet.rect.x, y: bullet.rect.y}
			if distance(bullet.origin, bulletPos) > bullet.bulletRange {
				//removeBulletFromList(&bullets, bullet)
				addBulletToRemoveList(bullet)
			}
		}

		for _, bullet := range bullets {
			var bulletRadians float64 = (float64(bullet.rect.rotation+90) / 180.0) * 3.14159
			bullet.rect.x = bullet.rect.x + (15 * 0.116 * math.Cos(bulletRadians))
			bullet.rect.y = bullet.rect.y + (15 * 0.116 * math.Sin(bulletRadians))
		}

		for _, bullet := range bullets {
			var bulletRemoved = false
			var bulletShooterP *Player
			var bulletShooterNPC *Npc

			//check if bullet.shooter is a player
			if p, ok := bullet.shooter.(*Player); ok {
				bulletShooterP = p
			} else {
				/* not player */
			}

			//check if bullet.shooter is an NPC
			if npc, ok := bullet.shooter.(*Npc); ok {
				bulletShooterNPC = npc
			} else {
				/* not player */
			}

			// Checkl bullets for collision with players
			for _, player := range players {
				if compareRects(player.rect, bullet.rect) == true {

					if bulletShooterP != player {

						//Remove bullet once it hits a player
						//removeBulletFromList(bullet)
						//bulletsToRemove = append(bulletsToRemove, bullet)
						addBulletToRemoveList(bullet)
						bulletRemoved = true

						//calculate damage dealing
						var damage = BASE_DAMAGE_VALUE + (bullet.damage * 5)

						//Player takes damage to shield until zero, then takes health damage
						var diff = player.shield - float64(damage)
						if diff >= 0 {
							player.shield -= float64(damage)
						} else {
							player.shield = 0
							player.health += diff
						}

						//player.Health = player.Health - 10

						if player.health <= 0 {
							player.rect.x = 0
							player.rect.y = 0
							player.health = 100

							//update shooter's scraps
							if bulletShooterP != nil {
								var shooter = *bulletShooterP
								shooter.scraps += 100
								shooter.xp += 100
								*bulletShooterP = shooter
							}
						}
					}
				}

			}

			// Check bullets for collision with npcs
			if bulletRemoved == false {
				for _, npc := range npcs {

					if compareRects(npc.rect, bullet.rect) == true && bulletShooterNPC == nil {

						//Remove bullet once it hits a player
						addBulletToRemoveList(bullet)

						//calculate damage dealing
						var damage = BASE_DAMAGE_VALUE + (bullet.damage * 5)

						npc.health -= float64(damage)

						//player.Health = player.Health - 10

						if npc.health <= 0 {
							println("remove npc")
							//remove npc
							removeNpcFromList(npc)

							//only update bullet shooters shit if it was shot by a player
							if bulletShooterP != nil {
								var shooter = *bulletShooterP
								//update shooter's scraps
								shooter.scraps += (int32(rand.Intn(51)) + 50)
								shooter.xp += 20
								*bulletShooterP = shooter

								//drop item randomly
								dropItemRandomly(bulletShooterP, 75)
							}
						}

					}

				}
			}
		}

		time.Sleep((time.Second / time.Duration(60)))
	}
}