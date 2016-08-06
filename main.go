package main

import "time"

//"github.com/vova616/chipmunk/vect"

// var alive = true
//
var WORLD_LIMIT = float64(1000)
var FRAME_WAIT_TIME float64 = 16

var serverInput = make(chan PlayerDataObject, 1000)
var serverOutput = make(chan PlayerDataObject, 1000)

func main() {

	//defer profile.Start(profile.CPUProfile).Stop()

	var server = NewServer("localhost:7777") //192.168.2.36

	//------Start server ops---------
	go server.listenForPlayers()
	go sendServerOutput()
	//
	initializeEntityManager()

	for {
		w := ForLoopWaiter{start: time.Now()}

		//process the input from server
		processServerInput()
		//add all current entities to the channel so
		//they can be processed by a gouroutine
		updateEntitiesChannels()
		//remove entites
		//!updates entities table
		removeEnities()
		// update entites cell Data and their positions
		// !requires updating the hash map
		updateEntitiesPositionAndCellData()
		//create packets
		processServerOutput()
		//send packets
		sendServerOutput()

		w.waitForTime(FRAME_WAIT_TIME)
	}

	// //Update Display (ALWAYS ON MAIN THREAD)
	// //InitializeDisplay()
}

//---ForLoopWaiter Definitoin
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
