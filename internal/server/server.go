package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/jordanorelli/astro-domu/internal/errors"
	"github.com/jordanorelli/astro-domu/internal/wire"
	"github.com/jordanorelli/blammo"
)

type Server struct {
	*blammo.Log
	Host          string
	Port          int
	lastSessionID int
}

func (s *Server) Start() error {
	if s.Host == "" {
		s.Host = "127.0.0.1"
	}
	if s.Port == 0 {
		s.Port = 12805
	}
	if s.Log == nil {
		stdout := blammo.NewLineWriter(os.Stdout)
		stderr := blammo.NewLineWriter(os.Stderr)

		options := []blammo.Option{
			blammo.DebugWriter(stdout),
			blammo.InfoWriter(stdout),
			blammo.ErrorWriter(stderr),
		}

		s.Log = blammo.NewLog("astro", options...).Child("server")
	}

	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("server failed to start a listener: %w", err)
	}
	s.Log.Info("listening for TCP traffic on %q", addr)

	go s.runHTTPServer(lis)

	return nil
}

func (s *Server) runHTTPServer(lis net.Listener) {
	err := http.Serve(lis, s)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.Error("error in http.Serve: %v", err)
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Error("upgrade error: %v", err)
		return
	}

	defer func() {
		s.Info("closing connection")
		if err := conn.Close(); err != nil {
			s.Error("error closing connection: %v", err)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.lastSessionID++
	sn := session{
		Log:    s.Log.Child("sessions").Child(strconv.Itoa(s.lastSessionID)),
		id:     s.lastSessionID,
		conn:   conn,
		outbox: make(chan wire.Response),
	}
	go sn.pump(ctx)

	for {
		t, r, err := conn.NextReader()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				s.Info("received close frame from client")
			} else {
				s.Error("read error: %v", err)
			}
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
			var req wire.Request
			if err := json.Unmarshal(text, &req); err != nil {
				s.Error("unable to parse request: %v", err)
				sn.outbox <- wire.ErrorResponse(0, "unable to parse request: %v", err)
				break
			}
			sn.outbox <- wire.NewResponse(req.Seq, wire.OK{})
		case websocket.BinaryMessage:
			sn.outbox <- wire.ErrorResponse(0, "unable to parse binary frames")
		}
	}
}
