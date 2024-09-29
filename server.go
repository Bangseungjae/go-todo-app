package main

import (
	"context"
	"errors"
	"golang.org/x/sync/errgroup"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	srv *http.Server
	l   net.Listener
}

func NewServer(l net.Listener, mux http.Handler) *Server {
	return &Server{
		srv: &http.Server{Handler: mux},
		l:   l,
	}
}

func (s *Server) Run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		// http.ErrServerClosed는
		// http.Server.Shutdown()이 정상 종료한 것을 보여주는 것으로 아무 문제없다.
		if err := s.srv.Serve(s.l); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("failed to close: %+v", err)
		}
		return nil
	})

	<-ctx.Done()
	if err := s.srv.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}
	// 정상 종료를 기다립니다.
	return eg.Wait()
}
