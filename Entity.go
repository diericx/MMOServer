package main

import (
	"math"
	"math/rand"

	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
	"github.com/zemirco/uid"
)

var entities = make(map[string]Entity)

//Entity The entity object
type Entity struct {
	ID string
	//These positions are only used for exporting
	X             float32
	Y             float32
	CurrentPlayer bool
	mov           Vector2
	body          *chipmunk.Body
}

//NewEntity create a new Entity object
func NewEntity() *Entity {

	e := Entity{
		ID: uid.New(5),
		X:  0,
		Y:  0,
	}

	//create shape
	shape := chipmunk.NewBox(vect.Vector_Zero, vect.Float(1), vect.Float(1))
	shape.SetElasticity(0.95)

	//create a body for the shape
	e.body = chipmunk.NewBody(vect.Float(1), shape.Moment(float32(1)))
	e.body.SetPosition(vect.Vect{vect.Float(0), 0.0})
	e.body.SetAngle(vect.Float(rand.Float32() * 2 * math.Pi))
	e.body.AddShape(shape)
	space.AddBody(e.body)

	entities[e.ID] = e

	return &e

}
