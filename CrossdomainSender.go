package main

import (
    "log"
    "net"
    "fmt"
)

var shouldListenOn843 = true

func listenOn843() {
    fmt.Printf( "<<Listening on :843 to send Crossdomain Policy>>\n" )

    listener, err := net.Listen("tcp", ":843")
    if err != nil {
        log.Fatal(err)
    }
    for shouldListenOn843 {
        c, err := listener.Accept()

        if err != nil {
            log.Fatal(err)
        }
        
        fmt.Fprint(c, "<?xml version=\"1.0\"?><cross-domain-policy><allow-access-from domain=\"*\" to-ports=\"7770-7780\"/> </cross-domain-policy>")
        
        c.Close()

    }
}