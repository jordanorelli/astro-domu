package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jordanorelli/blammo"
)

type server struct {
	*blammo.Log
	host  string
	port  int
	world *room
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
			s.Info("received: %s", text)
		case websocket.BinaryMessage:

		}
	}
}
