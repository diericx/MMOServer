//000bf? ►?Action?playerUpdate?BulletIDs??BulletRots??BulletXs??BulletYs??Gear??   ?Healthd?ID?1?IsNPCrOtherPlayerIDs??OtherPlayerRots??OtherPlayerXs??OtherPlayerYs??Rotation ?X?        ?Y?
//0007e??Action?playerUpdate?BulletIDs??BulletRots??BulletXs??BulletYs??Gear??   ?Healthd?ID?1?IsNPC"Rotation ?X?        ?Y?
package main 

import (
    "io"
    //"io/ioutil"
    "log"
    "net"
    "fmt"
    "time"
    "math/rand"
    "math"
    "strconv"
    "strings"
    "bufio"
    "os"
    //"strings"
    //"net/http"
    "ugorji/go/codec"
)

type player struct {
    rect rectangle
    RWC io.ReadWriteCloser
    ID string
    Level int 
    XP int
    Shooting bool
    Infamy int
    Health float64
    HealthCap int 
    HealthRegen int 
    Energy float64
    EnergyCap int 
    EnergyRegen int 
    Shield float64
    ShieldCap int 
    ShieldRegen int 
    FireRate int 
    FireRateCooldown int
    Damage int 
    Speed int 
    WeaponCooldownCap float64
    WeaponCooldown float64
    WeaponBulletCount int
    Scraps int32
    X float64
    Y float64
    xMovement float64
    yMovement float64
    Gear []string
    Inventory []string
}

type point struct {
    x float64
    y float64
}

type gearSet struct {
    hull int
    laser int
    wing int
    jet int
}

type Npc struct {
    rect rectangle
    ID int
    Type int
    Health float64
    rotation int
    shotTime int
    shotCooldown int
    shotType string
}

type rectangle struct {
    y float64
    x float64
    width float64
    height float64
    rotation int
    points []point
}

type bullet struct {
    shooter *player
    damage int
    ID int
    rect rectangle
}

type Message struct {
    Action string
    Data string
}

type Update struct {
    //client player data
    Action string
    Type string
    Value int
    ID string
    Level int 
    XP int
    Infamy int
    Shooting bool
    FireRate int
    Exp int
    Health float64
    HealthCap int
    HealthRegen int
    Energy float64 
    EnergyCap int
    EnergyRegen int
    Shield float64
    ShieldCap int
    ShieldRegen int
    Speed int
    Damage int
    Scraps int32
    X float64 //change to float64
    Y float64
    Rotation int
    Gear []string
    IsNPC bool
    //inventory data
    Inventory []string
    //other player data
    OtherPlayerIDs []string
    OtherPlayerXs []float64
    OtherPlayerYs []float64
    OtherPlayerRots []int
    OtherPlayerHlths []float64
    OtherPlayerGearSets [][]string
    //bullet data
    BulletIDs []int
    BulletXs []float64
    BulletYs []float64
    BulletRots []int
    //NPC data
    NpcIDs []int
    NpcTypes []int
    NpcXs []float64
    NpcYs []float64
    NpcHlths []float64
}

type BulletUpdate struct {
    ID string
    Damage string
    X float64 //change to float64
    Y float64
    Rotation int
}

type BulletPacket struct {
    Action string
    Bullets []*BulletUpdate
}

type Gear struct {
    Success bool
    Cockpit int
    Lasers int
    Wings int
    Jets int
}

type DamageTaken struct {
    ID string
    Action string
    BulletID int
}

type Shoot struct {
    Action string
    ID string
    X float64
    Y float64
    Rotation int
}

const listenAddr = ":7777"

const baseAddr = "http://192.168.1.18:3000/api/v1/"

var (
    mh codec.MsgpackHandle
)

var partner = make(chan io.ReadWriteCloser)

var players []*player
var bullets []*bullet
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
    // speed := data[0].(int)
    // fmt.Println(speed)

    // for i := 100; i < 110; i++ {
    //     newEnemy := new(Enemy)
    //     newEnemy.ID = strconv.Itoa(i)
    //     newEnemy.health = 100
    //     newEnemy.rect.rotation = 0
    //     var x = ((rand.Float64() * 10) - 5)
    //     var y = ((rand.Float64() * 10) - 5)
    //     newEnemy.rect = createRect(x, y, 1, 1)

    //     enemies = append(enemies, newEnemy)
    // }

    //go serverRuntime()
    go getConsoleInput()
    go moveBullets()
    go movePlayers()
    go updatePlayerStats()
    go updateNPCs()
    go chat()
    matchmake()
}

func serverRuntime() {
    //fmt.Printf(strconv.Itoa( len(players) ) )
    time.Sleep( 1  * (time.Second / time.Duration(1)) )
}

func randomFloat(min, max float64) float64 {
    //return rand.Intn(max - min) + min
    return min + (rand.Float64() * ((max - min) + 1))
}

func randomInt(min, max int) int {
    return rand.Intn(max - min) + min
}

func FloatToString(input_num float64) string {
    // to convert a float number to a string
    return strconv.FormatFloat(input_num, 'f', 6, 64)
}

func dotProduct(pointA point, pointB point) float64 {
    return math.Abs(pointA.x * pointB.x) + math.Abs(pointA.y * pointB.y)
}

func movePoints(rect rectangle ) {
    for _, p := range rect.points {
        p.x = p.x + rect.x
        p.y = p.y + rect.y
    }
}

func normalize(p point) point {
    var magnitude = math.Sqrt( p.x*p.x + p.y*p.y )
    if( magnitude > 0 ) {
        p.x = p.x / magnitude
        p.y = p.y / magnitude
    }
    return p
}

func spawnNPCs() {
    for i := 0; i < 110; i++ {
        newNPC := new(Npc)
        newNPC.ID = rand.Intn(10000)
        newNPC.Type = 1
        newNPC.shotTime = -1
        newNPC.shotCooldown = -1
        newNPC.Health = 50
        newNPC.rect.rotation = 0
        var x = ((rand.Float64() * ARENA_SIZE ) - (ARENA_SIZE/2))
        var y = ((rand.Float64() * ARENA_SIZE ) -  (ARENA_SIZE/2))
        newNPC.rect = createRect(x, y, 3, 3)

        npcs = append(npcs, newNPC)
    }

    //create test NPC type 2
    newNPC := new(Npc)
    newNPC.ID = rand.Intn(10000)
    newNPC.Type = 2
    newNPC.shotTime = 1000
    newNPC.shotCooldown = 1000
    newNPC.Health = 50
    newNPC.rect.rotation = 0
    var x float64 = 0
    var y float64= 0
    newNPC.rect = createRect(x, y, 3, 3)

    npcs = append(npcs, newNPC)
}

func getConsoleInput() {
    for {
        reader := bufio.NewReader(os.Stdin)
        fmt.Print(">")
        input, _ := reader.ReadString('\n')
        input = strings.TrimSpace(input)
        //do soemthing with it
        var keywords []string
        keywords = findAllKeywords(input)
        //if the user types give ("give player scraps 100")
        if (keywords[0] == "give") {
            //make sure they input the correct amount of parameters
            if ( len(keywords) ==  4) {
                playerID := findPlayerIndexByID(keywords[1])
                if (playerID >= 0) {
                    if (keywords[2] == "scraps") {
                        //try to convert value to string
                        i, err := strconv.ParseInt(keywords[3], 10, 32)
                        if err != nil {
                            fmt.Println("\"" + keywords[3] + "\" is not a valid integer!")
                        } else {
                            players[playerID].Scraps += int32(i)
                            fmt.Println("Succesfully gave player \"" + keywords[1] + "\" " + keywords[3] + " scraps!")
                        }
                    } else {
                        fmt.Println("\"" + keywords[2] + "\" is not a known command!")
                    }
                } else {
                    fmt.Println("Player with id \"" + keywords[1] + "\" not found!")
                }
            } else {
                fmt.Println("Not enough parameters supplied for \"give\" command!")
            }
        } else {
            fmt.Println("\"" + keywords[0] + "\" is not a known command!")
        }
    }
}

func findAllKeywords(val string) []string {
    input := val
    foundParams := make([]string, 0)

    for {
        i := strings.Index(input, " ")
        if (i == -1) {
            foundParams = append(foundParams, input)
            break
        } else {
            beforeSpace := input[:i]
            afterSpace := input[i+1:]
            foundParams = append(foundParams, beforeSpace)
            input = afterSpace
        }

    }

    return foundParams
}

func compareRects(objRect rectangle, bulletRect rectangle) bool {

    var pRectRot rectangle = objRect
    rotateRectsPoints(pRectRot, (float64(objRect.rotation) / 180.0) * 3.14159 )

    var bRectRot rectangle = bulletRect
    rotateRectsPoints(bRectRot, (float64(bulletRect.rotation) / 180.0) * 3.14159 )

    //CHECK X

    var v point
    v.x = pRectRot.width/2
    v.y = 0
    v = rotatePoint(v, (float64(objRect.rotation) / 180.0) * 3.14159 )

    v = normalize(v)

    var a point = pRectRot.points[2]

    var av float64 = dotProduct(a, v)

    var rv float64 = 0

    for i := 0; i < 4; i++ {
        var poop float64 = dotProduct(bRectRot.points[i], v)
        if (poop > rv) {
            rv = poop;
        }
    }

    movePoints(pRectRot)
    movePoints(bRectRot)

    var c point
    c.x = math.Abs(pRectRot.x - bRectRot.x)
    c.y = math.Abs(pRectRot.y - bRectRot.y)

    var cv float64 = dotProduct(c, v)
    var result bool = cv > av + rv

    //fmt.Printf("%t\n", result)

    //CHECK Y
    pRectRot = objRect
    rotateRectsPoints(pRectRot, (float64(objRect.rotation) / 180.0) * 3.14159 )

    bRectRot = bulletRect
    rotateRectsPoints(bRectRot, (float64(bulletRect.rotation) / 180.0) * 3.14159 )

    var w point
    w.y = pRectRot.height/2
    w.x = 0
    w = rotatePoint(w, (float64(objRect.rotation) / 180.0) * 3.14159 )

    w = normalize(w)

    var a2 = pRectRot.points[2]

    var aw float64 = dotProduct(a2, w)

    var rw float64 = 0

    for i := 0; i < 4; i++ {
        var poop float64 = dotProduct(bRectRot.points[i], w)
        if (poop > rw) {
            rw = poop;
        }
    }

    var cw float64 = dotProduct(c, w)
    var result2 bool = cw > aw + rw

    // fmt.Printf("%v %v \n", cw, aw + rw)

    // fmt.Printf("%t\n", result2)

    return !result && !result2 
}

func updateNPCs() {
    for {
        for _, npc := range npcs {
            if (npc.Type == 2) {
                npc.rect.x = rand.Float64()*5 - 2.5
                npc.rect.y = rand.Float64()*5 - 2.5
            }

            if (npc.shotTime != -1) {
                if (npc.shotCooldown > 0) {
                    npc.shotCooldown -= 1
                } else if (npc.shotCooldown <= 0) {
                    npc.shotCooldown = npc.shotTime
                    handleShot("radialShotgunShot", nil, 5, npc.rect, &bullets)
                }
            }
        }
        time.Sleep( (time.Second / time.Duration(1000)) )
    }
}

func updatePlayerStats() {
    for {
        for _, player := range players {

            //update XP and Level
            var currentXPCap float64 = float64(BASE_XP) * (math.Pow( float64(player.Level) , float64(LEVEL_XP_FACTOR) ) )
            var currentXPCapRounded = int(currentXPCap) //round number
            if (player.XP >= currentXPCapRounded) {
                var diff = player.XP - currentXPCapRounded
                player.Level += 1
                player.XP = diff
            }

            //update health stat
            hullHealthCapAttr := getItemAttribute(player.Gear[0], "healthCap")
            var healthCap = float64(BASE_HEALTH_CAP_VALUE) + (10 * float64(player.HealthCap) ) + hullHealthCapAttr
            var healthRegen = BASE_HEALTH_REGEN_VALUE + (0.002 * float64(player.HealthRegen) )
            if (player.Health < healthCap) {
                player.Health += healthRegen
                if (player.Health > healthCap) {
                    player.Health = healthCap
                }
            }

            //update shield stat
            var shieldCap = BASE_SHIELD_CAP_VALUE + (10 * float64(player.ShieldCap) )
            var shieldRegen = BASE_SHIELD_REGEN_VALUE + (0.01 * float64(player.ShieldRegen) )
            if (player.Shield < shieldCap) {
                player.Shield += shieldRegen
                if (player.Shield > shieldCap) {
                    player.Shield = shieldCap
                }
            }

            //update energy stat
            var energyCap = BASE_ENERGY_CAP_VALUE + (10 * float64(player.EnergyCap) )
            var energyRegen = BASE_ENERGY_REGEN_VALUE + (0.01 * float64(player.EnergyRegen) )
            if (player.Energy < energyCap) {
                player.Energy += energyRegen
                if (player.Energy > energyCap) {
                    player.Energy = energyCap
                }
            }

            //shoot
            var fireRate = BASE_FIRE_RATE_VALUE - (10 * player.FireRate )
            if (player.Shooting) {
                player.FireRateCooldown -= 1

                if (player.FireRateCooldown <= 0) {
                    player.FireRateCooldown = fireRate

                    // spawn new bullet
                    handleShot("singleShot", player, player.Damage, player.rect, &bullets)
                }
            } else {
                player.FireRateCooldown = 0
            }
        }
        time.Sleep( (time.Second / time.Duration(1000)) )
    }
}

func movePlayers() {
    for {
        for _, player := range players {

            wingSpeedAttr := getItemAttribute(player.Gear[2], "speed")

            var speed = BASE_SPEED_VALUE + (0.5 * float64(player.Speed + int(wingSpeedAttr) ) ) 
            player.rect.x = player.rect.x + (player.xMovement* (speed/100) )
            player.rect.y = player.rect.y + (player.yMovement* (speed/100) )


            if (player.rect.x >= ARENA_SIZE) {
                player.rect.x = ARENA_SIZE
            } else if (player.rect.x <= -ARENA_SIZE) {
                player.rect.x = -ARENA_SIZE
            }

            if (player.rect.y >= ARENA_SIZE) {
                player.rect.y = ARENA_SIZE
            } else if (player.rect.y <= -ARENA_SIZE) {
                player.rect.y = -ARENA_SIZE
            }
            //player.rect.rotation = player.Rotation
        }
        time.Sleep( (time.Second / time.Duration(300)) )
    }

}

func moveBullets() {
    for {

        for _, bullet := range bullets {
            var bulletRadians float64 = (float64(bullet.rect.rotation+90) / 180.0) * 3.14159
            bullet.rect.x = bullet.rect.x + (15 * 0.116 * math.Cos( bulletRadians ) )
            bullet.rect.y = bullet.rect.y + (15 * 0.116 * math.Sin( bulletRadians ) )
        }

        for _, bullet := range bullets {
            var bulletRemoved = false
        
            // Checkl bullets for collision with players
            for _, player := range players {
                if ( compareRects(player.rect, bullet.rect) == true ) {
                    var shouldContinue = false

                    if (bullet.shooter == nil) {
                        shouldContinue = true
                    } else {
                        if (bullet.shooter != player) {
                            shouldContinue = true
                        }
                    }

                    if (shouldContinue) {

                        //Remove bullet once it hits a player
                        removeBulletFromList(bullet)
                        bulletRemoved = true

                        //calculate damage dealing
                        var damage = BASE_DAMAGE_VALUE + (bullet.damage * 5)

                        //Player takes damage to shield until zero, then takes health damage
                        var diff = player.Shield - float64(damage)
                        if (diff >= 0) {
                            player.Shield -= float64(damage)
                        } else {
                            player.Shield = 0
                            player.Health += diff
                        }

                        //player.Health = player.Health - 10

                        if (player.Health <= 0) {
                            player.rect.x = 0
                            player.rect.y = 0
                            player.Health = 100

                            //update shooter's scraps
                            if (bullet.shooter != nil) {
                                bullet.shooter.Scraps += 100
                                bullet.shooter.XP += 100
                            }
                        }
                    }
                }
            
            }

            // Check bullets for collision with npcs
            if (bulletRemoved == false) {
                for _, npc := range npcs {
                    if ( compareRects(npc.rect, bullet.rect) == true ) {

                        //Remove bullet once it hits a player
                        removeBulletFromList(bullet)

                        //calculate damage dealing
                        var damage = BASE_DAMAGE_VALUE + (bullet.damage * 5)

                        npc.Health -= float64(damage)

                        //player.Health = player.Health - 10

                        if (npc.Health <= 0) {
                            removeNpcFromList(npc)

                            //only update bullet shooters shit if it was shot by a player
                            if (bullet.shooter != nil) {
                                //update shooter's scraps
                                bullet.shooter.Scraps += ( int32(rand.Intn(51)) + 50 )
                                bullet.shooter.XP += 20

                                //drop item randomly
                                dropItemRandomly(bullet.shooter, 75)
                            }
                        }

                    }
                
                }
            }
        }

        time.Sleep( (time.Second / time.Duration(60)) )
    }
}

func matchmake() {
    fmt.Printf( "Hosting match making server\n" )

    listener, err := net.Listen("tcp", listenAddr)
    if err != nil {
        log.Fatal(err)
    }
    for !shouldQuit {
        c, err := listener.Accept()
        // c.SetReadBuffer(1)
        // c.SetWriteBuffer(1)
        if err != nil {
            log.Fatal(err)
        }
        go match(c)
    }
}

func createRect(x float64, y float64, width float64, height float64) rectangle {

    var newRect rectangle 

    newRect.x = x
    newRect.y = y
    newRect.width = width
    newRect.height = height

    var w2 = width/2
    var h2 = height/2

    var point1 point
    point1.x = -w2
    point1.y = -h2
    newRect.points = append(newRect.points, point1)

    var point2 point
    point2.x = w2
    point2.y = -h2
    newRect.points = append(newRect.points, point2)

    var point3 point
    point3.x = w2
    point3.y = h2
    newRect.points = append(newRect.points, point3)

    var point4 point
    point4.x = -w2
    point4.y = h2
    newRect.points = append(newRect.points, point4)

    return newRect

}

func match(c io.ReadWriteCloser) {
    // setup new player and its stats
    newPlayer := new(player)
    newPlayer.RWC = c
    newPlayer.ID = ""
    newPlayer.Infamy = 0
    newPlayer.Health = 100
    newPlayer.HealthCap = 0
    newPlayer.HealthRegen = 0
    newPlayer.Energy = 50
    newPlayer.EnergyCap = 0
    newPlayer.EnergyRegen = 0
    newPlayer.Shield = 10
    newPlayer.ShieldCap = 0
    newPlayer.ShieldRegen = 0 //per tenth of a second
    newPlayer.FireRate = 0
    newPlayer.FireRateCooldown = 0
    newPlayer.Damage = 0
    newPlayer.Speed = 0
    newPlayer.Scraps = 0
    newPlayer.WeaponCooldownCap = 0.5
    newPlayer.WeaponCooldown = 0
    newPlayer.WeaponBulletCount = 1

    newPlayer.rect = createRect(0, 0, 3, 3)

    newPlayer.Gear = []string{
        "H1", 
        "L1", 
        "W1", 
        "T1"}

    newPlayer.Inventory = make([]string, 8)

    players = append(players, newPlayer)
    fmt.Printf("\nPlayer Joined!\n>")
    go getDataFromPlayer(newPlayer)
    
}

func rotatePoint(p point, angle float64) point {
    var newP point
    newP.x = ( p.x * math.Cos(angle) ) - ( p.y * math.Sin(angle) )
    newP.y = ( p.x * math.Sin(angle) ) - ( p.y * math.Cos(angle) )

    return newP
}

func rotateRectsPoints(r rectangle, angle float64) {
    for _, p := range r.points {
        p = rotatePoint(p, angle)
    }
}

func findPlayerIndex(p *player) int {
    var i = 0;
    var foundIndex = -1;
    for _, player := range players {
        if (p == player) {
            foundIndex = i;
        }
        i++
    }
    return foundIndex
}

func findPlayerIndexByID(pID string) int {
    var i = 0;
    var foundIndex = -1;
    for _, player := range players {
        if (player.ID == pID) {
            foundIndex = i;
        }
        i++
    }
    return foundIndex
}

func removePlayerFromList(p *player) {
    var foundIndex = findPlayerIndex(p);

    if (foundIndex != -1) {
        players = append(players[:foundIndex], players[foundIndex+1:]...)
    }
}

func removeBulletFromList(b *bullet) {
    var i = 0;
    var foundIndex = -1;
    for _, bullet := range bullets {
        if (b == bullet) {
            foundIndex = i;
        }
        i++
    }
    if (foundIndex != -1) {
        bullets = append(bullets[:foundIndex], bullets[foundIndex+1:]...)
    }
}

func removeNpcFromList(n *Npc) {
    var i = 0;
    var foundIndex = -1;
    for _, npc := range npcs {
        if (n == npc) {
            foundIndex = i;
        }
        i++
    }
    if (foundIndex != -1) {
        npcs = append(npcs[:foundIndex], npcs[foundIndex+1:]...)
    }
}

func isItemInPlayerInventory(inv []string, itemID string) bool {
    var response = false
    for _, item := range inv {
        if (item == itemID) {
            response = true
        }
    }

    return response
}

func removeItemFromInventory(inv *[]string, itemID string) {
    var j = 0;
    var foundIndex = -1;
    var inventory = *inv
    for _, item := range inventory {
        if (item == itemID) {
            foundIndex = j;
        }
        j++
    }
    if (foundIndex != -1) {
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
    if (index != -1) {
        var inventory =  *inv

        inventory[index] = itemID

        *inv = inventory
    }
}

func getNextOpenSlotInInventory(inv []string) int {
    var foundIndex = -1
    for i := 0; i < len(inv); i++ {
        if (inv[i] == "") {
           foundIndex = i
           break
        }
    }
    return foundIndex
}

func dropItemRandomly(player *player, chance int) {
    
    var randInt = rand.Intn(101)

    var openSlot = getNextOpenSlotInInventory(player.Inventory)

    if (randInt <= chance) {
        var itemType = rand.Intn(3)
        if (itemType == 0) {
            //hull
            var randItem = rand.Intn(NUMBER_OF_HULL_ITEMS)+1
            var randItemID = "H"+strconv.Itoa(randItem)
            addItemToInventory(&player.Inventory, openSlot, randItemID)
        } else if (itemType == 1) {
            //wings
            var randItem = rand.Intn(NUMBER_OF_WING_ITEMS)+1
            var randItemID = "W"+strconv.Itoa(randItem)
            addItemToInventory(&player.Inventory, openSlot, randItemID)
        } else if (itemType == 2) {
            //lasers
        }

        
        
    }
}

func getDataFromPlayer(player *player) {

    for {
        var shouldRemove = false

        buf := make([]byte, 1024)
        n, err := player.RWC.Read(buf)

        if (err != nil) {
            removePlayerFromList(player)
            break
        }

        // var stringData = string(buf[0:n])
        // fmt.Printf(stringData)
        // dec := json.NewDecoder(strings.NewReader(stringData))

        var res = &Update{}
        // dec.Decode(&res)

        // fmt.Printf("%v\n", buf[0:n])

        dec := codec.NewDecoder(player.RWC, &mh)
        dec = codec.NewDecoderBytes(buf[0:n], &mh)
        err = dec.Decode(res)

        if err == nil {
            // fmt.Printf("%v", res.X)
            // res := &Update{}

            //decoder.Decode(n)
            // json.Unmarshal([]byte(buf[0:n]), &res)
            //fmt.Printf(res.ID )
            if res.Action == "update" {

                player.ID = res.ID

                player.Shooting = res.Shooting

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
                var itemToEquip = player.Inventory[res.Value] //H1
                //put item currently wearing into variabls
                var currentItem = ""

                if (string(itemToEquip[0]) == "W"){
                    currentItem = player.Gear[2]
                    //replace equiped item
                    player.Gear[2] = itemToEquip
                } else if ( string(itemToEquip[0]) == "H" ) {
                    currentItem = player.Gear[0] //H2
                    //replace equiped item
                    player.Gear[0] = itemToEquip
                }

                //remove equipped item from inventory
                removeItemFromInventoryViaIndex(&player.Inventory, res.Value)
                //add item that was replaced
                addItemToInventory(&player.Inventory, res.Value, currentItem )

            }else if res.Action == "drop" {
                fmt.Printf("%v", res.Value)
                removeItemFromInventoryViaIndex(&player.Inventory, res.Value)
            } else if res.Action == "shoot" {
                player.Health = player.Health
            } else if (res.Action == "upgradeHealthCap") {
                if (player.Scraps >= 100 && player.HealthCap < HEALTHCAP_CAP) {
                    player.HealthCap += 1
                    player.Scraps -= 100;
                }
            } else if (res.Action == "upgradeHealthRegen") {
                if (player.Scraps >= 100 && player.HealthRegen < HEALTH_REGEN_CAP) {
                    player.HealthRegen += 1
                    player.Scraps -= 100;
                }
            } else if (res.Action == "upgradeSpeed") {
                if (player.Scraps >= 100 && player.Speed < MOVE_SPEED_CAP) {
                    player.Speed += 1
                    player.Scraps -= 100;
                }
            } else if (res.Action == "upgradeShieldCap") {
                if (player.Scraps >= 100 && player.ShieldCap < SHIELD_CAP) {
                    player.ShieldCap += 1
                    player.Scraps -= 100;
                }
            } else if (res.Action == "upgradeShieldRegen") {
                if (player.Scraps >= 100 && player.ShieldRegen < SHIELD_REGEN_CAP) {
                    player.ShieldRegen += 1
                    player.Scraps -= 100;
                }
            } else if (res.Action == "upgradeEnergyCap") {
                if (player.Scraps >= 100 && player.EnergyCap < ENERGY_CAP_CAP) {
                    player.EnergyCap += 1
                    player.Scraps -= 100;
                }
            } else if (res.Action == "upgradeEnergyRegen") {
                if (player.Scraps >= 100 && player.EnergyRegen < ENERGY_REGEN_CAP) {
                    player.EnergyRegen += 1
                    player.Scraps -= 100;
                }
            } else if (res.Action == "upgradeShieldRegen") {
                if (player.Scraps >= 100 && player.ShieldRegen < SHIELD_REGEN_CAP) {
                    player.ShieldRegen += 1
                    player.Scraps -= 100;
                }
            } else if (res.Action == "upgradeDamage") {
                if (player.Scraps >= 100 && player.Damage < DAMAGE_CAP) {
                    player.Damage += 1
                    player.Scraps -= 100;
                }
            } else if (res.Action == "upgradeFireRate") {
                if (player.Scraps >= 100 && player.FireRate < FIRE_RATE_CAP) {
                    player.FireRate += 1
                    player.Scraps -= 100;
                }
            } else if (res.Action == "jump") {
                if (player.Scraps >= 200) {
                    player.Scraps -= 200
                    player.rect.x = res.X;
                    player.rect.y = res.Y;
                }
            }
            // fmt.Printf( strconv.FormatFloat(res.X, 'f', 6, 64) )
        } else {
            shouldRemove = true
        }

        if shouldRemove == true {
            for i,otherPlayer := range players {
                if otherPlayer.RWC == player.RWC {
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
        header = "0" + header;
    }

    return header

}

// {"Action":"shoot", "ID":"87", "X":"0", "Y":"0", "Rotation":"22"}
func chat () {

    var speedMod = 30
    var sentWelcomeMsg = false
    for { 

        //var bulletPackets []*BulletUpdate

        for _, player := range players {

            bulletIDs := make([]int, 0);
            bulletXs := make([]float64, 0);
            bulletYs := make([]float64, 0);
            bulletRots := make([]int, 0);

            otherPlayerIDs := make([]string, 0);
            otherPlayerXs := make([]float64, 0);
            otherPlayerYs := make([]float64, 0);
            otherPlayerRots := make([]int, 0);
            otherPlayerHlths := make([]float64, 0);
            otherPlayerGearSets := make([][]string, 0);

            npcIDs := make([]int, 0);
            npcTypes := make([]int, 0);
            npcXs := make([]float64, 0);
            npcYs := make([]float64, 0);
            npcHlths := make([]float64, 0);

            if (player.ID != "") {

                //put all bullets into one array that are CLOSE TO THE PLAYER
                //WARNING: MAY CAUSE LAG
                for _, bullet := range bullets {

                    var dist = math.Sqrt( math.Pow(bullet.rect.x - player.rect.x, 2) + math.Pow(bullet.rect.y - player.rect.y, 2) )
                    if (dist <= PLAYER_LOAD_DIST) {
                        bulletIDs = append(bulletIDs, bullet.ID);
                        bulletXs = append(bulletXs, bullet.rect.x);
                        bulletYs = append(bulletYs, bullet.rect.y);
                        bulletRots = append(bulletRots, bullet.rect.rotation);
                    }

                }

                //get all data from other players
                //WARNING: MAY CAUSE LAG
                for _, otherPlayer := range players {
                    if (player.ID == "1" && player != otherPlayer) {
                        //fmt.Printf("X: %v, Y: %v\n", otherPlayer.rect.x, otherPlayer.rect.y)
                        //fmt.Printf("%v\n", player.Gear.wings)
                    }

                    var dist = math.Sqrt( math.Pow(otherPlayer.rect.x - player.rect.x, 2) + math.Pow(otherPlayer.rect.y - player.rect.y, 2) )
                    if (dist <= PLAYER_LOAD_DIST && player != otherPlayer) {
                        otherPlayerIDs = append(otherPlayerIDs, otherPlayer.ID);
                        otherPlayerXs = append(otherPlayerXs, otherPlayer.rect.x);
                        otherPlayerYs = append(otherPlayerYs, otherPlayer.rect.y);
                        otherPlayerRots = append(otherPlayerRots, otherPlayer.rect.rotation);
                        otherPlayerHlths = append(otherPlayerHlths, otherPlayer.Health);

                        gearSet := []string{otherPlayer.Gear[0], otherPlayer.Gear[1], otherPlayer.Gear[2], otherPlayer.Gear[3]}
                        otherPlayerGearSets = append(otherPlayerGearSets, gearSet);
                    }
                }

                //get all data from NPCs
                //WARNING: MAY CAUSE LAG
                for _, npc := range npcs {
                    var dist = math.Sqrt( math.Pow(npc.rect.x - player.rect.x, 2) + math.Pow(npc.rect.y - player.rect.y, 2) )
                    if (dist <= PLAYER_LOAD_DIST) {
                        npcIDs = append(npcIDs, npc.ID);
                        npcTypes = append(npcTypes, npc.Type);
                        npcXs = append(npcXs, npc.rect.x);
                        npcYs = append(npcYs, npc.rect.y);
                        npcHlths = append(npcHlths, npc.Health);
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
                    Action: "playerUpdate",
                    ID: player.ID,
                    Level: player.Level,
                    XP: player.XP,
                    Infamy: player.Infamy,
                    FireRate: player.FireRate,
                    Health: player.Health,
                    HealthCap: player.HealthCap,
                    HealthRegen: player.HealthRegen,
                    Energy: player.Energy,
                    EnergyCap: player.EnergyCap,
                    EnergyRegen: player.EnergyRegen,
                    Shield: player.Shield,
                    ShieldCap: player.ShieldCap,
                    ShieldRegen: player.ShieldRegen,
                    Speed: player.Speed,
                    Damage: player.Damage,
                    Scraps: player.Scraps,
                    X: player.rect.x,
                    Y: player.rect.y,
                    Rotation: player.rect.rotation,
                    Gear: player.Gear,
                    Inventory: player.Inventory,
                    IsNPC: false,
                    OtherPlayerIDs: otherPlayerIDs,
                    OtherPlayerXs: otherPlayerXs,
                    OtherPlayerYs: otherPlayerYs,
                    OtherPlayerRots: otherPlayerRots,
                    OtherPlayerHlths: otherPlayerHlths,
                    OtherPlayerGearSets: otherPlayerGearSets,
                    BulletIDs: bulletIDs,
                    BulletXs: bulletXs,
                    BulletYs: bulletYs,
                    BulletRots: bulletRots,
                    NpcIDs: npcIDs,
                    NpcTypes: npcTypes,
                    NpcXs: npcXs,
                    NpcYs: npcYs,
                    NpcHlths: npcHlths,
                }

                var newByteArray []byte
                enc := codec.NewEncoder(player.RWC, &mh)
                enc = codec.NewEncoderBytes(&newByteArray, &mh)
                enc.Encode(res1D)

                var stringMessage = string(newByteArray)
                //create header
                var header = intToBinaryString(len(stringMessage))
                //format message message
                stringMessage = header+stringMessage
                //print message
                if (player.ID == "1") {
                    //fmt.Printf(stringMessage + "\n")
                }
                // fmt.Printf("Header: %v\n", header);
                // fmt.Printf("Packet Length: %v\n", len(newByteArray));

                //send message
                go fmt.Fprint(player.RWC, stringMessage)
            } else {
                if (sentWelcomeMsg == false) {
                    sentWelcomeMsg = true
                    res1F := &Message{
                        Action: "message",
                        Data: "connected!",
                    }

                    var newByteArray []byte
                    enc := codec.NewEncoder(player.RWC, &mh)
                    enc = codec.NewEncoderBytes(&newByteArray, &mh)
                    enc.Encode(res1F)

                    var stringMessage = string(newByteArray)
                    //---create header---
                    var header = intToBinaryString(len(stringMessage))
                    //---add header---
                    stringMessage = header+stringMessage

                    stringMessage = "<?xml version=\"1.0\"?>\n<cross-domain-policy>\n<allow-access-from domain=\"*\" to-ports=\"7770-7780\"/>\n</cross-domain-policy>\n"
                    //send message
                    go fmt.Fprint(player.RWC, stringMessage)
                }
            }
        }
        time.Sleep( 1  * (time.Second / time.Duration(speedMod)) )
    }

}

func cp(w io.Writer, r io.Reader, errc chan<- error ) {
    _, err := io.Copy(w, r)
    errc <- err
}