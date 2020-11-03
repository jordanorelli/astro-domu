package sim

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jordanorelli/astro-domu/internal/errors"
	"github.com/jordanorelli/astro-domu/internal/wire"
	"github.com/jordanorelli/blammo"
)

type Server struct {
	*blammo.Log
	http  *http.Server
	world *world
}

func (s *Server) Start(host string, port int) error {
	if s.Log == nil {
		s.Log = defaultLog().Child("server")
	}

	s.world = newWorld(s.Log.Child("world"))
	go s.world.run(2)

	addr := fmt.Sprintf("%s:%d", host, port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("server failed to start a listener: %w", err)
	}
	s.Log.Info("listening for TCP traffic on %q", addr)

	go s.runHTTPServer(lis)

	return nil
}

func (s *Server) runHTTPServer(lis net.Listener) {
	srv := http.Server{
		Handler: s,
	}
	s.http = &srv
	err := srv.Serve(lis)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.Error("error in http.Serve: %v", err)
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log := s.Log.Child("login")
	upgrader := websocket.Upgrader{
		HandshakeTimeout: 3 * time.Second,
		ReadBufferSize:   2 << 12,
		WriteBufferSize:  2 << 12,
		Subprotocols:     []string{"astrodomu@v0"},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("upgrade error: %v", err)
		return
	}

	t, rd, err := conn.NextReader()
	if err != nil {
		log.Error("unable to get a reader: %v", err)
		conn.Close()
		return
	}

	if t != websocket.TextMessage {
		log.Error("first message is not text")
		// TODO: send websocket close frame here
		conn.Close()
		return
	}

	var req wire.Request
	if err := json.NewDecoder(rd).Decode(&req); err != nil {
		log.Error("unable to parse initial request: %v", err)
		// TODO: send websocket close frame here
		conn.Close()
		return
	}

	login, ok := req.Body.(*wire.Login)
	if !ok {
		log.Error("first request is not wire.Login, is %T", req.Body)
		// TODO: send websocket close frame here
		conn.Close()
		return
	}

	log.Info("login requested: %v", *login)

	failed := make(chan error, 1)
	s.world.connect <- connect{
		conn:   conn,
		login:  *login,
		failed: failed,
	}
	e := <-failed
	if e != nil {
		log.Error("connect failed: %v", err)
		// TODO: send websocket close frame here
		conn.Close()
		return
	}
}

func (s *Server) Shutdown() {
	s.Info("starting shutdown procedure")

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()

		if err := s.world.stop(); err != nil {
			s.Error("error stopping the simulation: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		log := s.Child("http")
		log.Info("shutting down http server")
		if err := s.http.Shutdown(context.Background()); err != nil {
			log.Error("error shutting down http server: %v", err)
		} else {
			log.Info("http server has shut down")
		}
	}()

	wg.Wait()
	s.Info("shutdown procedure complete")
}
