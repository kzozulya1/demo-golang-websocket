package ws

import (
	"app/pkg/gpool"
	"net"
)

//HandleWSConnection  handles new ws connection
func HandleWSConnection(conn net.Conn, pool *gpool.GPool) {
	//init ws read hanldler
	rHandler := NewReadHandler(conn)
	//wait for data is ready for read in conn. After ready - invoke callback rHandler.OnDataReady
	NetPollWait(conn, rHandler.OnDataReady, pool)
}
