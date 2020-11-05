package sim

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jordanorelli/astro-domu/internal/math"
	"github.com/jordanorelli/astro-domu/internal/wire"
	"github.com/jordanorelli/blammo"
)

type player struct {
	*blammo.Log
	//*room
	name    string
	outbox  chan wire.Response
	pending *Request
	avatar  *entity
	stop    chan bool
}

func (p *player) start(c chan Request, conn *websocket.Conn, r *room) {
	welcome := wire.Welcome{
		Rooms:   make(map[string]wire.Room),
		Players: make(map[string]wire.Player),
	}
	ents := make(map[int]wire.Entity)
	for id, e := range r.allEntities() {
		ents[id] = wire.Entity{
			ID:       id,
			Position: e.Position,
			Glyph:    e.Glyph,
		}
	}
	welcome.Rooms[r.name] = wire.Room{
		Name:     r.name,
		Rect:     r.Rect,
		Entities: ents,
	}
	for _, p := range r.players {
		wp := wire.Player{
			Name: p.name,
			Room: r.name,
		}
		if p.avatar != nil {
			wp.Avatar = p.avatar.ID
		}
		welcome.Players[p.name] = wp
	}
	p.Info("sending welcome to outbox")
	p.send(wire.Response{Re: 1, Body: welcome})
	p.Info("sent welcome, starting loops")
	p.stop = make(chan bool, 1)
	go p.readLoop(c, conn)
	go p.runLoop(conn)
}

func (p *player) readLoop(c chan Request, conn *websocket.Conn) {
	p.Info("readLoop started")
	defer p.Info("readLoop ended")

	for {
		_, b, err := conn.ReadMessage()
		if err != nil {
			p.Error("read error: %v", err)
			conn.Close()
			p.stop <- false
			return
		}
		p.Log.Child("received-frame").Info(string(b))

		var req wire.Request
		if err := json.Unmarshal(b, &req); err != nil {
			p.Error("unable to parse request: %v", err)
			continue
		}

		effect, ok := req.Body.(Effect)
		if !ok {
			p.Error("request is not an effect, is %T", req.Body)
			continue
		}
		c <- Request{
			From:  p.name,
			Seq:   req.Seq,
			Wants: effect,
		}
	}
}

func (p *player) runLoop(conn *websocket.Conn) {
	p.Info("runLoop started")
	defer p.Info("runLoop ended")

	for {
		select {
		case res := <-p.outbox:
			if err := sendResponse(conn, res); err != nil {
				p.Error(err.Error())
			}
		case sendCloseFrame := <-p.stop:
			if sendCloseFrame {
				p.Info("sending close frame")
				msg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
				if err := conn.WriteMessage(websocket.CloseMessage, msg); err != nil {
					p.Error("failed to write close message: %v", err)
				} else {
					p.Info("sent close frame")
				}
			}
			return
		}
	}
}

func (p *player) send(res wire.Response) bool {
	select {
	case p.outbox <- res:
		return true
	default:
		select {
		case <-p.outbox:
			return p.send(res)
		default:
			return false
		}
	}
}

func sendResponse(conn *websocket.Conn, res wire.Response) error {
	payload, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal outgoing response: %w", err)
	}

	if err := conn.SetWriteDeadline(time.Now().Add(3 * time.Second)); err != nil {
		return fmt.Errorf("failed to set write deadline: %w", err)
	}

	w, err := conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return fmt.Errorf("failed get a writer frame: %w", err)
	}
	if _, err := w.Write(payload); err != nil {
		return fmt.Errorf("failed write payload: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("failed to close write frame: %w", err)
	}
	return nil
}

type spawnPlayer struct{}

func (s spawnPlayer) exec(w *world, r *room, p *player, seq int) result {
	e := entity{
		ID:       <-w.nextID,
		Glyph:    '@',
		solid:    true,
		behavior: doNothing{},
	}
	p.avatar = &e

	for n, t := range r.tiles {
		x, y := n%r.Width, n/r.Width
		e.Position = math.Vec{x, y}
		if t.addEntity(&e) {
			return result{}
		}
	}
	return result{}
}

type Move math.Vec

func (Move) NetTag() string { return "move" }

func (m *Move) exec(w *world, r *room, p *player, seq int) result {
	pos := p.avatar.Position
	target := pos.Add(math.Vec(*m))
	p.Info("running move for player %s from %v to %v", p.name, p.avatar.Position, target)
	if !r.Contains(target) {
		return result{reply: wire.Errorf("target cell (%d, %d) is out of bounds", target.X, target.Y)}
	}

	currentTile := r.getTile(pos)
	nextTile := r.getTile(target)
	if !nextTile.addEntity(p.avatar) {
		return result{reply: wire.Errorf("target cell (%d, %d) is occupied", target.X, target.Y)}
	}
	currentTile.removeEntity(p.avatar.ID)
	p.avatar.Position = target
	return result{reply: wire.OK{}}
}

var lastEntityID = 0

func init() {
	wire.Register(func() wire.Value { return new(Move) })
}
