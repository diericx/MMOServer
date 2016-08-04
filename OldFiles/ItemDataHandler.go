package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strconv"
)

var data map[string]interface{}

func loadAllItemData() {
	//load file into byte array
	file, err := ioutil.ReadFile("items.txt")
	if err != nil {
		//if there is an error print it
		fmt.Println(err)
	} else {
		//else, do shit
		//unmarshal the json
		err := json.Unmarshal(file, &data)
		if err != nil {
			//if there is an error, print
			fmt.Println(err)
		} else {
			//if not, do shit
			fmt.Println(data)
			num := data["H1"].(map[string]interface{})
			fmt.Println(num["healthCap"])
		}
	}
}

//get attributes for an item
func getItemAttribute(itemID string, attributeID string) float64 {
	var value float64 = 0
	//get data for this item
	itemData := data[itemID].(map[string]interface{})
	//look for specified attribute
	for k := range itemData {
		if k == attributeID {
			value = itemData[attributeID].(float64)
		}
	}
	return value
}

//check if an item is int he player's inventory
func isItemInPlayerInventory(inv []string, itemID string) bool {
	var response = false
	for _, item := range inv {
		if item == itemID {
			response = true
		}
	}

	return response
}

//Remove an item from player's inventory via ITEM STRING
func removeItemFromInventory(inv *[]string, itemID string) {
	var j = 0
	var foundIndex = -1
	var inventory = *inv
	for _, item := range inventory {
		if item == itemID {
			foundIndex = j
		}
		j++
	}
	if foundIndex != -1 {
		inventory = append(inventory[:foundIndex], inventory[foundIndex+1:]...)
	}
	*inv = inventory
}

//Remove an item from player's inventory via ITEM INDEX
func removeItemFromInventoryViaIndex(inv *[]string, index int) {
	var inventory = *inv
	inventory[index] = ""
	*inv = inventory
}

//add an item to player's inventory
func addItemToInventory(inv *[]string, index int, itemID string) {
	if index != -1 {
		var inventory = *inv

		inventory[index] = itemID

		*inv = inventory
	}
}

//return next open slot in player's inventory
func getNextOpenSlotInInventory(inv []string) int {
	var foundIndex = -1
	for i := 0; i < len(inv); i++ {
		if inv[i] == "" {
			foundIndex = i
			break
		}
	}
	return foundIndex
}

//drop a random item in a player's inventory
func dropItemRandomly(player *Player, chance int) {

	var randInt = rand.Intn(101)

	var openSlot = getNextOpenSlotInInventory(player.inventory)

	if randInt <= chance {
		var itemType = rand.Intn(3)
		if itemType == 0 {
			//hull
			var randItem = rand.Intn(NUMBER_OF_HULL_ITEMS) + 1
			var randItemID = "H" + strconv.Itoa(randItem)
			addItemToInventory(&player.inventory, openSlot, randItemID)
		} else if itemType == 1 {
			//wings
			var randItem = rand.Intn(NUMBER_OF_WING_ITEMS) + 1
			var randItemID = "W" + strconv.Itoa(randItem)
			addItemToInventory(&player.inventory, openSlot, randItemID)
		} else if itemType == 2 {
			//lasers
		}

	}
}

// value := 0
// //get data for this item
// itemData := data[itemID].(map[string]interface{})
// //look for specified attribute
// for k := range itemData {
//        if (k == attributeID) {
//        	value =
//        }
//    }
// return itemData

// 	fmt.Println(data[itemID])
// itemData := data[itemID].(map[string]interface{})
// item1 := itemData["speed"].(float64)
// fmt.Println("%v", item1)
// for k := range itemData {
//        fmt.Println(k)
//    }
// return itemData
