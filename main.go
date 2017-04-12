package main

import "time"

type ForLoopWaiter struct {
	start time.Time
}

//FrameWaitTime time between frames
var FrameWaitTime float64 = 50

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
	for len(InputPacketChan) > 0 {
		packet := <-InputPacketChan

		//TODO - Put player creation on the main thread
		p := GetPlayer(packet.addr)
		if p == nil {
			p = NewPlayer(packet.addr)
			println("New Player Connected!")
		}
		println(packet.Id)
		if packet.Id == 1 {
			// movePacket.e = p.e
			println(packet.X, packet.Y)
			p.e.X += float32(packet.X) * 1
			p.e.Y += float32(packet.Y) * 1
			// MovePacketChan <- movePacket
			// err = msgpack.Unmarshal(buf[:n], &input)
		}
		// println(p.X, p.Y)
		//apply movement
		// p.e.X += float32(p.X) * 10
		// p.e.Y += float32(p.Y) * 10
		//println(p.e.X)
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
