package main

var items = make(map[string]*Entity)

func NewItem(pos Vect2, size Vect2) *Entity {
	item := NewEntity(pos, size)
	item.entityType = "item"
	item.stats = Stats{}
	//funcs

	items[item.id.String()] = item

	return item
}

func NewStatAlterItem(pos Vect2, value int) *Entity {
	item := NewItem(pos, Vect2{x: 1, y: 1})
	item.resourceId = "glowing_orb"
	item.stats.energy = value
	item.onCollide = item.onStatAlterItemCollide
	return item
}

func (e *Entity) onStatAlterItemCollide(other *Entity) {
	e.active = false
	//actionObj := ServerActionObj{entity: e}
	e.RemoveSelf()
	//entitiesToRemove <- actionObj
}
