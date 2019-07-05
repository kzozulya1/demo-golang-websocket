package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/mailru/easygo/netpoll"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

		//Package netpoll provides a portable interface for network I/O event notification facility.
		//Its API is intended for monitoring multiple file descriptors to see if I/O is possible on any of them.

		//The Handle function creates netpoll.Desc for further use in Poller's methods:
		desc, err := netpoll.Handle(conn.UnderlyingConn(), netpoll.EventRead|netpoll.EventEdgeTriggered)
		if err != nil {
			fmt.Printf("Net poll handle error: %s\n", err.Error())
		}

		//The Poller describes os-dependent network poller:
		poller, err := netpoll.New(nil)
		if err != nil {
			fmt.Printf("Net poll new error: %s\n", err.Error())
		}

		desc = netpoll.Must(netpoll.HandleRead(conn.UnderlyingConn()))

		//fmt.Printf("before poll starts\n")
		//Read Event listening...
		poller.Start(desc, func(ev netpoll.Event) {
			//EventHup is indicates that some side of i/o operations (receive, send or
			// both) is closed.
			//So terminate current handler
			if ev&netpoll.EventReadHup != 0 {
				poller.Stop(desc)
				conn.UnderlyingConn().Close()
				return
			}

			//Websocket conn - read message
			_, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Printf("Gorilla WS: error read message: %s\n", err.Error())
			}

			//Echo message back to client
			err = conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				fmt.Printf("Gorilla WS: error write message: %s\n", err.Error())
			}

			// fmt.Printf("Net poll read all\n")

			// buf := bufio.NewReader(conn.UnderlyingConn())
			// fmt.Printf("[1]")
			// data, err := buf.ReadString('\n')
			// fmt.Printf("[1]")

			// if err != nil {
			// 	fmt.Printf("Net poll bufio read string error: %s\n", err.Error())
			// }
			// fmt.Printf("Read data: %s", data)

			// data, err := ioutil.ReadAll(conn.UnderlyingConn())
			// if err != nil {
			// 	fmt.Printf("Net poll ioutil.ReadAll error: %s\n", err.Error())
			// }
			// fmt.Printf("Read data: %q", data)

			//Echo data

			// err = conn.WriteMessage(websocket.TextMessage, []byte(data))
			// if err != nil {
			// 	fmt.Printf("Write  data back error: %s", err.Error())
			// }
		})

		//fmt.Printf("after poll start\n")

		//~~~~~~~~~~~~~~

		// for {
		// 	// Read message from browser
		// 	msgType, msg, err := conn.ReadMessage()
		// 	if err != nil {
		// 		return
		// 	}

		// 	// Print the message to the console
		// 	fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

		// 	// Write message back to browser
		// 	if err = conn.WriteMessage(msgType, msg); /*err = writeWrapper(conn, msg)*/ err != nil {
		// 		return
		// 	}
		// }
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "websockets.html")
	})

	port := os.Getenv("TCP_PORT")
	http.ListenAndServe(port, nil)
}
