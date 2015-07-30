package main 

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
)

var data map[string]interface{}

func loadAllItemData() {
	//load file into byte array
	file, err := ioutil.ReadFile("items.txt")
	if (err != nil) {
		//if there is an error print it
		fmt.Println(err)
	} else {
		//else, do shit
		//unmarshal the json
		err := json.Unmarshal(file, &data);
		if (err != nil) {
			//if there is an error, print
			fmt.Println(err)
		} else {
			//if not, do shit
			fmt.Println(data)
		}
	}
}

func getItemAttribute(itemID string, attributeID string) float64 {
	var value float64 = 0
	//get data for this item
	itemData := data[itemID].(map[string]interface{})
	//look for specified attribute
	for k := range itemData {
        if (k == attributeID) {
        	value = itemData[attributeID].(float64)
        }
    }
	return value
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