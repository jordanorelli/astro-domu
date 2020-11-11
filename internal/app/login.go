package app

import (
	"github.com/jordanorelli/astro-domu/internal/wire"
)

type login struct {
	name string
}

func (l login) exec(ui *UI) {
	ui.client = &wire.Client{
		Log:  ui.Child("client"),
		Host: "cdm.jordanorelli.com",
		Port: 12805,
	}

	n, err := ui.client.Dial()
	if err != nil {
		panic("unable to dial server: " + err.Error())
	}
	ui.notifications = n

	res, err := ui.client.Send(wire.Login{Name: l.name})
	if err != nil {
		panic("unable to login: " + err.Error())
	}
	welcome := res.Body.(*wire.Welcome)
	ui.Info("cool beans! a login response: %#v", welcome)
	ui.state.playerName = l.name
	if ui.state.room == nil {
		ui.state.room = new(wire.Room)
	}
	p := welcome.Players[l.name]
	room := welcome.Rooms[p.Room]
	e := room.Entities[p.Avatar]
	ui.state.avatar = &e
	ui.state.room = &room

	ui.root = inGameView
}
