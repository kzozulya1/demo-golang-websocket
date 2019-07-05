package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"app/pkg/gpool"
	wsConn "app/pkg/ws"

	"github.com/gobwas/ws"
)

const (
	maxReadGORoutines = 5
)

func main() {
	port := os.Getenv("TCP_PORT")
	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	//Allocate go routine pool
	pool := gpool.New(maxReadGORoutines)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Error accept connection: %s", err.Error())
			continue
		}
		_, err = ws.Upgrade(conn)
		if err != nil {
			fmt.Printf("Error upgrade connection %s", err.Error())
			continue
		}

		go wsConn.HandleWSConnection(conn, pool)
	}
}
