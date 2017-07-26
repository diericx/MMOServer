package main

import (
	"time"

	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
)

var chanBufSize = 1024

var space *chipmunk.Space

//Create Channels
var inputChan = make(chan InputPacket, chanBufSize)

//FrameWaitTime time between frames
var updateFrameWaitTime float64 = 20
var updateDeltaTime = float32(updateFrameWaitTime) / 1000
var sendFrameWaitTime float64 = 100
var sendDeltaTime = float32(sendFrameWaitTime) / 1000

//ForLoopWaiter Holds start time for waiting in a loop
type ForLoopWaiter struct {
	start time.Time
}

//Vector2 A Vector2
type Vector2 struct {
	x float32
	y float32
}

func main() {
	InitConnection()
	initPhysics()

	go Listen()
	go Send()

	//set up ball
	var _ = NewBall(vect.Vect{0, 0})

	for {
		w := ForLoopWaiter{start: time.Now()}

		processInput()
		updateEntities()
		space.Step(vect.Float(updateDeltaTime))

		w.waitForTime(updateFrameWaitTime)
	}

}

func initPhysics() {

	space = chipmunk.NewSpace()
	space.Gravity = vect.Vect{0, 0}
	space.SetDamping(vect.Float(0.2))

}

func processInput() {
	//Process inputChan
	for len(inputChan) > 0 {
		inputPacket := <-inputChan
		inputPacket.entity.mov.x = float32(inputPacket.X)
		inputPacket.entity.mov.y = float32(inputPacket.Y)
	}
}

func updateEntities() {
	//Update player entities
	for _, p := range players {
		// update exported values
		p.e.X = float32(p.e.body.Position().X)
		p.e.Y = float32(p.e.body.Position().Y)
		// set velocity
		p.e.body.AddVelocity(p.e.mov.x/2, p.e.mov.y/2)
		//p.e.X += p.e.mov.x * updateDeltaTime * 5
		//p.e.Y += p.e.mov.y * updateDeltaTime * 5
	}

	//update balls
	for _, b := range balls {
		// update exported values
		b.e.X = float32(b.e.body.Position().X)
		b.e.Y = float32(b.e.body.Position().Y)
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
