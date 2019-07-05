package ws

import (
	"errors"
	"fmt"
	"io"
	"net"

	"app/pkg/gpool"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

//ReadHandler maintain state and read/write ops
type ReadHandler struct {
	state  ws.State
	reader *wsutil.Reader
	writer *wsutil.Writer
	conn   net.Conn
}

//NewReadHandler constructs new WsReadHandler object
func NewReadHandler(conn net.Conn) *ReadHandler {
	return &ReadHandler{
		ws.StateServerSide,
		wsutil.NewReader(conn, ws.StateServerSide),
		wsutil.NewWriter(conn, ws.StateServerSide, ws.OpText),
		conn,
	}
}

//OnDataReady callback - called when net poll is waked by new data in websocket
func (wrh *ReadHandler) OnDataReady() error {
	header, err := wrh.reader.NextFrame()
	if err != nil {
		if gpool.HTTPisClosedConnError(err) {
			return errors.New("connection is closed")
		}
		return fmt.Errorf("Next frame read error: %s", err.Error())
	}

	// Reset writer to write frame with right operation code.
	wrh.writer.Reset(wrh.conn, wrh.state, header.OpCode)

	if _, err = io.Copy(wrh.writer, wrh.reader); err != nil {
		return fmt.Errorf("Io copy to writer error: %s", err.Error())
	}

	if err = wrh.writer.Flush(); err != nil {
		return fmt.Errorf("Writer flush error: %s", err.Error())
	}
	return nil
}
