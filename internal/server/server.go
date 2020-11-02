package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jordanorelli/astro-domu/internal/errors"
	"github.com/jordanorelli/astro-domu/internal/server/sim"
	"github.com/jordanorelli/astro-domu/internal/wire"
	"github.com/jordanorelli/blammo"
)

type Server struct {
	*blammo.Log
	http  *http.Server
	world *sim.World

	sync.Mutex
	lastSessionID  int
	sessions       map[int]*session
	waitOnSessions sync.WaitGroup
}

func (s *Server) Start(host string, port int) error {
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

	s.world = sim.NewWorld(s.Log.Child("world"))
	go s.world.Run(3)

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

func (s *Server) createSession(conn *websocket.Conn) *session {
	s.Lock()
	defer s.Unlock()

	s.lastSessionID++
	sn := &session{
		Log:    s.Log.Child("sessions").Child(strconv.Itoa(s.lastSessionID)),
		id:     s.lastSessionID,
		start:  time.Now(),
		conn:   conn,
		outbox: make(chan wire.Response),
		done:   make(chan bool, 1),
	}
	if s.sessions == nil {
		s.sessions = make(map[int]*session)
	}
	s.waitOnSessions.Add(1)
	s.sessions[sn.id] = sn
	// sn.entityID = s.world.SpawnPlayer(sn.id)
	// s.Info("created session %d, %d sessions active", sn.id, len(s.sessions))
	return sn
}

// dropSession removes a session from the server. This should only be called as
// a result of the connection's read loop terminating
func (s *Server) dropSession(sn *session) {
	s.Lock()
	defer s.Unlock()

	close(sn.done)
	delete(s.sessions, sn.id)
	s.waitOnSessions.Add(-1)

	s.Info("dropped session %d after %v time connected, %d sessions active", sn.id, time.Since(sn.start), len(s.sessions))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Error("upgrade error: %v", err)
		return
	}

	sn := s.createSession(conn)
	go sn.run()
	sn.read(s.world.Inbox)
	s.dropSession(sn)

	sn.Info("closing connection")
	if err := conn.Close(); err != nil {
		s.Error("error closing connection: %v", err)
	}
}

func (s *Server) Shutdown() {
	s.Info("starting shutdown procedure")

	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()

		if err := s.world.Stop(); err != nil {
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

	go func() {
		defer wg.Done()

		log := s.Child("sessions")
		s.Lock()
		numSessions := len(s.sessions)
		if numSessions > 0 {
			log.Info("broadcasting shutdown to %d active sessions", numSessions)
			for id, sn := range s.sessions {
				log.Info("sending done signal to session: %d", id)
				sn.done <- true
			}
		} else {
			log.Info("no active sessions")
		}
		s.Unlock()

		if numSessions > 0 {
			log.Info("waiting on %d connected sessions to shut down", numSessions)
		}
		s.waitOnSessions.Wait()
		log.Info("all sessions have shut down")
	}()
	wg.Wait()
	s.Info("shutdown procedure complete")
}
