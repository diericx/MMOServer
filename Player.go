package main

import (
	"math"
	"net"

	"github.com/melvinmt/firebase"
	"github.com/vova616/chipmunk/vect"
)

type FirebaseUserData struct {
	Username string
}

//var players []*Entity
var players = make(map[string]*Entity)

var DEFAULT_PLAYER_SHOOT_TIME float64 = 100
var DEFAULT_PLAYER_SPEED float64 = 500
var DEFAULT_PLAYER_MAX_SPEED float64 = 500

func NewPlayer(addr *net.UDPAddr, location Vect2, size Vect2) *Entity {
	newPlayer := NewEntity(nil, location, size)
	newPlayer.entityType = "Player"
	newPlayer.body.UserData = newPlayer
	newPlayer.value = 1
	newPlayer.addr = addr
	newPlayer.targetX = 0
	newPlayer.targetY = 0
	newPlayer.speed = DEFAULT_PLAYER_SPEED
	newPlayer.maxSpeed = DEFAULT_PLAYER_MAX_SPEED
	newPlayer.speedDamp = vect.Float(3)
	newPlayer.shootTime = 100
	newPlayer.shootCooldown = 0
	newPlayer.active = true //TODO

	newPlayer.body.size = Vect2{x: 100, y: 100}
	//newPlayer.body.Shapes = make([]*chipmunk.Shape, 0)
	//newPlayer.body.AddShape(newPlayer.shape)

	newPlayer.shoot = newPlayer.playerShoot
	newPlayer.onUpdate = func() {
		//set player shit to default
		newPlayer.shootTime = DEFAULT_PLAYER_SHOOT_TIME
		newPlayer.speed = DEFAULT_PLAYER_SPEED
		newPlayer.maxSpeed = DEFAULT_PLAYER_MAX_SPEED
		//update score
		if newPlayer.score > 0 {
			newPlayer.score -= 0.05
		} else {
			newPlayer.score = 0
		}
		//Update leaderboard
		var foundSelf = false
		for i := 0; i < len(leaderboard); i++ {
			if newPlayer.username == "" {
				break
			}
			if leaderboard[i] == newPlayer {
				if foundSelf {
					leaderboard[i] = nil
				}
				foundSelf = true
				continue
			}
			if foundSelf {
				continue
			}
			if leaderboard[i] == nil {
				if !foundSelf {
					leaderboard[i] = newPlayer
				}
				break
			} else {
				if newPlayer.score > leaderboard[i].score {
					leaderboard[i] = newPlayer
					break
				}
			}
		}
		//check value
		if newPlayer.value == 0 {
			//DIE
			newPlayer.removeFromLeaderboard()
		}
		//check health
		if newPlayer.health <= 0 {
			//Drop shit

			//remove from chain
			newPlayer.removeFromChain()

			//deal out score
			for k, v := range newPlayer.damageTaken {
				//if the player hasn't left... add to its score
				if entities[k] != nil {
					entities[k].score += v
				}
			}

			//DIE
			newPlayer.health = 100
			newPlayer.score = 0
			newPlayer.damageTaken = make(map[string]float64)
			newPlayer.body.pos = Vect2{x: 0, y: 0}
		}
		//check power
		if newPlayer.power < newPlayer.powerMax {
			newPlayer.power++
		}
		if newPlayer.power >= newPlayer.powerMax {
			newPlayer.power = newPlayer.powerMax
		}
		//make sure player stays in bounds
		var playerPos = newPlayer.body.Position()
		if playerPos.x < -WORLD_LIMIT {
			playerPos.x = -WORLD_LIMIT
		}
		if playerPos.x > WORLD_LIMIT {
			playerPos.x = WORLD_LIMIT
		}
		if playerPos.y < -WORLD_LIMIT {
			playerPos.y = -WORLD_LIMIT
		}
		if playerPos.y > WORLD_LIMIT {
			playerPos.y = WORLD_LIMIT
		}

		//set player's velocity and Position
		newPlayer.body.pos = playerPos

		if newPlayer.shootCooldown > 0 {
			newPlayer.shootCooldown -= 10
		}

		if newPlayer.shootCooldown < 0 {
			newPlayer.shootCooldown = 0
		}
	}

	newPlayer.onCollide = func(other *Entity) {
		println("COLLIDED")
		if other.entityType == "bullet" {
			//take damage from Bullet
			newPlayer.takeDamage(other)
			//remove bullet
			other.value = 0
		}
	}

	players[newPlayer.id.String()] = newPlayer

	return newPlayer
}

func (p *Entity) playerShoot(e *Entity) {
	if p.shootCooldown != 0 {
		return
	}

	var damageVect = Vect2{x: 10, y: 20}
	var newBullet = NewBullet(Vect2{x: 0, y: 0}, Vect2{x: 32, y: 16}, 60, damageVect)
	newBullet.body.pos = p.body.pos
	newBullet.body.rotation = p.body.rotation + (math.Pi / 2)
	newBullet.body.vel = Vect2{x: math.Cos(newBullet.body.rotation) * 100, y: math.Sin(newBullet.body.rotation) * 100}
	newBullet.body.continuous = true
	newBullet.origin = p

	p.shootCooldown = p.shootTime

}

func (p *Entity) authorize(url string, token string) {
	//ref := firebase.NewReference(url).Auth(token)
	// f := firego.New(url, nil)
	// var v map[string]interface{}
	// if err := f.Value(&v); err != nil {
	// 	println(err)
	// }
	// p.active = true
	// fmt.Printf("%s\n", v["username"])
	// println(v["username"])
	ref := firebase.NewReference(url).Auth(token).Export(false)
	response := FirebaseUserData{}
	if err := ref.Value(&response); err != nil {
		println(err)
	}
	p.username = response.Username
	p.active = true
}
