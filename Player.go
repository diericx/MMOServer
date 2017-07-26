package main

import "net"

var players = make(map[string]Player)

var redTeamPlayerCount = 0
var blueTeamPlayerCount = 0

//Player player struct
type Player struct {
	addr *net.UDPAddr
	e    *Entity
}

//NewPlayer Create new player object
func NewPlayer(addr *net.UDPAddr) *Player {

	p := Player{
		addr: addr,
		e:    NewEntity("player", true),
	}

	if redTeamPlayerCount <= blueTeamPlayerCount {
		p.e.Team = true //True = redTeam
		redTeamPlayerCount++
	} else {
		p.e.Team = false
		blueTeamPlayerCount++
	}

	players[addr.String()] = p

	return &p
}

//GetPlayer Returns the player found at this location or nil
func GetPlayer(addr *net.UDPAddr) *Player {
	var p = players[addr.String()]

	if p.addr.String() == addr.String() {
		return &p
	}

	return nil
}
