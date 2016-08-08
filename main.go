package main

import "time"

func main() {

	var FRAME_WAIT_TIME float64 = 33

	//create some entities
	// for i := 0; i < 10; i++ {
	// 	NewEntity(Vect2{x: float64(i * 500), y: 0}, Vect2{x: 10, y: 10})
	// }

	go listenForPackets()

	for {
		w := ForLoopWaiter{start: time.Now()}

		processServerInput()
		processServerOutput()
		updateEntities()
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
