package main

var items = make(map[string]*Entity)

func NewItem(pos Vect2, size Vect2) *Entity {
	item := NewEntity(pos, size)
	item.entityType = "player"
	//funcs

	items[item.id.String()] = item

	return item
}
