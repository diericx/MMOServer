package main

import "time"

type ForLoopWaiter struct {
	start time.Time
}

func main() {
	InitConnection()
	go Send()
	Listen()
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
