package main

import (
	"math/rand"
	"time"
)

var s1 = rand.NewSource(time.Now().UnixNano())
var r1 = rand.New(s1)

func main() {

	var FRAME_WAIT_TIME float64 = 33

	fillEnergyCheckpointArray()

	//create some entities
	// for i := 0; i < 10; i++ {
	// 	NewStatAlterItemEntity(Vect2{x: float64(i * 2), y: 0}, 100)
	// }
	s := Stats{
		Damage: 100,
	}
	newWeapon := NewItemPickupEntity(Vect2{x: 3, y: 3}, "Energy Blaster", "weapon", s)
	newWeapon.inventory[0].Rng[0] = 100
	newWeapon.inventory[0].Rng[1] = 101

	s2 := Stats{
		Defense: 100,
	}
	NewItemPickupEntity(Vect2{x: -3, y: -3}, "Cold Shoulders", "shoulder", s2)
	//NewStatAlterItem(Vect2{x: 0, y: 0}, 100)

	go listenForPackets()

	for {
		w := ForLoopWaiter{start: time.Now()}

		processServerInput()
		updateEntities()
		processServerOutput()
		sendServerOutput()

		w.waitForTime(FRAME_WAIT_TIME)
	}
}

type ForLoopWaiter struct {
	start time.Time
}

func (flw ForLoopWaiter) waitForTime(maxMilliToWait float64) {
	var deltaTime = time.Since(flw.start)
	var delta = deltaTime.Seconds()
	var deltaMilli = delta * 1000
	deltaMilli = maxMilliToWait - deltaMilli

	if deltaMilli > 0 {
		time.Sleep(time.Duration(deltaMilli) * time.Millisecond)
	}
}
