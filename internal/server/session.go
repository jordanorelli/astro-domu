package server

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jordanorelli/astro-domu/internal/wire"
	"github.com/jordanorelli/blammo"
)

type session struct {
	*blammo.Log
	id       int
	entityID int
	start    time.Time
	conn     *websocket.Conn
	outbox   chan wire.Response
	done     chan bool
}

// run is the session run loop.
func (sn *session) run() {
	for {
		select {
		case res := <-sn.outbox:
			if err := sn.sendResponse(res); err != nil {
				sn.Error(err.Error())
			}
		case sendCloseFrame := <-sn.done:
			sn.Info("saw done signal")
			if sendCloseFrame {
				sn.Info("sending close frame")
				msg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
				if err := sn.conn.WriteMessage(websocket.CloseMessage, msg); err != nil {
					sn.Error("failed to write close message: %v", err)
				} else {
					sn.Info("sent close frame")
				}
			}
			return
		}
	}
}

// read reads for messages on the underlying websocket.
func (sn *session) read() {
	for {
		t, b, err := sn.conn.ReadMessage()
		if err != nil {
			if v, ok := err.(*websocket.CloseError); ok {
				switch v.Code {
				case websocket.CloseNormalClosure,
					websocket.CloseGoingAway:
					sn.Info("received close frame with code %d (%v)", v.Code, err)
				case websocket.CloseNoStatusReceived,
					websocket.CloseProtocolError,
					websocket.CloseUnsupportedData,
					websocket.CloseAbnormalClosure,
					websocket.CloseInvalidFramePayloadData,
					websocket.ClosePolicyViolation,
					websocket.CloseMessageTooBig,
					websocket.CloseMandatoryExtension,
					websocket.CloseInternalServerErr,
					websocket.CloseServiceRestart,
					websocket.CloseTryAgainLater,
					websocket.CloseTLSHandshake:
					sn.Error("received close frame with code %d (%v)", v.Code, err)
				default:
					sn.Error("received close frame with code %d (%v)", v.Code, err)
				}
				return
			}
			sn.Error("unexpected read error: %v", err)
			return
		}

		switch t {
		case websocket.TextMessage:
			sn.Log.Child("received-frame").Info(string(b))
			var req wire.Request
			if err := json.Unmarshal(b, &req); err != nil {
				sn.Error("unable to parse request: %v", err)
				sn.outbox <- wire.ErrorResponse(0, "unable to parse request: %v", err)
				break
			}
			sn.outbox <- wire.NewResponse(req.Seq, wire.OK{})
		case websocket.BinaryMessage:
			sn.outbox <- wire.ErrorResponse(0, "unable to parse binary frames")
		}
	}
}

// sendResponse sends an individual response on the underlying websocket connection
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
