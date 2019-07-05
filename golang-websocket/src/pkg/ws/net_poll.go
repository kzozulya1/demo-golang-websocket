package ws

import (
	"errors"
	"fmt"
	"net"

	"app/pkg/gpool"

	"github.com/mailru/easygo/netpoll"
)

var (
	errNetPollHandle = errors.New("Can't handle Events Read and EdgeTriggered and allocate Desc struct")
	errNetPollNew    = errors.New("Can't allocate new net poller")
)

var goRoutCounter = 0

//NetPollWait initialized net epoll object, and waits to data readiness
func NetPollWait(conn net.Conn, onDataReadyCb func() error, gRoutPool *gpool.GPool) error {

	//The Handle function creates netpoll.Desc for further use in Poller's methods:
	desc, err := netpoll.Handle(conn, netpoll.EventRead|netpoll.EventEdgeTriggered)
	if err != nil {
		return errNetPollHandle
	}

	//The Poller describes os-dependent network poller:
	poller, err := netpoll.New(nil)
	if err != nil {
		return errNetPollNew
	}

	desc = netpoll.Must(netpoll.HandleRead(conn))

	//Wait for system read readiness event
	poller.Start(desc, func(ev netpoll.Event) {
		//EventHup is indicates that some side of i/o operations (receive, send or both) is closed.
		if ev&netpoll.EventReadHup != 0 {
			poller.Stop(desc)
			conn.Close()
			return
		}

		//Data is ready, so read it with callback
		gRoutPool.Schedule(func() {
			if err := onDataReadyCb(); err != nil {
				fmt.Printf("Error in onDataReady callback: %s", err.Error())
			}
		})

	})
	return nil
}
