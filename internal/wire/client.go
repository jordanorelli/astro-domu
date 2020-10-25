package wire

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jordanorelli/blammo"
)

type Client struct {
	*blammo.Log
	Host string
	Port int

	lastSeq int
	conn    *websocket.Conn

	// outbox is the set of requests that we'd like to send. The send loop will
	// read off of this channel and write these values to the underlying
	// websocket connection.
	outbox chan Request
}

// Dial dials the server specified by the client. The returned read-only
// channel is a channel of responses from the server that are not replies to a
// request sent by the client.
func (c *Client) Dial() (<-chan Response, error) {
	c.outbox = make(chan Request)
	dialer := websocket.Dialer{
		HandshakeTimeout: 3 * time.Second,
		ReadBufferSize:   32 * 1024,
		WriteBufferSize:  32 * 1024,
		Subprotocols:     []string{"astrodomu@v0"},
	}

	path := url.URL{
		Host:   fmt.Sprintf("%s:%d", c.Host, c.Port),
		Scheme: "ws",
		Path:   "/",
	}

	c.Info("dialing: %s", path.String())
	conn, _, err := dialer.Dial(path.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("dial error: %w", err)
	}
	c.Info("connected to server")

	c.conn = conn
	done := make(chan bool, 2)
	notifications := make(chan Response)
	go c.readLoop(notifications)
	go c.writeLoop(done)
	return notifications, nil
}

func (c *Client) Send(v Value) {
	c.lastSeq++
	d := 3 * time.Second
	timeout := time.NewTimer(d)
	select {
	case c.outbox <- NewRequest(c.lastSeq, v):
		timeout.Stop()
	case <-timeout.C:
		c.Error("send timed out after %v", d)
	}
}

func (c *Client) readLoop(notifications chan<- Response) {
	defer close(notifications)

	for {
		_, r, err := c.conn.NextReader()
		if err != nil {
			c.Error("unable to get a reader frame: %v", err)
			return
		}
		b, err := ioutil.ReadAll(r)
		if err != nil {
			c.Error("unable to read frame: %v", err)
			return
		}
		var res Response
		if err := json.Unmarshal(b, &res); err != nil {
			c.Error("unable to parse frame data: %v", err)
			continue
		}
		c.Child("read-frame").Info(string(b))
		if res.Re <= 0 {
			notifications <- res
		}
	}
}

func (c *Client) writeLoop(done chan bool) {
	for {
		select {
		case req := <-c.outbox:
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				c.Error("unable to get a writer frame: %v", err)
				return
			}
			b, err := json.Marshal(req)
			if err != nil {
				c.Error("unable to marshal outgoing response: %v", err)
				break
			}
			if _, err := w.Write(b); err != nil {
				c.Error("failed to write payload: %v", err)
				break
			}
			if err := w.Close(); err != nil {
				c.Error("failed to close write frame: %v", err)
				break
			}
			c.Child("write-frame").Info(string(b))
		case shouldClose := <-done:
			if shouldClose {
				msg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
				if err := c.conn.WriteMessage(websocket.CloseMessage, msg); err != nil {
					c.Error("failed to write close message: %v", err)
				}
				c.Info("closing connection")
				if err := c.conn.Close(); err != nil {
					c.Error("failed to close connection: %v", err)
				}
				c.Info("connection closed")
				c.conn = nil
				return

			}
		}
	}

}
