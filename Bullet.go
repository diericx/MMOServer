package main

//holding array
var bullets = make(map[string]*Entity)

func NewBullet(pos Vect2, size Vect2, origin *Entity) *Entity {
	b := NewEntity(pos, size)
	b.origin = origin
	b.body.pos = pos
	b.body.size = size
	b.entityType = "bullet"
	b.resourceId = "default_bullet"
	b.onUpdate = b.bulletUpdateFunc
	b.onCollide = b.bulletCollideFunc
	b.value = origin.stats_calc.damage * float64(rnGenerator(origin.equipped["weapon"].rng))
	b.expireCounter = 30

	bullets[b.id.String()] = b

	return b
}

func rnGenerator(rng []int) int {
	var diff = rng[1] - rng[0]
	if diff <= 0 {
		return 0
	}

	return r1.Intn(rng[1]-rng[0]) + rng[0]
}

func (b *Entity) bulletUpdateFunc() {
	// b.body.pos.x += b.body.vel.x
	// b.body.pos.y += b.body.vel.y
}

func (b *Entity) bulletCollideFunc(other *Entity) {
	b.active = false
}
