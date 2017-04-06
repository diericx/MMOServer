package main

import "time"

type ForLoopWaiter struct {
	start time.Time
}

//FrameWaitTime time between frames
var FrameWaitTime float64 = 33

func main() {
	InitConnection()
	go Listen()

	for {
		w := ForLoopWaiter{start: time.Now()}
		ProcessPlayerInput()
		Send()
		w.waitForTime(FrameWaitTime)
	}
}

//proccess all player input packets
func ProcessPlayerInput() {
	//process movement packets
	for len(MovePacketChan) > 0 {
		p := <-MovePacketChan
		//apply movement
		p.e.X += float32(p.X)
		p.e.Y += float32(p.Y)
		println(p.e.X)
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
