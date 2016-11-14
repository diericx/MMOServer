package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var items = make(map[string]*Entity)
var itemData map[string]interface{}

type Item struct {
	StatsObj   Stats
	Rng        []int
	Name       string
	ResourceId string
	ItemType   string
}

func NewItem(name string, itemType string) Item {
	i := Item{}
	i.Name = name
	i.ResourceId = "default_shoulder"
	i.ItemType = itemType
	i.Rng = make([]int, 2)
	return i
}

func NewItemEntity(pos Vect2, size Vect2) *Entity {
	e := NewEntity(pos, size)
	e.entityType = "item"
	e.stats_calc1 = Stats{}
	//funcs

	items[e.id] = e

	return e
}

func NewStatAlterItemEntity(pos Vect2, value int) *Entity {
	e := NewItemEntity(pos, Vect2{x: 1, y: 1})
	e.entityType = "item-stat-alter"
	e.resourceId = "glowing_orb"
	e.stats_calc1.Energy = value
	e.onCollide = e.onItemStatAlterCollide
	return e
}

func NewItemPickupEntity(pos Vect2, name string, itemType string, resourceId string, s Stats) *Entity {
	e := NewItemEntity(pos, Vect2{x: 1, y: 1})
	e.entityType = "item-pickup"
	e.resourceId = resourceId
	e.onCollide = e.onItemPickupCollide
	e.inventory[0] = NewItem(name, itemType)
	e.inventory[0].ResourceId = resourceId
	e.inventory[0].StatsObj = s
	return e
}

func NewDefaultEquippedArray() map[string]Item {
	equ := make(map[string]Item)
	equ["weapon"] = NewItem("Marc Laser", "weapon")
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

func loadItemData() {
	file, err := ioutil.ReadFile("item-data.txt")
	if err != nil {
		//if there is an error print it
		fmt.Println(err)
	} else {
		//else, do shit
		//unmarshal the json
		err := json.Unmarshal(file, &itemData)
		if err != nil {
			//if there is an error, print
			fmt.Println(err)
		} else {
			//if not, do shit
			fmt.Println(itemData)
			//num := data["H1"].(map[string]interface{})
			//fmt.Println(num["healthCap"])
		}
	}
}
