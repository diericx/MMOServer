package main

import (
	//"fmt"
	// "encoding/json"
	// "io/ioutil"
	"math/rand"
)

func handleShot(shotType string, shooter interface{}, damage int, rect rectangle, bullets *[]*bullet) {

	var bulletShooterP *Player
	var bulletShooterNPC *Npc

	//check if bullet.shooter is a player
	if p, ok := shooter.(*Player); ok {
		bulletShooterP = p
	} else {
		/* not player */
	}

	//check if bullet.shooter is an NPC
	if npc, ok := shooter.(*Npc); ok {
		bulletShooterNPC = npc
	} else {
		/* not player */
	}

	if shotType == "singleShot" {
		fireSingleShot(bulletShooterP, bullets)
	} else if shotType == "radialShotgunShot" {
		fireRadialShotgunShot(bulletShooterNPC, bullets)
	}
	// newBullet := new (bullet)
	// newBullet.ID = rand.Intn(1000)
	// newBullet.rect = createRect(0, 0, 0.17, 0.5)
	// newBullet.rect.rotation = 0
	// newBullet.shooter = nil
	// bullets = append(bullets, newBullet)
}

func fireSingleShot(player *Player, bullets *[]*bullet) {
	newBullet := new(bullet)
	newBullet.ID = rand.Intn(1000)
	newBullet.damage = player.Damage
	newBullet.rect = createRect(player.rect.x, player.rect.y, 0.17, 0.5)
	newBullet.rect.rotation = player.rect.rotation
	newBullet.shooter = player
	*bullets = append(*bullets, newBullet)
}

func fireRadialShotgunShot(shooter *Npc, bullets *[]*bullet) {
	var shooterObj = *shooter
	for i := 0; i < 8; i++ {
		newBullet := new(bullet)
		newBullet.ID = rand.Intn(1000)
		newBullet.damage = shooterObj.damage
		newBullet.rect = createRect(shooterObj.rect.x, shooterObj.rect.y, 0.17, 0.5)
		newBullet.rect.rotation = i * 45
		newBullet.shooter = shooter
		*bullets = append(*bullets, newBullet)
	}
}
