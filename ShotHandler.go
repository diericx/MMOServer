package main

import (
	"fmt"
	// "encoding/json"
	// "io/ioutil"
	"math/rand"
)

func handleShot(shotType string, player *player, damage int, rect rectangle, bullets *[]*bullet) {

	if shotType == "singleShot" {
		fireSingleShot(player, bullets)
	} else if shotType == "radialShotgunShot" {
		fireRadialShotgunShot(rect, damage, bullets)
	}
	// newBullet := new (bullet)
	// newBullet.ID = rand.Intn(1000)
	// newBullet.rect = createRect(0, 0, 0.17, 0.5)
	// newBullet.rect.rotation = 0
	// newBullet.shooter = nil
	// bullets = append(bullets, newBullet)
}

func fireSingleShot(player *player, bullets *[]*bullet) {
	newBullet := new(bullet)
	newBullet.ID = rand.Intn(1000)
	newBullet.damage = player.Damage
	newBullet.rect = createRect(player.rect.x, player.rect.y, 0.17, 0.5)
	newBullet.rect.rotation = player.rect.rotation
	newBullet.shooter = player
	*bullets = append(*bullets, newBullet)
}

func fireRadialShotgunShot(rect rectangle, damage int, bullets *[]*bullet) {
	newBullet := new(bullet)
	newBullet.ID = rand.Intn(1000)
	newBullet.damage = damage
	newBullet.rect = createRect(rect.x, rect.y, 0.17, 0.5)
	newBullet.rect.rotation = rect.rotation
	newBullet.shooter = nil
	*bullets = append(*bullets, newBullet)
}
