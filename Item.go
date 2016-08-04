package main

import "github.com/vova616/chipmunk/vect"

type Item struct {
	name   string
	buff   Buff
	entity *Entity
}

//item target
var target *Entity

func NewItem(loc vect.Vect) *Entity {
	newGameObject := NewGameObject(loc, vect.Vect{X: 5, Y: 5})
	var buff = Buff{fireRate: -50, speed: -250}
	newGameObject.items[0] = Item{name: "Laser MKII", buff: buff}
	newGameObject.tag = newGameObject.items[0].name
	newGameObject.value = 1

	target = newGameObject

	newGameObject.onCollide = func(e *Entity) {
		if e.entityType == "Player" {
			for i, _ := range e.items {
				if e.items[i].name == "" {
					e.items[i] = newGameObject.items[0]
					target = e
					println("Collected 'Laser MKII' !")
					newGameObject.value = 0
					break
				}

			}

			newGameObject.removeSelf()
		}
	}

	newGameObject.onUpdate = func() {
		if newGameObject.value <= 0 {
			newGameObject.removeSelf()
		}
	}

	return newGameObject
}
