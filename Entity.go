package main

import (
	"net"

	"github.com/satori/go.uuid"
)

type Vect2 struct {
	x float64
	y float64
}

type Body struct {
	pos        Vect2
	size       Vect2
	vel        Vect2
	continuous bool
	angle      float64
	points     []Vect2
}

type Entity struct {
	id         uuid.UUID
	entityType string
	body       Body
	addr       *net.UDPAddr
}

//holding array
var entities = make(map[string]*Entity)

func NewEntity(pos Vect2, size Vect2) *Entity {
	newEntity := Entity{}

	newEntity.id = uuid.NewV4()
	newEntity.body.pos = pos
	newEntity.body.size = size

	entities[newEntity.id.String()] = &newEntity

	return &newEntity
}
