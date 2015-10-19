//000bf? ►?Action?playerUpdate?BulletIDs??BulletRots??BulletXs??BulletYs??Gear??   ?Healthd?ID?1?IsNPCrOtherPlayerIDs??OtherPlayerRots??OtherPlayerXs??OtherPlayerYs??Rotation ?X?        ?Y?
//0007e??Action?playerUpdate?BulletIDs??BulletRots??BulletXs??BulletYs??Gear??   ?Healthd?ID?1?IsNPC"Rotation ?X?        ?Y?
package main

import (
	"io"
	//"io/ioutil"
	"fmt"
	//"log"
	"math/rand"
	"net"
	//"strconv"
	"time"
	//"strings"
	//"net/http"
)

type Player struct {
	rect               Rectangle
	addr               net.UDPAddr
	lastUpdate         time.Time
	id                 string
	level              int
	xp                 int
	skillPoints        int
	shooting           bool
	infamy             int
	health             float64
	healthCap          int
	healthRegen        int
	energy             float64
	energyCap          int
	energyRegen        int
	shield             float64
	shieldCap          int
	shieldRegen        int
	fireRate           int
	fireRateCooldown   int
	damage             int
	speed              int
	weaponCooldownCap  float64
	weaponCooldown     float64
	weaponBulletCount  int
	scraps             int32
	targetNPC          Npc
	targetNPC_rotation float64
	x                  float64
	y                  float64
	xMovement          float64
	yMovement          float64
	gear               []string
	inventory          []string
}

type Point struct {
	x float64
	y float64
}

type Npc struct {
	rect         Rectangle
	origin       Vector2
	id           int
	npcType      int
	health       float64
	alive        bool
	damage       int
	rotation     int
	bulletRange  int
	shotTime     int
	shotCooldown int
	shotType     string
}

type Rectangle struct {
	y        float64
	x        float64
	width    float64
	height   float64
	rotation int
	points   []Point
}

type Bullet struct {
	shooter     interface{}
	damage      int
	ID          int
	origin      Vector2
	rect        Rectangle
	bulletRange float64
}

type Vector2 struct {
	x float64
	y float64
}

type Message struct {
	Action string
	Data   string
}

type Update struct {
	//client player data
	Action      string
	Type        string
	Value       int
	ID          string
	Level       int
	XP          int
	SkillPoints int
	Infamy      int
	Shooting    bool
	FireRate    int
	Exp         int
	Health      float64
	HealthCap   int
	HealthRegen int
	Energy      float64
	EnergyCap   int
	EnergyRegen int
	Shield      float64
	ShieldCap   int
	ShieldRegen int
	Speed       int
	Damage      int
	Scraps      int32
	X           float64 //change to float64
	Y           float64
	Rotation    int
	Gear        []string
	IsNPC       bool
	//target marker data
	TargetNPC_rotation      float64
	TargetPlayers_rotations []float64
	//inventory data
	Inventory []string
	//other player data
	OtherPlayerIDs      []string
	OtherPlayerXs       []float64
	OtherPlayerYs       []float64
	OtherPlayerRots     []int
	OtherPlayerHlths    []float64
	OtherPlayerGearSets [][]string
	//bullet data
	BulletIDs  []int
	BulletXs   []float64
	BulletYs   []float64
	BulletRots []int
	//NPC data
	NpcIDs   []int
	NpcTypes []int
	NpcXs    []float64
	NpcYs    []float64
	NpcHlths []float64
}

type Gear struct {
	Success bool
	Cockpit int
	Lasers  int
	Wings   int
	Jets    int
}

type DamageTaken struct {
	ID       string
	Action   string
	BulletID int
}

type Shoot struct {
	Action   string
	ID       string
	X        float64
	Y        float64
	Rotation int
}

var partner = make(chan io.ReadWriteCloser)

var m = make(map[string][]interface{})

var shouldQuit = false

//CONSTANTS
var PLAYER_LOAD_DIST float64 = 30
var ARENA_SIZE float64 = 200
var HEALTHCAP_CAP int = 10
var HEALTH_REGEN_CAP int = 10
var MOVE_SPEED_CAP int = 10
var SHIELD_CAP int = 10
var SHIELD_REGEN_CAP int = 10
var ENERGY_CAP_CAP int = 10
var ENERGY_REGEN_CAP int = 10
var DAMAGE_CAP int = 10
var BULLET_SPEED_CAP int = 10
var FIRE_RATE_CAP int = 30

var BASE_HEALTH_CAP_VALUE int = 100
var BASE_HEALTH_REGEN_VALUE float64 = 0.001
var BASE_DAMAGE_VALUE int = 10
var BASE_SPEED_VALUE float64 = 3
var BASE_ENERGY_CAP_VALUE float64 = 50
var BASE_ENERGY_REGEN_VALUE float64 = 0.05
var BASE_SHIELD_CAP_VALUE float64 = 10
var BASE_SHIELD_REGEN_VALUE float64 = 0.01
var BASE_FIRE_RATE_VALUE int = 500

var HEALTH_MOD = 10
var ENERGY_MOD = 10

var NUMBER_OF_WING_ITEMS = 2
var NUMBER_OF_HULL_ITEMS = 2

var METEOR_MIN_DIST float64 = 75
var METEOR_MAX_DIST float64 = 210
var METEOR_MAX_AMMOUNT int = 300

var NPC_2_MIN_DIST float64 = 150
var NPC_2_MAX_DIST float64 = 100
var NPC_2_MAX_AMMOUNT int = 25
var NPC_2_MAX_MOVE_DIST float64 = 0.1
var NPC_2_SIGHT float64 = 25

//leveling
var BASE_XP = 50
var LEVEL_XP_FACTOR = 1

func main() {
	rand.Seed(time.Now().Unix())

	spawnNPCs()
	loadAllItemData()

	//go listenOn843()
	go getConsoleInput()
	go moveBullets()
	go movePlayers()
	go updatePlayerStats()
	go updateNPCs()
	go sendData()
	listenForPlayers()
}

func listenForPlayers() {
	/* Lets prepare an address at any address at port 10001*/
	ServerAddr, err := net.ResolveUDPAddr("udp", listenAddr)
	CheckError(err)

	/* Now listen at selected port */
	serverConn, err = net.ListenUDP("udp", ServerAddr)
	CheckError(err)
	defer serverConn.Close()

	buf := make([]byte, 1024)

	for !shouldQuit {

		n, addr, err := serverConn.ReadFromUDP(buf)

		var foundPlayer = playerAlreadyExists(addr)

		if foundPlayer == nil {
			println("Must create player")
			var player = instantiatePlayer(addr)
			getDataFromPlayer(player, buf, n)
		} else {
			getDataFromPlayer(foundPlayer, buf, n)
		}

		//fmt.Println("Received ", string(buf[0:n]), " from ", addr)

		if err != nil {
			fmt.Println("Error: ", err)
		}
	}
}

// {"Action":"shoot", "ID":"87", "X":"0", "Y":"0", "Rotation":"22"}
