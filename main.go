package main

import "time"

var chanBufSize = 1024

//Create Channels
var inputChan = make(chan InputPacket, chanBufSize)

//FrameWaitTime time between frames
var FrameWaitTime float64 = 33

type ForLoopWaiter struct {
	start time.Time
}

func main() {
	InitConnection()

	go Listen()

	for {
		w := ForLoopWaiter{start: time.Now()}
		processInput()
		Send()
		w.waitForTime(FrameWaitTime)
	}

}

func processInput() {
	//Process inputChan
	for len(inputChan) > 0 {
		inputPacket := <-inputChan
		inputPacket.entity.X += float32(inputPacket.X)
		inputPacket.entity.Y += float32(inputPacket.Y)
	}
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
