package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jordanorelli/astro-domu/internal/wire"
	"github.com/jordanorelli/blammo"
)

type client struct {
	*blammo.Log
	host    string
	port    int
	lastSeq int
	conn    *websocket.Conn
	outbox  chan wire.Request
}

func (c *client) run(ctx context.Context) {
	c.outbox = make(chan wire.Request)

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
	go c.readLoop()

	c.Debug("dial response status: %d", res.StatusCode)
	for k, vals := range res.Header {
		c.Debug("dial response header: %s = %s", k, strings.Join(vals, ","))
	}

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

			if _, err := w.Write(payload); err != nil {
				c.Error("failed to write payload of length %d: %v", len(payload), err)
				break
			}
			c.Child("sent-frame").Info(string(payload))

			if err := w.Close(); err != nil {
				c.Error("failed to close websocket write frame: %v", err)
			}
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

func (c *client) send(v wire.Value) {
	c.lastSeq++
	c.outbox <- wire.NewRequest(c.lastSeq, v)
}

func (c *client) readLoop() {
	for {
		_, r, err := c.conn.NextReader()
		if err != nil {
			return
		}
		b, err := ioutil.ReadAll(r)
		if err != nil {
			return
		}
		c.Log.Child("received-frame").Info(string(b))
	}
}
