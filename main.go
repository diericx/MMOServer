package main

import (
	"hash/fnv"
	"math"
	"math/rand"
	"strconv"
	"time"
)

var s1 = rand.NewSource(time.Now().UnixNano())
var r1 = rand.New(s1)

func main() {

	var FRAME_WAIT_TIME float64 = 33

	fillEnergyCheckpointArray()
	loadItemData()

	//create some entities
	for i := 0; i < 20; i++ {
		NewStatAlterItemEntity(Vect2{x: math.Cos(float64(i)) * 5, y: math.Sin(float64(i)) * 5}, 100)
	}

	s := Stats{
		Damage: 100,
	}
	newWeapon := NewItemPickupEntity(Vect2{x: 3, y: 3}, "Energy Blaster", "weapon", "default_item", s)
	newWeapon.inventory[0].Rng[0] = 100
	newWeapon.inventory[0].Rng[1] = 101

	s2 := Stats{
		Defense: 100,
	}
	NewItemPickupEntity(Vect2{x: -3, y: -3}, "Cold Shoulders", "shoulder", "shoulder1", s2)
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

func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}

func Hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
