package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/websocket"

	_ "github.com/mattn/go-sqlite3"
)

type Server struct {
	addr string
}

type Socket struct {
	io.ReadWriter
	done   chan bool
	closed bool
}

type Packet struct {
	Action   string
	Username string
	Password string
	Message  string
}

func main() {

	db, err := sql.Open("sqlite3", "./Diericx.db")
	checkErr(err)
	//
	// // insert
	// stmt, err := db.Prepare("INSERT INTO users(username, password) values(?,?)")
	// checkErr(err)
	//
	// res, err := stmt.Exec("astaxie", "xxxx")
	// checkErr(err)
	//
	// id, err := res.LastInsertId()
	// checkErr(err)
	//
	// fmt.Println(id)
	// // update
	// stmt, err = db.Prepare("update users set username=? where uid=?")
	// checkErr(err)
	//
	// res, err = stmt.Exec("astaxieupdate", id)
	// checkErr(err)
	//
	// affect, err := res.RowsAffected()
	// checkErr(err)
	//
	// fmt.Println(affect)

	// query
	rows, err := db.Query("SELECT * FROM users")
	checkErr(err)

	for rows.Next() {
		var uid int
		var username string
		var password string
		err = rows.Scan(&uid, &username, &password)
		checkErr(err)
		fmt.Println(uid)
		fmt.Println(username)
		fmt.Println(password)
	}

	db.Close()

	// delete
	// stmt, err = db.Prepare("delete from users where uid=?")
	// checkErr(err)
	//
	// res, err = stmt.Exec(id)
	// checkErr(err)
	//
	// affect, err = res.RowsAffected()
	// checkErr(err)
	//
	// fmt.Println(affect)

	newServer := new(Server)
	newServer.addr = "localhost:7878"

	newServer.listenForConnections()

}

func (s Server) listenForConnections() {
	println("waiting for connection...")
	http.Handle("/", websocket.Handler(socketHandler))
	http.ListenAndServe(s.addr, nil) //192.168.2.36
}

func socketHandler(ws *websocket.Conn) {
	println("Got connection!")
	s := Socket{ws, make(chan bool), false}
	fmt.Fprint(s, "Hi!")
	//add the new player to the data stream so it will be created
	go handlePackets(s)
	<-s.done
}

func handlePackets(s Socket) {
	println("got here")
	for !s.closed {
		//serialize the packet
		buf := make([]byte, 500)
		n, err := s.Read(buf)
		if err != nil {
			println("Disconnected")
			return
		}

		var packet Packet
		//println(string(buf[0:n]))
		json.Unmarshal(buf[0:n], &packet)

		var rPacket = handlePacket(packet)
		b, _ := json.Marshal(rPacket)

		fmt.Fprintf(s, string(b))

	}
	println("stopped")
}

func handlePacket(p Packet) Packet {
	db, err := sql.Open("sqlite3", "./Diericx.db")
	checkErr(err)

	var rPacket Packet

	if p.Action == "login" {

		rPacket.Action = "login-return"

		println(p.Username)

		rows, err := db.Query("SELECT count(1) from users where (username = ?) AND (password = ?);", p.Username, p.Password)
		checkErr(err)

		for rows.Next() {
			var count int
			err = rows.Scan(&count)
			checkErr(err)
			println(count)
			if count > 0 {
				rPacket.Message = "Success"
			} else if count <= 0 {
				rPacket.Message = "Fail"
			}
		}

	}

	return rPacket

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
