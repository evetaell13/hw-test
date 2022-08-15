package internalhttp

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Server struct { // TODO
	logger     Logger
	app        Application
	httpServer *http.Server
}

type Logger interface { // TODO
	Info(msg string)
	Error(msg string)
	Warn(msg string)
	Debug(msg string)
}

type Application interface { // TODO
}

func NewServer(logger Logger, app Application) *Server {
	return &Server{
		logger:     logger,
		app:        app,
		httpServer: &http.Server{Addr: ":8081"},
	}
}

func (s *Server) Start(ctx context.Context) error {
	// TODO
	<-ctx.Done()
	// TODO config for port
	http.HandleFunc("/", s.getRoot)
	http.HandleFunc("/hello", s.getHello)
	go func() {
		if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()
	// TODO use context
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return nil
}

func (s *Server) getRoot(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w,
		fmt.Sprintf("IP: %v\n time: %v\n address: %v %v \n",
			r.RemoteAddr,
			time.Now().UTC(),
			r.Host, r.URL.Path,
		),
	)
}

func (s *Server) getHello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, `<h1>Hello!</h1>`)
}
