package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jordanorelli/blammo"
)

// session represents the server side of a client's session. i.e., a single
// connection along with its associated state.
type session struct {
	*blammo.Log
	id     int
	conn   *websocket.Conn
	outbox chan response
}

// pump is the session send loop. Pump should pump the session's outbox
// messages to the underlying connection until the context is closed.
func (sn *session) pump(ctx context.Context) {
	for {
		select {
		case res := <-sn.outbox:
			payload, err := json.Marshal(res)
			if err != nil {
				sn.Error("failed to marshal outgoing response: %v", err)
				break
			}

			if err := sn.conn.SetWriteDeadline(time.Now().Add(3 * time.Second)); err != nil {
				sn.Error("failed to set write deadline: %v", err)
				break
			}
			w, err := sn.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				sn.Error("failed get a writer frame: %v", err)
				break
			}
			if _, err := w.Write(payload); err != nil {
				sn.Error("failed write payload: %v", err)
				break
			}
			if err := w.Close(); err != nil {
				sn.Error("failed to close write frame: %v", err)
				break
			}
			sn.Child("sent-frame").Info(string(payload))
		case <-ctx.Done():
			sn.Info("parent context done, shutting down write pump")
			return
		}
	}
}
