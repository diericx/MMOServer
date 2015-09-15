package main

import (
	"fmt"
	// "encoding/json"
	// "io/ioutil"
	"math/rand"
)

func handleShot(shotType string, shooter interface{}, damage int, rect Rectangle, bullets *[]*Bullet) {

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

}

func fireSingleShot(shooter *Player, bullets *[]*Bullet) {
	fmt.Println("fire shot")
	newBullet := new(Bullet)
	newBullet.ID = rand.Intn(1000)
	newBullet.damage = shooter.damage
	newBullet.origin = Vector2{x: shooter.rect.x, y: shooter.rect.y}
	newBullet.rect = createRect(shooter.rect.x, shooter.rect.y, 0.17, 0.5)
	newBullet.rect.rotation = shooter.rect.rotation
	newBullet.shooter = shooter
	*bullets = append(*bullets, newBullet)
}

func fireRadialShotgunShot(shooter *Npc, bullets *[]*Bullet) {
	var shooterObj = *shooter
	for i := 0; i < 8; i++ {
		newBullet := new(Bullet)
		newBullet.ID = rand.Intn(1000)
		newBullet.damage = shooterObj.damage
		newBullet.origin = Vector2{x: shooter.rect.x, y: shooter.rect.y}
		newBullet.rect = createRect(shooterObj.rect.x, shooterObj.rect.y, 0.17, 0.5)
		newBullet.rect.rotation = i * 45
		newBullet.shooter = shooter
		*bullets = append(*bullets, newBullet)
	}
}
