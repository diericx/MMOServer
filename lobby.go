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
    "encoding/json"
    //"net/http"
    "ugorji/go/codec"
)

type player struct {
    rect rectangle
    RWC io.ReadWriteCloser
    ID string
    Health int
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
    ID string
    rect rectangle
}

type Message struct {
    Action string
    Data string
}

type Update struct {
    Action string
    ID string
    Health string
    X float64 //change to float64
    Y float64
    Rotation int
    Gear []int
    IsNPC bool
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

func sendToEveryoneBut(toSend string, rwc io.ReadWriteCloser) {
    for _, player := range players {
        if (player.RWC != rwc) {
            go fmt.Fprintln(player.RWC, toSend )
        }
    }
}

func sendToEveryone(toSend string) {
    for _, player := range players {
        go fmt.Fprintln(player.RWC, toSend )
    }
}

func sendDamageTakenPacket(p *player) {
    packet := &DamageTaken{
        ID: p.ID,
        Action: "damageTaken",
        BulletID: 1,
    }
    // fmt.Printf(p.ID)
    packetString, _ := json.Marshal(packet)
    go sendToEveryone(string(packetString))
    //go fmt.Fprintln(p.RWC, string(packetString) )
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

func moveBullets() {
    for {

        for _, bullet := range bullets {
            var bulletRadians float64 = (float64(bullet.rect.rotation+90) / 180.0) * 3.14159
            bullet.rect.x = bullet.rect.x + (15 * 0.016 * math.Cos( bulletRadians ) )
            bullet.rect.y = bullet.rect.y + (15 * 0.016 * math.Sin( bulletRadians ) )
        }

        for _, player := range players {
            player.rect.x = player.rect.x + (player.xMovement*0.1)
            player.rect.y = player.rect.y + (player.yMovement*0.1)
            //player.rect.rotation = player.Rotation
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
                    }

                    // sendDamageTakenPacket(player)
                }
            
            }

            for _, e := range enemies { 
                if ( compareRects(e.rect, bullet.rect) == true ) {
                    fmt.Printf("BULLET HIT Enemy ")
                    var toRemove int = -1
                    for i, bullet2 := range bullets {
                        if (bullet == bullet2) {
                            toRemove = i
                        }
                    }
                    bullets[toRemove] = bullets[len(bullets)-1]
                    bullets = bullets[0:len(bullets)-1]

                    e.health = e.health - 10

                    if (e.health <= 0) {
                        var toRemove int = -1
                        for i, e2 := range enemies {
                            if (e == e2) {
                                toRemove = i
                            }
                        }
                        enemies[toRemove] = enemies[len(enemies)-1]
                        enemies = enemies[0:len(enemies)-1]    
                    }
                }
            }
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
    newPlayer.Health = 100

    newPlayer.rect = createRect(0, 0, 1, 1)

    newPlayer.gear.cockpit = -1

    players = append(players, newPlayer)
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

func getDataFromPlayer(player *player) {

    for {
        var shouldRemove = false

        buf := make([]byte, 1024)
        n, err := player.RWC.Read(buf)

        // var stringData = string(buf[0:n])
        // fmt.Printf(stringData)
        // dec := json.NewDecoder(strings.NewReader(stringData))

        var res = &Update{}
        // dec.Decode(&res)

        dec := codec.NewDecoder(player.RWC, &mh)
        dec = codec.NewDecoderBytes(buf[0:n], &mh)
        dec.Decode(res)

        if err == nil {
            // fmt.Printf("%v", res.X)
            // res := &Update{}

            //decoder.Decode(n)
            // json.Unmarshal([]byte(buf[0:n]), &res)
            // fmt.Printf("\n x=%v", res.X )
            if res.Action == "update" {

                player.ID = res.ID

                player.xMovement = res.X
                player.yMovement = res.Y

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

                // fmt.Println(string(newByteArray))
                //sendToEveryoneBut( string(newByteArray), player.RWC ) 

                var resX float64
                var resY float64

                resX = res.X
                resY = res.Y

                // var newBullet bullet
                newBullet := new (bullet)
                newBullet.ID = strconv.Itoa( rand.Intn(1000) )
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

// {"Action":"shoot", "ID":"87", "X":"0", "Y":"0", "Rotation":"22"}
func chat () {

    var speedMod = 30
    for { 
        for _, player := range players {
            //fmt.Printf("%v", player.ID)

            for _, bullet := range bullets {
                //create update packet
                gear := []int{0, 0, 0, 0}
                res1D := &Update{
                    Action: "bulletUpdate",
                    ID: bullet.ID,
                    Health: strconv.Itoa(0),
                    X: bullet.rect.x,
                    Y: bullet.rect.y,
                    Rotation: bullet.rect.rotation,
                    Gear: gear,
                    IsNPC: false,
                }         

                var newByteArray []byte
                enc := codec.NewEncoder(player.RWC, &mh)
                enc = codec.NewEncoderBytes(&newByteArray, &mh)
                enc.Encode(res1D)

                var stringMessage = string(newByteArray)
                var diff = 100 - len(stringMessage)

                for i := 1; i < diff; i ++ {
                    stringMessage += "$"
                }  

                go fmt.Fprintln(player.RWC, stringMessage)  
            }

            // sendEnemiesInArea(player)

            for _, otherPlayer := range players {
                if (otherPlayer.ID != "") {
                    //create new gear obj for the other players current gear set
                    // otherPlayersGear := gearSet{
                    //     cockpit: otherPlayer.gear.cockpit,
                    //     lasers: otherPlayer.gear.lasers,
                    //     wings: otherPlayer.gear.wings,
                    //     jets: otherPlayer.gear.jets,
                    // }

                    gear := []int{otherPlayer.gear.cockpit, otherPlayer.gear.lasers, otherPlayer.gear.wings, otherPlayer.gear.jets}

                    //new table that has multiple updates 

                    //create update packet
                    res1D := &Update{
                        Action: "playerUpdate",
                        ID: otherPlayer.ID,
                        Health: strconv.Itoa(otherPlayer.Health),
                        X: otherPlayer.rect.x,
                        Y: otherPlayer.rect.y,
                        Rotation: otherPlayer.rect.rotation,
                        Gear: gear,
                        IsNPC: false,
                    }

                    var newByteArray []byte
                    enc := codec.NewEncoder(player.RWC, &mh)
                    enc = codec.NewEncoderBytes(&newByteArray, &mh)
                    enc.Encode(res1D)

                    var stringMessage = string(newByteArray)
                    var diff = 100 - len(stringMessage)

                    for i := 1; i < diff; i ++ {
                        stringMessage += "$"
                    } 

                    // dec := codec.NewDecoder(player.RWC, &mh)
                    // dec = codec.NewDecoderBytes(b, &mh)
                    // err2 := dec.Decode(&Update) 
                    //send update packet
                    // res1B, _ := json.Marshal(res1D)
                    // go fmt.Fprintln(player.RWC, string(res1B)+"\n" )
                    // fmt.Println( string(b) )
                    // fmt.Printf(strconv.Itoa(len(stringMessage) ))
                    go fmt.Fprintln(player.RWC, stringMessage)
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
                    var diff = 100 - len(stringMessage)

                    for i := 1; i < diff; i ++ {
                        stringMessage += "$"
                    } 

                    go fmt.Fprintln(player.RWC, stringMessage)
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