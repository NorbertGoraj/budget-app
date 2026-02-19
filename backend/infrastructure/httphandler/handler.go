package httphandler

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"
)

const shutdownTimeout = 10 * time.Second

// Handler wraps an *http.Server and ties its lifecycle to a context.
// Call Serve() in a goroutine; it blocks until ctx is cancelled, then
// drains in-flight requests before returning.
//
// Usage:
//
//	go httphandler.New(srv, ctx, wg).Serve()
type Handler struct {
	server *http.Server
	ctx    context.Context
	wg     *sync.WaitGroup
}

func New(server *http.Server, ctx context.Context, wg *sync.WaitGroup) *Handler {
	return &Handler{server: server, ctx: ctx, wg: wg}
}

func (h *Handler) Serve() {
	go func() {
		log.Printf("HTTP server listening on %s", h.server.Addr)
		if err := h.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server %s error: %v", h.server.Addr, err)
		}
	}()

	h.wg.Add(1)
	<-h.ctx.Done()

	log.Printf("shutting down HTTP server %s...", h.server.Addr)
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := h.server.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server %s forced to shutdown: %v", h.server.Addr, err)
	} else {
		log.Printf("HTTP server %s shutdown", h.server.Addr)
	}

	h.wg.Done()
}
