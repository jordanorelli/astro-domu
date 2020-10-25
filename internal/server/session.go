package server

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jordanorelli/astro-domu/internal/wire"
	"github.com/jordanorelli/blammo"
)

type session struct {
	*blammo.Log
	id     int
	conn   *websocket.Conn
	outbox chan wire.Response
}

// pump is the session send loop. Pump should pump the session's outbox
// messages to the underlying connection until the context is closed.
func (sn *session) pump(ctx context.Context) {
	for {
		select {
		case res := <-sn.outbox:
			if err := sn.sendResponse(res); err != nil {
				sn.Error(err.Error())
			}
		case <-ctx.Done():
			sn.Info("parent context done, shutting down write pump")
			return
		}
	}
}

func (sn *session) sendResponse(res wire.Response) error {
	payload, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal outgoing response: %w", err)
	}

	if err := sn.conn.SetWriteDeadline(time.Now().Add(3 * time.Second)); err != nil {
		return fmt.Errorf("failed to set write deadline: %w", err)
	}

	w, err := sn.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return fmt.Errorf("failed get a writer frame: %w", err)
	}
	if _, err := w.Write(payload); err != nil {
		return fmt.Errorf("failed write payload: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("failed to close write frame: %w", err)
	}
	sn.Child("sent-frame").Info(string(payload))
	return nil
}
