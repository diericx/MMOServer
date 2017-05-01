package main

import (
	"github.com/zemirco/uid"
)

var entities = make(map[string]Entity)

//Entity The entity object
type Entity struct {
	ID            string
	X             float32
	Y             float32
	CurrentPlayer bool
	mov           Vector2
}

//NewEntity create a new Entity object
func NewEntity() *Entity {

	e := Entity{
		ID: uid.New(5),
		X:  0,
		Y:  0,
	}

	entities[e.ID] = e

	return &e

}
