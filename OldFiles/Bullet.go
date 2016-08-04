package main

import (
	"time"

	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
)

type Bullet struct {
	origin     *Player
	gameObject GameObject
	life       int
}

type BulletCollisions struct {
	e *Entity
	b *Bullet
}

func (bc BulletCollisions) CollisionPreSolve(arbiter *chipmunk.Arbiter) bool {
	//println("PreSolve")
	return true
}
func (bc BulletCollisions) CollisionPostSolve(arbiter *chipmunk.Arbiter) {
	//println("PostSolve")
}
func (bc BulletCollisions) CollisionExit(arbiter *chipmunk.Arbiter) {
	//println("Exit")
}

func (bc BulletCollisions) CollisionEnter(arbiter *chipmunk.Arbiter) bool {
	var bodyBEntity = arbiter.BodyB.UserData.(*Entity)
	if arbiter.BodyB != bc.e.origin.body {
		//bullet is not colliding with its owner
		bodyBEntity.takeDamage(bc.e.damage)
		RemoveBullet(bc.b)
	}
	return true
}

var bullets []*Bullet

func NewBullet(origin *Entity, location vect.Vect, size vect.Vect) *Bullet {
	var newBullet Bullet

	newBullet.gameObject = NewBulletGameObject(location, size, 60, 10 )

	newBullet.gameObject.entity.body.SetMass(0.01)

	bullets = append(bullets, &newBullet)

	return &newBullet
}

func (b *Bullet) update() {
	b.gameObject.value -= 1
	if b.gameObject.value <= 0 {
		RemoveBullet(b)
	}
}

func RemoveBullet(b *Bullet) {
	//TODO Remve entity object from list too
	for i, otherBullet := range bullets {
		if otherBullet.entity.id == b.entity.id {
			space.RemoveBody(otherBullet.gameObject.entity.body)
			bullets = append(bullets[:i], bullets[i+1:]...)
			return
		}
	}
}

func RemoveBulletViaBody(b *chipmunk.Body) {
	//TODO Remve entity object from list too
	//WARNING: Linear Time Function
	for i, otherPlayer := range players {
		if otherPlayer.gameObject.entity.body == b {
			space.RemoveBody(otherPlayer.gameObject.entity.body)
			players = append(players[:i], players[i+1:]...)
			return
		}
	}
}
