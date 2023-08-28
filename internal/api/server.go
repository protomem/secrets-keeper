package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/protomem/secrets-keeper/internal/config"
	"github.com/protomem/secrets-keeper/internal/storage"
	"github.com/protomem/secrets-keeper/pkg/closer"
	"github.com/protomem/secrets-keeper/pkg/logging"
	"github.com/protomem/secrets-keeper/pkg/logging/stdlog"
)

type Server struct {
	conf   config.Config
	logger logging.Logger

	store *storage.Storage

	router *mux.Router
	server *http.Server

	closer *closer.Closer
}

func New(conf config.Config) (*Server, error) {
	const op = "server.New"
	var err error
	ctx := context.Background()

	logger, err := stdlog.New(conf.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("%w: init logger: %s", err, op)
	}

	logger.Debug("server configured ...", "config", conf)

	store, err := storage.New(ctx, logger, conf.Database)
	if err != nil {
		return nil, fmt.Errorf("%w: init storage: %s", err, op)
	}

	err = store.Migrate(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: migrate: %s", err, op)
	}

	router := mux.NewRouter()
	server := &http.Server{
		Addr:    conf.BindAddr,
		Handler: router,
	}

	return &Server{
		conf:   conf,
		logger: logger.With("module", "server"),
		store:  store,
		router: router,
		server: server,
		closer: closer.New(),
	}, nil
}

func (s *Server) Run() error {
	const op = "server.Run"
	var err error
	ctx := context.Background()

	s.registerOnShutdown()
	s.setupRoutes()

	errs := make(chan error)

	go s.startServer(ctx, errs)
	go s.gracefulShutdown(ctx, errs)

	s.logger.Info("server started ...", "addr", s.conf.BindAddr)
	defer s.logger.Info("server stopped")

	err = <-errs
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Server) registerOnShutdown() {
	s.closer.Add(s.server.Shutdown)
	s.closer.Add(s.store.Close)
	s.closer.Add(s.logger.Sync)
}

func (s *Server) setupRoutes() {
	s.router.Use(s.requestID())
	s.router.Use(s.CORS())

	s.router.Handle("/health", s.handleHealthCheck()).Methods(http.MethodGet)

	s.router.Handle("/api/secrets/{key}", s.handleGetSecret()).Methods(http.MethodGet, http.MethodOptions)
	s.router.Handle("/api/secrets", s.handleCreateSecret()).Methods(http.MethodPost, http.MethodOptions)
}

func (s *Server) startServer(_ context.Context, errs chan error) {
	err := s.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		errs <- fmt.Errorf("start server: %w", err)
	}
}

func (s *Server) gracefulShutdown(ctx context.Context, errs chan error) {
	<-wait()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := s.closer.Close(ctx)
	if err != nil {
		errs <- fmt.Errorf("graceful shutdown: %w", err)
	}

	errs <- nil
}

func wait() <-chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	return ch
}
