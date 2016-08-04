package main

import (
	"math/rand"
	// "fmt"
)

var updateFunctions = make(map[int]func(e *Entity))

func update0(e *Entity) {

}

func update1(e *Entity) {
	// fmt.Println("UPDATEE!!  1")
}

func update2(e *Entity) {
	// fmt.Println("UPDATEE!!  2")
	moveEnemyRandomly(e)
}

func update3(e *Entity) {
	// fmt.Println("UPDATEE!!  3")
	moveEnemyRandomly(e)
}

func setupFunctions() {

	updateFunctions[0] = update0
	updateFunctions[1] = update1
	updateFunctions[2] = update2
	updateFunctions[3] = update3
}


////MOVE FUNCTIONS
func moveEnemyRandomly(e *Entity) {
	// if e.rect.pos.x == e.origin.x && e.rect.pos.y == e.origin.y {
		//move randomly if already on position
		var randMultiplierX = rand.Float64() * 8
		var randMultiplierY = rand.Float64() * 8
		// fmt.Println(randMultiplierX)
		e.rect.pos.x = e.origin.x + NPC_2_MAX_MOVE_DIST*randMultiplierX
		e.rect.pos.y = e.origin.y + NPC_2_MAX_MOVE_DIST*randMultiplierY
	// } else {
		// e.rect.pos.x =  NPC_2_MAX_MOVE_DIST*randMultiplierX
		// e.rect.pos.y =  NPC_2_MAX_MOVE_DIST*randMultiplierY

	// }
}