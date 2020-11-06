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

	conn *websocket.Conn

	outbox   chan *pending
	resolved chan Response
	done     chan bool
}

// Dial dials the server specified by the client. The returned read-only
// channel is a channel of responses from the server that are not replies to a
// request sent by the client.
func (c *Client) Dial() (<-chan Response, error) {
	c.outbox = make(chan *pending)
	c.resolved = make(chan Response)
	c.done = make(chan bool, 1)

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
	notifications := make(chan Response)
	go c.readLoop(notifications)
	go c.writeLoop()
	return notifications, nil
}

func (c *Client) Send(v Value) (Response, error) {
	d := 3 * time.Second
	timeout := time.NewTimer(d)

	done := make(chan struct{})
	p := pending{v: v, done: done}

	select {
	case c.outbox <- &p:
		timeout.Stop()
	case <-timeout.C:
		return Response{}, fmt.Errorf("send timed out after %v", d)
	}

	select {
	case <-done:
		if p.err == nil {
			if err, ok := p.res.Body.(error); ok {
				return p.res, err
			}
		}
		return p.res, p.err
	case <-timeout.C:
		return Response{}, fmt.Errorf("send timed out (2) after %v", d)
	}
}

func (c *Client) Close() { c.done <- true }

func (c *Client) readLoop(notifications chan<- Response) {
	defer close(notifications)

	for {
		c.Info("waiting for a reader frame")
		_, r, err := c.conn.NextReader()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				c.Info("received close frame from server")
				break
			}
			c.Error("unable to get a reader frame: %v", err)
			break
		}
		b, err := ioutil.ReadAll(r)
		if err != nil {
			c.Error("unable to read frame: %v", err)
			continue
		}
		var res Response
		if err := json.Unmarshal(b, &res); err != nil {
			c.Error("unable to parse frame data: %v", err)
			continue
		}
		c.Child("read-frame").Info(string(b))
		if res.Re <= 0 {
			notifications <- res
		} else {
			c.resolved <- res
		}
	}
	c.done <- false

	c.Info("closing connection")
	if err := c.conn.Close(); err != nil {
		c.Error("error on connection close: %v", err)
	} else {
		c.Info("closed connection cleanly")
	}
}

func (c *Client) writeLoop() {
	sent := make(map[int]*pending)

	for seq := 1; true; seq++ {
		select {
		case p := <-c.outbox:
			req := Request{seq, p.v}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				p.err = fmt.Errorf("unable to get a writer frame: %w", err)
				close(p.done)
				return
			}
			b, err := json.Marshal(req)
			if err != nil {
				p.err = fmt.Errorf("unable to marshal outgoing response: %w", err)
				close(p.done)
				break
			}
			if _, err := w.Write(b); err != nil {
				p.err = fmt.Errorf("failed to write payload: %w", err)
				close(p.done)
				break
			}
			if err := w.Close(); err != nil {
				p.err = fmt.Errorf("failed to close write frame: %w", err)
				close(p.done)
				break
			}
			c.Child("write-frame").Info(string(b))
			sent[seq] = p

		case res := <-c.resolved:
			p, ok := sent[res.Re]
			if !ok {
				c.Error("saw response for unknown seq %d")
				break
			}
			delete(sent, res.Re)
			p.res = res
			close(p.done)

		case shouldClose := <-c.done:
			c.Info("write loop sees done signal: %t", shouldClose)

			if shouldClose {
				c.Info("sending close frame")
				msg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
				if err := c.conn.WriteMessage(websocket.CloseMessage, msg); err != nil {
					c.Error("failed to write close message: %v", err)
				} else {
					c.Info("sent close frame")
				}
			}
			return
		}
	}
}

type pending struct {
	v    Value
	res  Response
	err  error
	done chan struct{}
}
