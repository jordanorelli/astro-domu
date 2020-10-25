package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/jordanorelli/blammo"
)

type server struct {
	*blammo.Log
	host          string
	port          int
	lastSessionID int
}

func (s *server) listen() error {
	return http.ListenAndServe(fmt.Sprintf("%s:%d", s.host, s.port), s)
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Error("upgrade error: %v", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.lastSessionID++
	sn := session{
		Log:    s.Log.Child("sessions").Child(strconv.Itoa(s.lastSessionID)),
		id:     s.lastSessionID,
		conn:   conn,
		outbox: make(chan response),
	}
	go sn.pump(ctx)

	for {
		t, r, err := conn.NextReader()
		if err != nil {
			s.Error("read error: %v", err)
			return
		}

		switch t {
		case websocket.TextMessage:
			text, err := ioutil.ReadAll(r)
			if err != nil {
				s.Error("readall error: %v", err)
				break
			}
			sn.Log.Child("received-frame").Info(string(text))
			var body requestBody
			if err := json.Unmarshal(text, &body); err != nil {
				s.Error("unable to parse request: %v", err)
				sn.outbox <- errorResponse(0, fmt.Errorf("unable to parse request: %v", err))
				break
			}
			sn.outbox <- ok(body.Seq)
		case websocket.BinaryMessage:
			sn.outbox <- errorResponse(0, fmt.Errorf("unable to parse binary frames"))
		}
	}
}
