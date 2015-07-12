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
    //"strings"
    //"net/http"
    "ugorji/go/codec"
)

type player struct {
    rect rectangle
    RWC io.ReadWriteCloser
    ID string
    Health int
    Speed float64
    Scraps int32
    X float64
    Y float64
    xMovement float64
    yMovement float64
    gear gearSet
}

type point struct {
    x float64
    y float64
}

type gearSet struct {
    cockpit int
    lasers int
    wings int
    jets int
}

type Enemy struct {
    rect rectangle
    ID string
    health int
    x float64
    y float64
    rotation int
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
    ID string
    Health int
    Scraps int32
    X float64 //change to float64
    Y float64
    Rotation int
    Gear []int
    IsNPC bool
    //other player data
    OtherPlayerIDs []string
    OtherPlayerXs []float64
    OtherPlayerYs []float64
    OtherPlayerHlths []int
    //bullet data
    BulletIDs []int
    BulletXs []float64
    BulletYs []float64
    BulletRots []int
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

type EnemyUpdate struct {
    Action string
    ID string
    Health int
    X float64
    Y float64
    Rotation string
    IsNPC bool
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

const listenAddr = "192.168.0.167:7777"

const baseAddr = "http://192.168.1.18:3000/api/v1/"

var (
    mh codec.MsgpackHandle
)

var partner = make(chan io.ReadWriteCloser)

var players []*player
var bullets []*bullet
var enemies []*Enemy

var shouldQuit = false

//CONSTANTS
var PLAYER_LOAD_DIST float64 = 20
var ARENA_SIZE float64 = 100

func main() {
    rand.Seed(time.Now().Unix())

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
    go moveBullets()
    go moveEnemies()
    go movePlayers()
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

func sendEnemiesInArea(p *player) {
    for _, e := range enemies {
        packet := &EnemyUpdate{
            Action: "update",
            ID: e.ID,
            Health: e.health,
            X: e.rect.x,
            Y: e.rect.y,
            Rotation: strconv.Itoa(e.rotation),
            IsNPC: true,
        }

        var newByteArray []byte
        enc := codec.NewEncoder(p.RWC, &mh)
        enc = codec.NewEncoderBytes(&newByteArray, &mh)
        enc.Encode(packet)

        var stringMessage = string(newByteArray)
        stringMessage += "\n"
        var diff = 100 - len(stringMessage)

        for i := 1; i < diff; i ++ {
            stringMessage += "$"
        } 

        go fmt.Fprintln(p.RWC, stringMessage)
    }
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

func movePlayers() {
    for {
        for _, player := range players {
            player.rect.x = player.rect.x + (player.xMovement*player.Speed)
            player.rect.y = player.rect.y + (player.yMovement*player.Speed)


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
        
            // rotateRectsPoints(player.rect, (float64(player.Rotation) / 180.0) * 3.14159 )
            
                // fmt.Printf("\n %v, %v", bullet.rect.x, bullet.rect.y)
            for _, player := range players {
                if ( compareRects(player.rect, bullet.rect) == true && bullet.shooter != player ) {
                    fmt.Printf("BULLET HIT Player")
                    var toRemove int = -1
                    for i, bullet2 := range bullets {
                        if (bullet == bullet2) {
                            toRemove = i
                        }
                    }
                    bullets[toRemove] = bullets[len(bullets)-1]
                    bullets = bullets[0:len(bullets)-1]

                    player.Health = player.Health - 10

                    if (player.Health <= 0) {
                        player.rect.x = 0
                        player.rect.y = 0
                        player.Health = 100

                        //update shooter's scraps
                        bullet.shooter.Scraps += 100
                    }

                }
            
            }

            // for _, e := range enemies { 
            //     if ( compareRects(e.rect, bullet.rect) == true ) {
            //         fmt.Printf("BULLET HIT Enemy ")
            //         var toRemove int = -1
            //         for i, bullet2 := range bullets {
            //             if (bullet == bullet2) {
            //                 toRemove = i
            //             }
            //         }
            //         bullets[toRemove] = bullets[len(bullets)-1]
            //         bullets = bullets[0:len(bullets)-1]

            //         e.health = e.health - 10

            //         if (e.health <= 0) {
            //             var toRemove int = -1
            //             for i, e2 := range enemies {
            //                 if (e == e2) {
            //                     toRemove = i
            //                 }
            //             }
            //             enemies[toRemove] = enemies[len(enemies)-1]
            //             enemies = enemies[0:len(enemies)-1]    
            //         }
            //     }
            // }
        }

        time.Sleep( (time.Second / time.Duration(60)) )
    }
}

func moveEnemies() {
    for {

        for _, e := range enemies {
            e.rect.x = e.rect.x + ((rand.Float64() * 2) - 1)
            e.rect.y = e.rect.y + ((rand.Float64() * 2) - 1)
        }

        time.Sleep( (time.Second / time.Duration(10)) )
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
    // var newPlayer player
    newPlayer := new(player)
    newPlayer.RWC = c
    newPlayer.ID = ""
    newPlayer.Speed = 0.03
    newPlayer.Health = 100
    newPlayer.Scraps = 0

    newPlayer.rect = createRect(0, 0, 1, 1)

    newPlayer.gear.cockpit = -1

    players = append(players, newPlayer)
    fmt.Printf("APPENDING PLAYER\n")
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

func removePlayerFromList(p *player) {
    var i = 0;
    var foundIndex = -1;
    for _, player := range players {
        if (p == player) {
            foundIndex = i;
        }
        i++
    }
    if (foundIndex != -1) {
        players = append(players[:foundIndex], players[foundIndex+1:]...)
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

                player.xMovement = res.X
                player.yMovement = res.Y

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
            } else if res.Action == "shoot" {
                player.Health = player.Health
                res1D := &Shoot{
                    Action: "shoot",
                    ID: player.ID,
                    X: res.X,
                    Y: res.Y,
                    Rotation: res.Rotation,
                }  

                var newByteArray []byte
                enc := codec.NewEncoder(player.RWC, &mh)
                enc = codec.NewEncoderBytes(&newByteArray, &mh)
                enc.Encode(res1D)

                //fmt.Println(string(newByteArray))

                var resX float64
                var resY float64

                resX = res.X
                resY = res.Y

                // var newBullet bullet
                newBullet := new (bullet)
                newBullet.ID = rand.Intn(1000)
                newBullet.rect = createRect(resX, resY, 0.17, 0.5)
                newBullet.rect.rotation = res.Rotation
                newBullet.shooter = player

                bullets = append(bullets, newBullet)            
            } else if res.Action == "input" {
                // player.rect.y = player.rect.y + (0.
                // if (res.Y == 1) {
                //     player.rect.y = player.rect.y + 0.5
                // }
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
    for { 

        bulletIDs := make([]int, 0);
        bulletXs := make([]float64, 0);
        bulletYs := make([]float64, 0);
        bulletRots := make([]int, 0);

        otherPlayerIDs := make([]string, 0);
        otherPlayerXs := make([]float64, 0);
        otherPlayerYs := make([]float64, 0);
        otherPlayerHlths := make([]int, 0);

        //var bulletPackets []*BulletUpdate

        for _, player := range players {
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
                        fmt.Printf("X: %v, Y: %v\n", otherPlayer.rect.x, otherPlayer.rect.y)
                    }

                    var dist = math.Sqrt( math.Pow(otherPlayer.rect.x - player.rect.x, 2) + math.Pow(otherPlayer.rect.y - player.rect.y, 2) )
                    if (dist <= PLAYER_LOAD_DIST && player != otherPlayer) {
                        otherPlayerIDs = append(otherPlayerIDs, otherPlayer.ID);
                        otherPlayerXs = append(otherPlayerXs, otherPlayer.rect.x);
                        otherPlayerYs = append(otherPlayerYs, otherPlayer.rect.y);
                        otherPlayerHlths = append(otherPlayerHlths, otherPlayer.Health);
                    }
                }

                //create new gear obj for the other players current gear set
                // otherPlayersGear := gearSet{
                //     cockpit: otherPlayer.gear.cockpit,
                //     lasers: otherPlayer.gear.lasers,
                //     wings: otherPlayer.gear.wings,
                //     jets: otherPlayer.gear.jets,
                // }

                gear := []int{player.gear.cockpit, player.gear.lasers, player.gear.wings, player.gear.jets}

                //new table that has multiple updates 

                //create update packet
                res1D := &Update{
                    Action: "playerUpdate",
                    ID: player.ID,
                    Health: player.Health,
                    Scraps: player.Scraps,
                    X: player.rect.x,
                    Y: player.rect.y,
                    Rotation: player.rect.rotation,
                    Gear: gear,
                    IsNPC: false,
                    OtherPlayerIDs: otherPlayerIDs,
                    OtherPlayerXs: otherPlayerXs,
                    OtherPlayerYs: otherPlayerYs,
                    OtherPlayerHlths: otherPlayerHlths,
                    BulletIDs: bulletIDs,
                    BulletXs: bulletXs,
                    BulletYs: bulletYs,
                    BulletRots: bulletRots,
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
                //send message
                go fmt.Fprint(player.RWC, stringMessage)
            }
        }
        time.Sleep( 1  * (time.Second / time.Duration(speedMod)) )
    }

}

func cp(w io.Writer, r io.Reader, errc chan<- error ) {
    _, err := io.Copy(w, r)
    errc <- err
}