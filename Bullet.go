package main

//holding array
var bullets = make(map[string]*Entity)

func NewBullet(pos Vect2, size Vect2) *Entity {
	b := NewEntity(pos, size)
	b.body.pos = pos
	b.body.size = size
	b.entityType = "bullet"
	b.onUpdate = b.bulletUpdateFunc

	bullets[b.id.String()] = b

	return b
}

func (b *Entity) bulletUpdateFunc() {
	b.body.pos.x += b.body.vel.x
	b.body.pos.y += b.body.vel.y
}
