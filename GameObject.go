package main

//"github.com/vova616/chipmunk"

var gameObjects []*Entity

//Create DEFAULT gameObject
func NewGameObject(location Vect2, size Vect2) *Entity {
	newGameObject := NewEntity(nil, location, size)
	newGameObject.body.UserData = newGameObject
	newGameObject.entityType = "GameObject"
	gameObjects = append(gameObjects, newGameObject)
	return newGameObject
}

func NewBullet(loc Vect2, size Vect2, life int, damageVect Vect2) *Entity {
	newGameObject := NewGameObject(loc, size)
	newGameObject.value = life
	newGameObject.setDamage(damageVect)
	newGameObject.entityType = "Bullet"

	newGameObject.body.continuous = true

	newGameObject.onCollide = func(e *Entity) {
		if e != newGameObject.origin && ((e.entityType == "Npc" && newGameObject.origin.entityType != "Npc") || e.entityType == "Player") {
			e.takeDamage(newGameObject)
			newGameObject.value = 0
		}
	}

	newGameObject.onUpdate = func() {
		newGameObject.value--
	}

	return newGameObject
}
