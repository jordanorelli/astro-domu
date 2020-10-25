package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jordanorelli/blammo"
)

type client struct {
	*blammo.Log
	host    string
	port    int
	lastSeq int
	conn    *websocket.Conn
	outbox  chan request
}

func (c *client) run(ctx context.Context) {
	c.outbox = make(chan request)

	dialer := websocket.Dialer{
		HandshakeTimeout: 3 * time.Second,
		ReadBufferSize:   32 * 1024,
		WriteBufferSize:  32 * 1024,
		Subprotocols:     []string{"astrodomu@v0"},
	}

	path := url.URL{
		Host:   fmt.Sprintf("%s:%d", c.host, c.port),
		Scheme: "ws",
		Path:   "/",
	}

	conn, res, err := dialer.Dial(path.String(), nil)
	if err != nil {
		c.Error("dial error: %v", err)
		return
	}
	c.conn = conn

	c.Debug("dial response status: %d", res.StatusCode)
	for k, vals := range res.Header {
		c.Debug("dial response header: %s = %s", k, strings.Join(vals, ","))
	}

	tick := time.NewTicker(time.Second)
	defer tick.Stop()

	for {
		select {
		case req := <-c.outbox:
			payload, err := json.Marshal(req)
			if err != nil {
				c.Error("unable to marshal a request: %v", err)
				break
			}

			w, err := conn.NextWriter(websocket.TextMessage)
			if err != nil {
				c.Error("unable to get a websocket frame writer: %v", err)
				break
			}

			n, err := w.Write(payload)
			if err != nil {
				c.Error("failed to write payload of length %d: %v", len(payload), err)
				break
			}
			c.Info("wrote %d bytes for payload of length %d", n, len(payload))

			if err := w.Close(); err != nil {
				c.Error("failed to close websocket write frame: %v", err)
			}
		case <-tick.C:
			w, err := conn.NextWriter(websocket.TextMessage)
			if err != nil {
				c.Error("unable to get a websocket frame writer: %v", err)
				break
			}

			w.Write([]byte("hey"))
			w.Close()
		case <-ctx.Done():
			c.Info("parent context done, sending close message")
			msg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
			if err := conn.WriteMessage(websocket.CloseMessage, msg); err != nil {
				c.Error("failed to write close message: %v", err)
			}
			c.Info("closing connection")
			if err := conn.Close(); err != nil {
				c.Error("failed to close connection: %v", err)
			}
			c.Info("connection closed")
			c.conn = nil
			return
		}
	}
}

func (c *client) send(cmd string, args map[string]interface{}) {
	eargs := make(map[string]json.RawMessage, len(args))
	for k, v := range args {
		b, err := json.Marshal(v)
		if err != nil {
			c.Error("failed to marshal an argument: %v", err)
			return
		}
		eargs[k] = json.RawMessage(b)
	}
	c.lastSeq++
	req := request{
		Seq:     c.lastSeq,
		Command: cmd,
		Args:    eargs,
	}
	c.outbox <- req
}
