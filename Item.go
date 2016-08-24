package main

var items = make(map[string]*Entity)

type Item struct {
	StatsObj Stats
	Rng      []int
	Name     string
}

func NewItem(name string) Item {
	i := Item{}
	i.Name = name
	i.Rng = make([]int, 2)
	return i
}

func NewItemEntity(pos Vect2, size Vect2) *Entity {
	e := NewEntity(pos, size)
	e.entityType = "item"
	e.stats = Stats{}
	//funcs

	items[e.id.String()] = e

	return e
}

func NewStatAlterItemEntity(pos Vect2, value int) *Entity {
	e := NewItemEntity(pos, Vect2{x: 1, y: 1})
	e.entityType = "item-stat-alter"
	e.resourceId = "glowing_orb"
	e.stats.Energy = value
	e.onCollide = e.onItemStatAlterCollide
	return e
}

func NewItemPickupEntity(pos Vect2, name string, s Stats) *Entity {
	e := NewItemEntity(pos, Vect2{x: 1, y: 1})
	e.inventory[0] = Item{
		StatsObj: s,
		Name:     name,
	}
	e.entityType = "item-pickup"
	e.resourceId = "default_item"
	e.onCollide = e.onItemPickupCollide
	return e
}

func NewDefaultEquippedArray() map[string]Item {
	equ := make(map[string]Item)
	equ["weapon"] = NewItem("Marc Laser")
	equ["weapon"].Rng[0] = 5
	equ["weapon"].Rng[1] = 10
	return equ
}

func (e *Entity) onItemStatAlterCollide(other *Entity) {
	e.active = false
	e.RemoveSelf()
}

func (e *Entity) onItemPickupCollide(other *Entity) {
	e.active = false
	e.RemoveSelf()
}
