package http

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type StdServerSet func(*http.Server)

//Deprecated: Use NewServer instead.
func NewServerUseGin(addr string, routers *gin.Engine, middlewares ...gin.HandlerFunc) *ServerUseGin {
	routers.Use(middlewares...)
	return &ServerUseGin{
		Routers: routers,
		addr:    addr,
	}
}

func NewServer(addr string, routers *gin.Engine) *ServerUseGin {
	return &ServerUseGin{
		Routers: routers,
		addr:    addr,
	}
}

type ServerUseGin struct {
	Routers    *gin.Engine
	httpServer *http.Server
	addr       string
}

func (s *ServerUseGin) Start(sets ...StdServerSet) error {
	s.httpServer = &http.Server{
		Addr:    s.addr,
		Handler: s.Routers,
	}
	for _, set := range sets {
		set(s.httpServer)
	}
	err := s.httpServer.ListenAndServe()
	if err != nil {
		return fmt.Errorf("http server stop: %w", err)
	}

	return nil
}

func (s *ServerUseGin) Stop(duration time.Duration) {
	if s.httpServer == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	_ = s.httpServer.Shutdown(ctx)
}
