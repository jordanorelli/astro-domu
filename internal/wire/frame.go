package wire

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/jordanorelli/astro-domu/internal/math"
)

type Frame struct {
	RoomName string         `json:"room_name"`
	RoomSize math.Rect      `json:"room_size"`
	Entities map[int]Entity `json:"entities"`
	Players  map[string]int `json:"players"`
}

func (Frame) NetTag() string { return "frame" }

func (f Frame) Diff(next Frame) *Delta {
	var delta Delta

	if f.RoomSize != next.RoomSize {
		rs := next.RoomSize
		delta.RoomSize = &rs
	}

	for id, e := range next.Entities {
		if old, ok := f.Entities[id]; !ok {
			// a new entity
			delta.addEntity(e)
		} else {
			// an existing entity
			if e != old {
				delta.addEntity(e)
			}
		}
	}

	for id, _ := range f.Entities {
		// entity removed
		if _, ok := next.Entities[id]; !ok {
			delta.nullEntity(id)
		}
	}

	for name, id := range next.Players {
		if oldID, ok := f.Players[name]; !ok {
			// a new player
			delta.addPlayer(name, id)
		} else {
			// an existing player
			if oldID != id {
				delta.addPlayer(name, id)
			}
		}
	}

	for name, _ := range f.Players {
		if _, ok := next.Players[name]; !ok {
			delta.nullPlayer(name)
		}
	}
	if !delta.IsEmpty() {
		return &delta
	}
	return nil
}

type Delta struct {
	RoomSize *math.Rect      `json:"room_size,omitempty"`
	Entities map[int]*Entity `json:"entities,omitempty"`
	Players  map[string]*int `json:"players,omitempty"`
}

func (d Delta) String() string {
	first := true
	var buf bytes.Buffer
	fmt.Fprint(&buf, "Δ{")
	if d.RoomSize != nil {
		fmt.Fprintf(&buf, "size(%d,%d@%d,%d)", d.RoomSize.Origin.X, d.RoomSize.Origin.Y, d.RoomSize.Width, d.RoomSize.Height)
		first = false
	}
	if len(d.Entities) > 0 {
		if !first {
			buf.WriteString(",")
		}
		buf.WriteString("entities<")
		parts := make([]string, 0, len(d.Entities))
		for id, e := range d.Entities {
			if e == nil {
				parts = append(parts, fmt.Sprintf("-%d", id))
			} else {
				parts = append(parts, strconv.Itoa(id))
			}
		}
		buf.WriteString(strings.Join(parts, ","))
		buf.WriteString(">")
		first = false
	}
	if len(d.Players) > 0 {
		if !first {
			buf.WriteString(",")
		}
		buf.WriteString("players<")
		parts := make([]string, 0, len(d.Players))
		for name, id := range d.Players {
			if id == nil {
				parts = append(parts, fmt.Sprintf("-%s", name))
			} else {
				parts = append(parts, fmt.Sprintf("+%s", name))
			}
		}
		buf.WriteString(strings.Join(parts, ","))
		buf.WriteString(">")
		first = false
	}
	return buf.String()
}

func (Delta) NetTag() string { return "delta" }

func (d *Delta) addEntity(e Entity) {
	if d.Entities == nil {
		d.Entities = make(map[int]*Entity)
	}
	d.Entities[e.ID] = &e
}

func (d *Delta) nullEntity(id int) {
	if d.Entities == nil {
		d.Entities = make(map[int]*Entity)
	}
	d.Entities[id] = nil
}

func (d *Delta) addPlayer(name string, id int) {
	if d.Players == nil {
		d.Players = make(map[string]*int)
	}
	d.Players[name] = &id
}

func (d *Delta) nullPlayer(name string) {
	if d.Players == nil {
		d.Players = make(map[string]*int)
	}
	d.Players[name] = nil
}

func (d Delta) IsEmpty() bool {
	return d.RoomSize == nil && len(d.Entities) == 0 && len(d.Players) == 0
}

func init() {
	Register(func() Value { return new(Frame) })
	Register(func() Value { return new(Delta) })
}
