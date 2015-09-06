//000bf? ►?Action?playerUpdate?BulletIDs??BulletRots??BulletXs??BulletYs??Gear??   ?Healthd?ID?1?IsNPCrOtherPlayerIDs??OtherPlayerRots??OtherPlayerXs??OtherPlayerYs??Rotation ?X?        ?Y?
//0007e??Action?playerUpdate?BulletIDs??BulletRots??BulletXs??BulletYs??Gear??   ?Healthd?ID?1?IsNPC"Rotation ?X?        ?Y?
package main

import (
	"io"
	"os"
	//"io/ioutil"
	"fmt"
	//"log"
	"math"
	"math/rand"
	"net"
	"strconv"
	"time"
	//"strings"
	//"net/http"
	"ugorji/go/codec"
)

type Player struct {
	rect              Rectangle
	addr              net.UDPAddr
	id                string
	level             int
	xp                int
	shooting          bool
	infamy            int
	health            float64
	healthCap         int
	healthRegen       int
	energy            float64
	energyCap         int
	energyRegen       int
	shield            float64
	shieldCap         int
	shieldRegen       int
	fireRate          int
	fireRateCooldown  int
	damage            int
	speed             int
	weaponCooldownCap float64
	weaponCooldown    float64
	weaponBulletCount int
	scraps            int32
	x                 float64
	y                 float64
	xMovement         float64
	yMovement         float64
	gear              []string
	inventory         []string
}

type Point struct {
	x float64
	y float64
}

type Npc struct {
	rect         Rectangle
	id           int
	npcType      int
	health       float64
	damage       int
	rotation     int
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
	shooter interface{}
	damage  int
	ID      int
	rect    Rectangle
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

const listenAddr = ":7777"

const baseAddr = "http://192.168.1.18:3000/api/v1/"

var (
	mh codec.MsgpackHandle
)

var partner = make(chan io.ReadWriteCloser)

var players []*Player
var bullets []*Bullet
var npcs []*Npc

var shouldQuit = false

//CONSTANTS
var PLAYER_LOAD_DIST float64 = 30
var ARENA_SIZE float64 = 100
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

//leveling
var BASE_XP = 100
var LEVEL_XP_FACTOR = 4

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
	matchmake()
}

func randomFloat(min, max float64) float64 {
	//return rand.Intn(max - min) + min
	return min + (rand.Float64() * ((max - min) + 1))
}

func randomInt(min, max int) int {
	return rand.Intn(max-min) + min
}

func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}

func spawnNPCs() {
	for i := 0; i < 110; i++ {
		newNPC := new(Npc)
		newNPC.id = rand.Intn(10000)
		newNPC.npcType = 1
		newNPC.shotTime = -1
		newNPC.shotCooldown = -1
		newNPC.health = 50
		newNPC.rect.rotation = 0
		var x = ((rand.Float64() * ARENA_SIZE) - (ARENA_SIZE / 2))
		var y = ((rand.Float64() * ARENA_SIZE) - (ARENA_SIZE / 2))
		newNPC.rect = createRect(x, y, 3, 3)

		npcs = append(npcs, newNPC)
	}

	//create test NPC type 2
	newNPC := new(Npc)
	newNPC.id = rand.Intn(10000)
	newNPC.npcType = 2
	newNPC.damage = 2
	newNPC.shotTime = 1000
	newNPC.shotCooldown = 1000
	newNPC.health = 50
	newNPC.rect.rotation = 0
	var x float64 = 0
	var y float64 = 0
	newNPC.rect = createRect(x, y, 3, 3)

	npcs = append(npcs, newNPC)
}

func updateNPCs() {
	for {
		for _, npc := range npcs {
			if npc.npcType == 2 {
				npc.rect.x = rand.Float64()*5 - 2.5
				npc.rect.y = rand.Float64()*5 - 2.5
			}

			if npc.shotTime != -1 {
				if npc.shotCooldown > 0 {
					npc.shotCooldown -= 1
				} else if npc.shotCooldown <= 0 {
					npc.shotCooldown = npc.shotTime
					handleShot("radialShotgunShot", npc, 0, npc.rect, &bullets)
				}
			}
		}
		time.Sleep((time.Second / time.Duration(1000)))
	}
}

func updatePlayerStats() {
	for {
		for _, player := range players {

			//update XP and Level
			var currentXPCap float64 = float64(BASE_XP) * (math.Pow(float64(player.level), float64(LEVEL_XP_FACTOR)))
			var currentXPCapRounded = int(currentXPCap) //round number
			if player.xp >= currentXPCapRounded {
				var diff = player.xp - currentXPCapRounded
				player.level += 1
				player.xp = diff
			}

			//update health stat
			hullHealthCapAttr := getItemAttribute(player.gear[0], "healthCap")
			var healthCap = float64(BASE_HEALTH_CAP_VALUE) + (10 * float64(player.healthCap)) + hullHealthCapAttr
			var healthRegen = BASE_HEALTH_REGEN_VALUE + (0.002 * float64(player.healthRegen))
			if player.health < healthCap {
				player.health += healthRegen
				if player.health > healthCap {
					player.health = healthCap
				}
			}

			//update shield stat
			var shieldCap = BASE_SHIELD_CAP_VALUE + (10 * float64(player.shieldCap))
			var shieldRegen = BASE_SHIELD_REGEN_VALUE + (0.01 * float64(player.shieldRegen))
			if player.shield < shieldCap {
				player.shield += shieldRegen
				if player.shield > shieldCap {
					player.shield = shieldCap
				}
			}

			//update energy stat
			var energyCap = BASE_ENERGY_CAP_VALUE + (10 * float64(player.energyCap))
			var energyRegen = BASE_ENERGY_REGEN_VALUE + (0.01 * float64(player.energyRegen))
			if player.energy < energyCap {
				player.energy += energyRegen
				if player.energy > energyCap {
					player.energy = energyCap
				}
			}

			//shoot
			var fireRate = BASE_FIRE_RATE_VALUE - (10 * player.fireRate)
			if player.shooting {
				player.fireRateCooldown -= 1

				if player.fireRateCooldown <= 0 {
					player.fireRateCooldown = fireRate

					// spawn new bullet
					handleShot("singleShot", player, player.damage, player.rect, &bullets)
				}
			} else {
				player.fireRateCooldown = 0
			}
		}
		time.Sleep((time.Second / time.Duration(1000)))
	}
}

func movePlayers() {
	for {
		for _, player := range players {

			wingSpeedAttr := getItemAttribute(player.gear[2], "speed")

			var speed = BASE_SPEED_VALUE + (0.5 * float64(player.speed+int(wingSpeedAttr)))
			player.rect.x = player.rect.x + (player.xMovement * (speed / 100))
			player.rect.y = player.rect.y + (player.yMovement * (speed / 100))

			if player.rect.x >= ARENA_SIZE {
				player.rect.x = ARENA_SIZE
			} else if player.rect.x <= -ARENA_SIZE {
				player.rect.x = -ARENA_SIZE
			}

			if player.rect.y >= ARENA_SIZE {
				player.rect.y = ARENA_SIZE
			} else if player.rect.y <= -ARENA_SIZE {
				player.rect.y = -ARENA_SIZE
			}
			//player.rect.rotation = player.Rotation
		}
		time.Sleep((time.Second / time.Duration(300)))
	}

}

func moveBullets() {
	for {

		for _, bullet := range bullets {
			var bulletRadians float64 = (float64(bullet.rect.rotation+90) / 180.0) * 3.14159
			bullet.rect.x = bullet.rect.x + (15 * 0.116 * math.Cos(bulletRadians))
			bullet.rect.y = bullet.rect.y + (15 * 0.116 * math.Sin(bulletRadians))
		}

		for _, bullet := range bullets {
			var bulletRemoved = false
			var bulletShooterP *Player
			var bulletShooterNPC *Npc

			//check if bullet.shooter is a player
			if p, ok := bullet.shooter.(*Player); ok {
				bulletShooterP = p
			} else {
				/* not player */
			}

			//check if bullet.shooter is an NPC
			if npc, ok := bullet.shooter.(*Npc); ok {
				bulletShooterNPC = npc
			} else {
				/* not player */
			}

			// Checkl bullets for collision with players
			for _, player := range players {
				if compareRects(player.rect, bullet.rect) == true {

					if bulletShooterP != player {

						//Remove bullet once it hits a player
						removeBulletFromList(bullet)
						bulletRemoved = true

						//calculate damage dealing
						var damage = BASE_DAMAGE_VALUE + (bullet.damage * 5)

						//Player takes damage to shield until zero, then takes health damage
						var diff = player.shield - float64(damage)
						if diff >= 0 {
							player.shield -= float64(damage)
						} else {
							player.shield = 0
							player.health += diff
						}

						//player.Health = player.Health - 10

						if player.health <= 0 {
							player.rect.x = 0
							player.rect.y = 0
							player.health = 100

							//update shooter's scraps
							if bulletShooterP != nil {
								var shooter = *bulletShooterP
								shooter.scraps += 100
								shooter.xp += 100
							}
						}
					}
				}

			}

			// Check bullets for collision with npcs
			if bulletRemoved == false {
				for _, npc := range npcs {

					if compareRects(npc.rect, bullet.rect) == true && bulletShooterNPC != npc {

						//Remove bullet once it hits a player
						removeBulletFromList(bullet)

						//calculate damage dealing
						var damage = BASE_DAMAGE_VALUE + (bullet.damage * 5)

						npc.health -= float64(damage)

						//player.Health = player.Health - 10

						if npc.health <= 0 {
							removeNpcFromList(npc)

							//only update bullet shooters shit if it was shot by a player
							if bulletShooterP != nil {
								var shooter = *bulletShooterP
								//update shooter's scraps
								shooter.scraps += (int32(rand.Intn(51)) + 50)
								shooter.xp += 20

								//drop item randomly
								dropItemRandomly(bulletShooterP, 75)
							}
						}

					}

				}
			}
		}

		time.Sleep((time.Second / time.Duration(60)))
	}
}

/* A Simple function to verify error */
func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}

func playerAlreadyExists(addr *net.UDPAddr) *Player {
	var foundPlayer *Player
	for _, player := range players {
		var playerAddress = &player.addr
		var inputAddress = addr
		if playerAddress.String() == inputAddress.String() {
			foundPlayer = player
		}
	}
	return foundPlayer
}

func matchmake() {
	fmt.Printf("Hosting match making server\n")

	// listener, err := net.Listen("tcp", listenAddr)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// for !shouldQuit {
	// 	c, err := listener.Accept()
	// 	// c.SetReadBuffer(1)
	// 	// c.SetWriteBuffer(1)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	go match(c)
	// }

	/* Lets prepare a address at any address at port 10001*/
	ServerAddr, err := net.ResolveUDPAddr("udp", listenAddr)
	CheckError(err)

	/* Now listen at selected port */
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	CheckError(err)
	defer ServerConn.Close()

	buf := make([]byte, 1024)

	for !shouldQuit {

		n, addr, err := ServerConn.ReadFromUDP(buf)

		var foundPlayer = playerAlreadyExists(addr)

		if foundPlayer == nil {
			fmt.Println("Must create player")
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

func instantiatePlayer(addr *net.UDPAddr) *Player {
	// setup new player and its stats
	newPlayer := new(Player)
	newPlayer.addr = *addr
	newPlayer.id = ""
	newPlayer.infamy = 0
	newPlayer.health = 100
	newPlayer.healthCap = 0
	newPlayer.healthRegen = 0
	newPlayer.energy = 50
	newPlayer.energyCap = 0
	newPlayer.energyRegen = 0
	newPlayer.shield = 10
	newPlayer.shieldCap = 0
	newPlayer.shieldRegen = 0 //per tenth of a second
	newPlayer.fireRate = 0
	newPlayer.fireRateCooldown = 0
	newPlayer.damage = 0
	newPlayer.speed = 0
	newPlayer.scraps = 0
	newPlayer.weaponCooldownCap = 0.5
	newPlayer.weaponCooldown = 0
	newPlayer.weaponBulletCount = 1

	newPlayer.rect = createRect(0, 0, 3, 3)

	newPlayer.gear = []string{
		"H1",
		"L1",
		"W1",
		"T1"}

	newPlayer.inventory = make([]string, 8)

	players = append(players, newPlayer)
	fmt.Printf("\nPlayer Joined!\n>")
	fmt.Println("%v", len(players))

	return newPlayer
	//go getDataFromPlayer(newPlayer)

}

func findPlayerIndex(p *Player) int {
	var i = 0
	var foundIndex = -1
	for _, player := range players {
		if p == player {
			foundIndex = i
		}
		i++
	}
	return foundIndex
}

func findPlayerIndexByID(pID string) int {
	var i = 0
	var foundIndex = -1
	for _, player := range players {
		if player.id == pID {
			foundIndex = i
		}
		i++
	}
	return foundIndex
}

func removePlayerFromList(p *Player) {
	var foundIndex = findPlayerIndex(p)

	if foundIndex != -1 {
		players = append(players[:foundIndex], players[foundIndex+1:]...)
	}
}

func removeBulletFromList(b *Bullet) {
	var i = 0
	var foundIndex = -1
	for _, bullet := range bullets {
		if b == bullet {
			foundIndex = i
		}
		i++
	}
	if foundIndex != -1 {
		bullets = append(bullets[:foundIndex], bullets[foundIndex+1:]...)
	}
}

func removeNpcFromList(n *Npc) {
	var i = 0
	var foundIndex = -1
	for _, npc := range npcs {
		if n == npc {
			foundIndex = i
		}
		i++
	}
	if foundIndex != -1 {
		npcs = append(npcs[:foundIndex], npcs[foundIndex+1:]...)
	}
}

func isItemInPlayerInventory(inv []string, itemID string) bool {
	var response = false
	for _, item := range inv {
		if item == itemID {
			response = true
		}
	}

	return response
}

func removeItemFromInventory(inv *[]string, itemID string) {
	var j = 0
	var foundIndex = -1
	var inventory = *inv
	for _, item := range inventory {
		if item == itemID {
			foundIndex = j
		}
		j++
	}
	if foundIndex != -1 {
		inventory = append(inventory[:foundIndex], inventory[foundIndex+1:]...)
	}
	*inv = inventory
}

func removeItemFromInventoryViaIndex(inv *[]string, index int) {
	var inventory = *inv
	inventory[index] = ""
	*inv = inventory
}

func addItemToInventory(inv *[]string, index int, itemID string) {
	if index != -1 {
		var inventory = *inv

		inventory[index] = itemID

		*inv = inventory
	}
}

func getNextOpenSlotInInventory(inv []string) int {
	var foundIndex = -1
	for i := 0; i < len(inv); i++ {
		if inv[i] == "" {
			foundIndex = i
			break
		}
	}
	return foundIndex
}

func dropItemRandomly(player *Player, chance int) {

	var randInt = rand.Intn(101)

	var openSlot = getNextOpenSlotInInventory(player.inventory)

	if randInt <= chance {
		var itemType = rand.Intn(3)
		if itemType == 0 {
			//hull
			var randItem = rand.Intn(NUMBER_OF_HULL_ITEMS) + 1
			var randItemID = "H" + strconv.Itoa(randItem)
			addItemToInventory(&player.inventory, openSlot, randItemID)
		} else if itemType == 1 {
			//wings
			var randItem = rand.Intn(NUMBER_OF_WING_ITEMS) + 1
			var randItemID = "W" + strconv.Itoa(randItem)
			addItemToInventory(&player.inventory, openSlot, randItemID)
		} else if itemType == 2 {
			//lasers
		}

	}
}

func getDataFromPlayer(player *Player, buf []byte, n int) {

	for {
		var shouldRemove = false

		// var stringData = string(buf[0:n])
		// fmt.Printf(stringData)
		// dec := json.NewDecoder(strings.NewReader(stringData))

		var res = &Update{}
		// dec.Decode(&res)

		// fmt.Printf("%v\n", buf[0:n])

		var r io.Reader

		dec := codec.NewDecoder(r, &mh)
		dec = codec.NewDecoderBytes(buf[0:n], &mh)
		err := dec.Decode(res)

		if err == nil {
			// fmt.Printf("%v", res.X)
			// res := &Update{}

			//decoder.Decode(n)
			// json.Unmarshal([]byte(buf[0:n]), &res)
			//fmt.Printf(res.ID )
			if res.Action == "update" {

				player.id = res.ID

				player.shooting = res.Shooting

				player.xMovement = res.X
				player.yMovement = res.Y

				player.rect.rotation = res.Rotation

				//print("%v", player.xMovement)

				// if (player.gear.cockpit == -1) {
				//     var link = baseAddr + "get_users_item_set?user_id=" + player.ID
				//     resp, err := http.Get(link)
				//     if err != nil {

				//     } else {
				//         // fmt.Printf("%v", resp)
				//         contents, err := ioutil.ReadAll(resp.Body)
				//         if err == nil {
				//             // fmt.Printf("%s\n", string(contents))
				//             dec2 := json.NewDecoder(strings.NewReader(string(contents)))
				//             var res = &Gear{}
				//             dec2.Decode(&res)

				//             player.gear.cockpit = res.Cockpit
				//             player.gear.lasers = res.Lasers
				//             player.gear.wings = res.Wings
				//             player.gear.jets = res.Jets
				//         } else {

				//         }
				//     }
				// }

				//test := string(buf[0:n])
			} else if res.Action == "equip" {

				//TO DO: ADD JSON SUPPORT

				//put item to equip data in variables
				var itemToEquip = player.inventory[res.Value] //H1
				//put item currently wearing into variabls
				var currentItem = ""

				if string(itemToEquip[0]) == "W" {
					currentItem = player.gear[2]
					//replace equiped item
					player.gear[2] = itemToEquip
				} else if string(itemToEquip[0]) == "H" {
					currentItem = player.gear[0] //H2
					//replace equiped item
					player.gear[0] = itemToEquip
				}

				//remove equipped item from inventory
				removeItemFromInventoryViaIndex(&player.inventory, res.Value)
				//add item that was replaced
				addItemToInventory(&player.inventory, res.Value, currentItem)

			} else if res.Action == "drop" {
				fmt.Printf("%v", res.Value)
				removeItemFromInventoryViaIndex(&player.inventory, res.Value)
			} else if res.Action == "shoot" {
				player.health = player.health
			} else if res.Action == "upgradeHealthCap" {
				if player.scraps >= 100 && player.healthCap < HEALTHCAP_CAP {
					player.healthCap += 1
					player.scraps -= 100
				}
			} else if res.Action == "upgradeHealthRegen" {
				if player.scraps >= 100 && player.healthRegen < HEALTH_REGEN_CAP {
					player.healthRegen += 1
					player.scraps -= 100
				}
			} else if res.Action == "upgradeSpeed" {
				if player.scraps >= 100 && player.speed < MOVE_SPEED_CAP {
					player.speed += 1
					player.scraps -= 100
				}
			} else if res.Action == "upgradeShieldCap" {
				if player.scraps >= 100 && player.shieldCap < SHIELD_CAP {
					player.shieldCap += 1
					player.scraps -= 100
				}
			} else if res.Action == "upgradeShieldRegen" {
				if player.scraps >= 100 && player.shieldRegen < SHIELD_REGEN_CAP {
					player.shieldRegen += 1
					player.scraps -= 100
				}
			} else if res.Action == "upgradeEnergyCap" {
				if player.scraps >= 100 && player.energyCap < ENERGY_CAP_CAP {
					player.energyCap += 1
					player.scraps -= 100
				}
			} else if res.Action == "upgradeEnergyRegen" {
				if player.scraps >= 100 && player.energyRegen < ENERGY_REGEN_CAP {
					player.energyRegen += 1
					player.scraps -= 100
				}
			} else if res.Action == "upgradeShieldRegen" {
				if player.scraps >= 100 && player.shieldRegen < SHIELD_REGEN_CAP {
					player.shieldRegen += 1
					player.scraps -= 100
				}
			} else if res.Action == "upgradeDamage" {
				if player.scraps >= 100 && player.damage < DAMAGE_CAP {
					player.damage += 1
					player.scraps -= 100
				}
			} else if res.Action == "upgradeFireRate" {
				if player.scraps >= 100 && player.fireRate < FIRE_RATE_CAP {
					player.fireRate += 1
					player.scraps -= 100
				}
			} else if res.Action == "jump" {
				if player.scraps >= 200 {
					player.scraps -= 200
					player.rect.x = res.X
					player.rect.y = res.Y
				}
			}
			// fmt.Printf( strconv.FormatFloat(res.X, 'f', 6, 64) )
		} else {
			shouldRemove = true
		}

		if shouldRemove == true {
			for i, otherPlayer := range players {
				if otherPlayer.addr.String() == player.addr.String() {
					players = append(players[:i], players[i+1:]...)
					return
				}
			}
		}
	}

}

func CToGoString(c []byte) string {
	n := -1
	for i, b := range c {
		if b == 0 {
			break
		}
		n = i
	}
	return string(c[:n+1])
}

func intToBinaryString(i int) string {
	//create header
	// var value = int64(i)
	// binary := strconv.FormatInt(value, 2)
	// var diff = 8 - len(binary)
	// for i := 0; i < diff; i++ {
	//     binary = "0" + binary;
	// }
	// return binary

	var header = strconv.FormatInt(int64(i), 16)

	var diff = 5 - len(header)
	for i := 0; i < diff; i++ {
		header = "0" + header
	}

	return header

}

// {"Action":"shoot", "ID":"87", "X":"0", "Y":"0", "Rotation":"22"}
func sendData() {

	var speedMod = 30
	var sentWelcomeMsg = false
	for {

		//var bulletPackets []*BulletUpdate

		for _, player := range players {

			bulletIDs := make([]int, 0)
			bulletXs := make([]float64, 0)
			bulletYs := make([]float64, 0)
			bulletRots := make([]int, 0)

			otherPlayerIDs := make([]string, 0)
			otherPlayerXs := make([]float64, 0)
			otherPlayerYs := make([]float64, 0)
			otherPlayerRots := make([]int, 0)
			otherPlayerHlths := make([]float64, 0)
			otherPlayerGearSets := make([][]string, 0)

			npcIDs := make([]int, 0)
			npcTypes := make([]int, 0)
			npcXs := make([]float64, 0)
			npcYs := make([]float64, 0)
			npcHlths := make([]float64, 0)

			if player.id != "" {

				//put all bullets into one array that are CLOSE TO THE PLAYER
				//WARNING: MAY CAUSE LAG
				for _, bullet := range bullets {

					var dist = math.Sqrt(math.Pow(bullet.rect.x-player.rect.x, 2) + math.Pow(bullet.rect.y-player.rect.y, 2))
					if dist <= PLAYER_LOAD_DIST {
						bulletIDs = append(bulletIDs, bullet.ID)
						bulletXs = append(bulletXs, bullet.rect.x)
						bulletYs = append(bulletYs, bullet.rect.y)
						bulletRots = append(bulletRots, bullet.rect.rotation)
					}

				}

				//get all data from other players
				//WARNING: MAY CAUSE LAG
				for _, otherPlayer := range players {
					if player.id == "1" && player != otherPlayer {
						//fmt.Printf("X: %v, Y: %v\n", otherPlayer.rect.x, otherPlayer.rect.y)
						//fmt.Printf("%v\n", player.Gear.wings)
					}

					var dist = math.Sqrt(math.Pow(otherPlayer.rect.x-player.rect.x, 2) + math.Pow(otherPlayer.rect.y-player.rect.y, 2))
					if dist <= PLAYER_LOAD_DIST && player != otherPlayer {
						otherPlayerIDs = append(otherPlayerIDs, otherPlayer.id)
						otherPlayerXs = append(otherPlayerXs, otherPlayer.rect.x)
						otherPlayerYs = append(otherPlayerYs, otherPlayer.rect.y)
						otherPlayerRots = append(otherPlayerRots, otherPlayer.rect.rotation)
						otherPlayerHlths = append(otherPlayerHlths, otherPlayer.health)

						gearSet := []string{otherPlayer.gear[0], otherPlayer.gear[1], otherPlayer.gear[2], otherPlayer.gear[3]}
						otherPlayerGearSets = append(otherPlayerGearSets, gearSet)
					}
				}

				//get all data from NPCs
				//WARNING: MAY CAUSE LAG
				for _, npc := range npcs {
					var dist = math.Sqrt(math.Pow(npc.rect.x-player.rect.x, 2) + math.Pow(npc.rect.y-player.rect.y, 2))
					if dist <= PLAYER_LOAD_DIST {
						npcIDs = append(npcIDs, npc.id)
						npcTypes = append(npcTypes, npc.npcType)
						npcXs = append(npcXs, npc.rect.x)
						npcYs = append(npcYs, npc.rect.y)
						npcHlths = append(npcHlths, npc.health)
					}
				}

				//create new gear obj for the other players current gear set
				// otherPlayersGear := gearSet{
				//     cockpit: otherPlayer.gear.cockpit,
				//     lasers: otherPlayer.gear.lasers,
				//     wings: otherPlayer.gear.wings,
				//     jets: otherPlayer.gear.jets,
				// }

				//new table that has multiple updates

				//create update packet
				res1D := &Update{
					Action:              "playerUpdate",
					ID:                  player.id,
					Level:               player.level,
					XP:                  player.xp,
					Infamy:              player.infamy,
					FireRate:            player.fireRate,
					Health:              player.health,
					HealthCap:           player.healthCap,
					HealthRegen:         player.healthRegen,
					Energy:              player.energy,
					EnergyCap:           player.energyCap,
					EnergyRegen:         player.energyRegen,
					Shield:              player.shield,
					ShieldCap:           player.shieldCap,
					ShieldRegen:         player.shieldRegen,
					Speed:               player.speed,
					Damage:              player.damage,
					Scraps:              player.scraps,
					X:                   player.rect.x,
					Y:                   player.rect.y,
					Rotation:            player.rect.rotation,
					Gear:                player.gear,
					Inventory:           player.inventory,
					IsNPC:               false,
					OtherPlayerIDs:      otherPlayerIDs,
					OtherPlayerXs:       otherPlayerXs,
					OtherPlayerYs:       otherPlayerYs,
					OtherPlayerRots:     otherPlayerRots,
					OtherPlayerHlths:    otherPlayerHlths,
					OtherPlayerGearSets: otherPlayerGearSets,
					BulletIDs:           bulletIDs,
					BulletXs:            bulletXs,
					BulletYs:            bulletYs,
					BulletRots:          bulletRots,
					NpcIDs:              npcIDs,
					NpcTypes:            npcTypes,
					NpcXs:               npcXs,
					NpcYs:               npcYs,
					NpcHlths:            npcHlths,
				}

				var w io.Writer
				var newByteArray []byte
				enc := codec.NewEncoder(w, &mh)
				enc = codec.NewEncoderBytes(&newByteArray, &mh)
				enc.Encode(res1D)

				var stringMessage = string(newByteArray)
				//create header
				var header = intToBinaryString(len(stringMessage))
				//format message message
				stringMessage = header + stringMessage
				//print message
				if player.id == "1" {
					//fmt.Printf(stringMessage + "\n")
				}
				// fmt.Printf("Header: %v\n", header);
				// fmt.Printf("Packet Length: %v\n", len(newByteArray));

				//send message
				//go fmt.Fprint(player.rwc, stringMessage)
			} else {
				if sentWelcomeMsg == false {
					sentWelcomeMsg = true
					res1F := &Message{
						Action: "message",
						Data:   "connected!",
					}

					var w io.Writer
					var newByteArray []byte
					enc := codec.NewEncoder(w, &mh)
					enc = codec.NewEncoderBytes(&newByteArray, &mh)
					enc.Encode(res1F)

					var stringMessage = string(newByteArray)
					//---create header---
					var header = intToBinaryString(len(stringMessage))
					//---add header---
					stringMessage = header + stringMessage

					stringMessage = "<?xml version=\"1.0\"?>\n<cross-domain-policy>\n<allow-access-from domain=\"*\" to-ports=\"7770-7780\"/>\n</cross-domain-policy>\n"
					//send message
					//go fmt.Fprint(player.rwc, stringMessage)
				}
			}
		}
		time.Sleep(1 * (time.Second / time.Duration(speedMod)))
	}

}

func cp(w io.Writer, r io.Reader, errc chan<- error) {
	_, err := io.Copy(w, r)
	errc <- err
}
