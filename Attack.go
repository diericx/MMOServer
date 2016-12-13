package main

import (
	"math"
	"time"
)

type Attack struct {
	id         int
	start      time.Time
	travelTime time.Duration
	value      int
	origin     *Entity
	target     *Entity
}

var attacks = make(map[int]Attack)

var attackIdIncrement = 0

func NewAttack(t *Entity, o *Entity, d time.Duration, v int) {
	a := Attack{}
	a.id = getAttackId()
	a.start = time.Now()
	a.travelTime = d
	a.value = v
	a.origin = o
	a.target = t

	println(t, o)
	attacks[a.id] = a
}

func updateAttacks() {
	attacksToRemove := []int{}

	for _, a := range attacks {
		//if the value is zero, pass over and remove this attack obj
		if attacks[a.id].value <= 0 {
			attacksToRemove = append(attacksToRemove, a.id)
			continue
		}
		now := time.Now()
		var diff = now.Sub(a.start)
		if diff >= a.travelTime {
			if a.target.origin != a.origin.origin { //if its not your planet
				if a.target.stats.Count > 0 { //and the count is greater than 0
					//subtract from the count
					a.target.SetCount(a.target.stats.Count - 1)
					a.value -= 1
				} else { //if the count is 0 or less
					//add to the count and make it yours
					a.target.SetCount(a.target.stats.Count + 1)
					a.target.SetOrigin(a.origin.origin)
					a.value -= 1
				}
			} else { //if it is your planet
				a.target.SetCount(a.target.stats.Count + 1)
				a.value -= 1
			}
		}
		attacks[a.id] = a
	}

	//Remove attack objects
	for _, v := range attacksToRemove {
		delete(attacks, v)
	}
}

func getAttackId() int {
	if attackIdIncrement == math.MaxInt32 {
		attackIdIncrement = 0
	}
	attackId := attackIdIncrement
	attackIdIncrement += 1
	return attackId
}
