package main

import (
	"math"
	"math/rand"
	"net"
	"time"

	"github.com/satori/go.uuid"
	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
	// "fmt"
)

type Vect2 struct {
	x float64
	y float64
}

type Body struct {
	pos        Vect2
	size       Vect2
	vel        Vect2
	continuous bool
	rotation   float64
	points     []Vect2
	UserData   interface{}
}

func (b Body) Position() Vect2 {
	return b.pos
}

type Entity struct {
	//body
	body Body
	//Player Server Data
	addr    *net.UDPAddr
	targetX int
	targetY int
	//Player Data
	active        bool
	id            uuid.UUID
	username      string
	lastUpdate    time.Time
	score         float64
	shootCooldown float64
	shootTime     float64
	power         int
	powerMax      int
	speed         float64
	maxSpeed      float64
	exp           float64
	expToLevel    float64
	level         int
	shooting      bool
	health        float64
	healthCap     int
	damage        float64
	damageTaken   map[string]float64
	//chain data
	height int
	parent *Entity
	child  *Entity
	//Other Data
	tag        string
	entityType string
	value      int
	key        int
	originObj  *Entity
	origin     *Entity
	speedDamp  vect.Float
	//Functions
	onCollide func(e *Entity)
	onRemove  func()
	onUpdate  func()
	shoot     func(e *Entity)
}

type BodyData struct {
	tag string
}

type Collision struct {
	e *Entity
}

func (pc Collision) CollisionPreSolve(arbiter *chipmunk.Arbiter) bool {
	//println("PreSolve")
	return true
}

func (pc Collision) CollisionPostSolve(arbiter *chipmunk.Arbiter) {
	//println("PostSolve")
}
func (pc Collision) CollisionExit(arbiter *chipmunk.Arbiter) {
	//println("Exit")
}

func (pc Collision) CollisionEnter(arbiter *chipmunk.Arbiter) bool {
	var other *Entity
	var bodyAEntity = arbiter.BodyA.UserData.(*Entity)
	var bodyBEntity = arbiter.BodyB.UserData.(*Entity)
	if bodyAEntity.id == pc.e.id {
		other = bodyBEntity
	} else if bodyBEntity.id == pc.e.id {
		other = bodyAEntity
	}
	//Do shit
	if other != pc.e.origin && pc.e.onCollide != nil && other.value != 0 && other.health != 0 {
		pc.e.onCollide(other)
	}
	return true
}

//----VARIABLES----
//data channels
var entityIn = make(chan *Entity, 500)
var entityOut = make(chan *Entity, 500)

var entityToRemoveIn = make(chan *Entity, 500)

//holding array
var entities = make(map[string]*Entity)

//leaderboard shit
var leaderboard [5]*Entity

//-----FUNCS-----
func NewEntity(origin *Entity, location Vect2, size Vect2) *Entity {
	newEntity := Entity{}

	newEntity.id = uuid.NewV4()
	newEntity.active = true
	newEntity.origin = origin
	newEntity.key = 9999999999999

	newEntity.body.size = Vect2{x: size.x, y: size.y}
	newEntity.body.UserData = &newEntity
	newEntity.body.pos = location

	newEntity.updateEntityCellData()

	newEntity.damageTaken = make(map[string]float64)
	newEntity.health = 100
	newEntity.healthCap = 100
	newEntity.power = 100
	newEntity.powerMax = 200

	newEntity.body.points = make([]Vect2, 4)

	entities[newEntity.id.String()] = &newEntity

	return &newEntity
}

func (e *Entity) setDamage(v Vect2) {
	var diff = v.y - v.x
	if int(diff) < int(0) || int(v.x) < int(0) {
		println("ERROR: INVALID DAMAGE INPUT")
	}
	var damage = rand.Intn(int(diff)) + int(v.x)
	e.damage = float64(damage)
}

func (e *Entity) distanceTo(e2 *Entity) float64 {
	var loc1 = e.body.Position()
	var loc2 = e2.body.Position()
	var deltaX = float64(loc2.x - loc1.x)
	var deltaY = float64(loc2.y - loc1.y)
	return math.Sqrt((deltaX * deltaX) + (deltaY * deltaY))
}

func (e *Entity) lookAt(e2 *Entity) float64 {
	var loc1 = e.body.Position()
	var loc2 = e2.body.Position()
	var deltaX = float64(loc2.x - loc1.x)
	var deltaY = float64(loc2.y - loc1.y)
	var angle = math.Atan2(deltaY, deltaX)
	return angle
}

func (e *Entity) takeDamage(other *Entity) {
	e.health -= other.damage
	if other.origin != nil {
		e.damageTaken[other.origin.id.String()] += other.damage
	}
}

func (e *Entity) removeSelf() {
	//deal out score
	for k, v := range e.damageTaken {
		//if the player hasn't left... add to its score
		if entities[k] != nil {
			if e.entityType == "Player" {
				entities[k].score += v
			} else {
				entities[k].score += (v / 2)
			}
		}
	}

	e.removeFromChain()

	if e.onRemove != nil {
		e.onRemove()
	}

	//remove from map
	removeFromMap(e.key, e.id.String())
	delete(entities, e.id.String())

	//remove from players and entites list
	delete(players, e.id.String())
	delete(entities, e.id.String())
}

func removeEnities() {
L:
	for {
		select {
		case e, ok := <-entityToRemoveIn:
			if ok {
				e.removeSelf()
			}
		default:
			break L
		}
	}
}

//---

func processEntity(in, out chan *Entity) {
	//L:
	for {

		e := <-in

		// if it should be removed or health is 0, remove self
		if e.value <= 0 {
			entityToRemoveIn <- e
		}
		if e.entityType != "Player" && e.health <= 0 {
			entityToRemoveIn <- e
		}
		//Move entity according to velocity
		//e.moveEntity(Vect2{x: e.body.vel.x, y: e.body.vel.y})
		//call on update
		if e.onUpdate != nil {
			e.onUpdate()
		}

		out <- e

	}
}

func (e *Entity) moveEntity(v2 Vect2) {
	e.body.pos.x += v2.x
	e.body.pos.y += v2.y
	if e.detectCollisions() {
		e.body.pos.x -= v2.x
		e.body.pos.y -= v2.y
	}
}

func processAllEntities(in, out chan *Entity) {
	for _, e := range entities {
		in <- e
	}

	for i := len(entities); i > 0; i-- {
		<-out
	}
}

func updateEntitiesChannels() {
	processAllEntities(entityIn, entityOut)
}

func initializeEntityManager() {
	for i := 0; i < 4; i++ {
		go processEntity(entityIn, entityOut)
	}
}

//----

func updateEntitiesPositionAndCellData() {
	for _, e := range entities {
		e.moveEntity(e.body.vel)
		if !e.body.continuous {
			e.body.vel.x = 0
			e.body.vel.y = 0
		}

		e.updateEntityCellData()
	}
}

//LEADERBOARD DATA
func clearLeaderboard() {
	for i := 0; i < len(leaderboard); i++ {
		for j := 0; j < len(leaderboard); j++ {
			if i > j {
				if leaderboard[i] == leaderboard[j] {
					leaderboard[j] = nil
					break
				}
			}
		}
	}
}

func (e *Entity) removeFromLeaderboard() {
	for i := 0; i < len(leaderboard); i++ {
		if leaderboard[i] == e {
			leaderboard[i] = nil
		}
	}
}

func findEntityById(id string) *Entity {
	for _, e := range entities {
		if e.id.String() == id {
			return e
		}
	}
	return nil
}

func (e *Entity) setParent(e2 *Entity) {
	if e.parent != e2 && e2.parent != e && e2.child == nil && !e2.entityIsInChain(e) {
		if e.parent != nil {
			e.parent.child = nil
		}
		e.parent = e2
		e.parent.child = e
		e.parent.adjustHeight(e.height + 1)
	}
}

func (e *Entity) entityIsInChain(e2 *Entity) bool {
	if e == e2 {
		return true
	}
	if e.parent == nil {
		return false
	}
	return e.parent.entityIsInChain(e2)
}

func (e *Entity) removeChild() {
	e.child = nil
	e.adjustHeight(e.height - 1)
}

func (e *Entity) removeFromChain() {
	e.height = 0
	//remove self from connection
	if e.parent != nil {
		if e.child != nil {
			var parent = e.parent
			e.parent.removeChild()
			e.child.setParent(parent)
		} else {
			e.parent.removeChild()
		}
	} else {
		if e.child != nil {
			e.child.parent = nil
		}
	}
}

// func (e *Entity) removeParent(recursive bool) {
// 	e.parent = nil
//
// }

func (e *Entity) adjustHeight(value int) {
	e.height = value
	if e.parent != nil {
		e.parent.adjustHeight(value + 1)
	}
}
