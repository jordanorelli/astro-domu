package main

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jordanorelli/blammo"
)

type client struct {
	*blammo.Log
	host string
	port int
}

func (c *client) run(ctx context.Context) {
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

	c.Debug("dial response status: %d", res.StatusCode)
	for k, vals := range res.Header {
		c.Debug("dial response header: %s = %s", k, strings.Join(vals, ","))
	}

	tick := time.NewTicker(time.Second)
	defer tick.Stop()

	for {
		select {
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
			return
		}
	}
}
