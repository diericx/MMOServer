package main

var bulletsToRemove []*Bullet

func clearBulletRemoveList(bullets *[]Bullet) {
	for _, bullet := range bulletsToRemove {
		removeBulletFromList(bullets, bullet)
	}
	bulletsToRemove = make([]*Bullet, 0)
}

func removeBulletFromList(bullets *[]Bullet, b *Bullet) {
	var i = 0
	var foundIndex = -1
	var bulletsCopy = *bullets
	for _, bullet := range *bullets {
		if b.entity.id == bullet.entity.id {
			foundIndex = i
		}
		i++
	}

	if foundIndex != -1 {
		*bullets = append(bulletsCopy[:foundIndex], bulletsCopy[foundIndex+1:]...)
	}
}
