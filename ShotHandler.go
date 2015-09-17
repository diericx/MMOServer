package main

import (
	//"fmt"
	// "encoding/json"
	// "io/ioutil"
	"math/rand"
)

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
		updateBulletAttributes(newBullet, npc.damage, npc.rect, npc.rect.rotation, npc.range)
	}

	return newBullet
}

func updateBulletAttributes(newBullet *Bullet, damage int, rect Rectangle, rotation int, range int) {
	newBullet.ID = rand.Intn(1000)
	newBullet.damage = damage
	newBullet.origin = Vector2{rect.x, rect.y}
	newBullet.rect = createRect(rect.x, rect.y, 0.17, 0.5)
	newBullet.rect.rotation = rotation
}