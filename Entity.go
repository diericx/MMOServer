package main

import (
	"github.com/zemirco/uid"
)

var entities = make(map[string]Entity)

//Entity The entity object
type Entity struct {
	Id string
	X  float32
	Y  float32
	Z  float32
}

//NewEntity create a new Entity object
func NewEntity() *Entity {

	e := Entity{
		Id: uid.New(5),
		X:  0,
		Y:  0,
		Z:  0,
	}

	entities[e.Id] = e

	return &e

}
