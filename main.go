package main

import (
	"hash/fnv"
	"math/rand"
	"strconv"
	"time"
)

var s1 = rand.NewSource(time.Now().UnixNano())
var r1 = rand.New(s1)

func main() {

	var FRAME_WAIT_TIME float64 = 33
	var spawnXaSD float64 = 500
	var spawnY float64 = 500

	for i := 0; i < 200; i++ {
		NewPlanet(Vect2{rand.Float64() * spawnX, rand.Float64() * spawnY}, Vect2{1, 1})
	}

	go listenForPackets()

	for {
		w := ForLoopWaiter{start: time.Now()}

		processServerInput()
		updateAttacks()
		updateEntities()
		processServerOutput()
		sendServerOutput()
		resetVariables()

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
