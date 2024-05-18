package rest

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"

	"github.com/emo2007/block-accounting/examples/license-api/internal/pkg/config"
	"github.com/emo2007/block-accounting/examples/license-api/internal/pkg/logger"
	"github.com/go-chi/chi/v5"
	mw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Server struct {
	*chi.Mux

	ctx context.Context

	log  *slog.Logger
	addr string
	tls  bool

	closeMu sync.RWMutex
	closed  bool
}

func NewServer(
	log *slog.Logger,
	conf config.RestConfig,
) *Server {
	s := &Server{
		log:  log,
		addr: conf.Address,
		tls:  conf.TLS,
	}

	s.buildRouter()

	return s
}

func (s *Server) Serve(ctx context.Context) error {
	s.ctx = ctx

	s.log.Info(
		"starting rest interface",
		slog.String("addr", s.addr),
		slog.Bool("tls", s.tls),
	)

	if s.tls {
		return http.ListenAndServeTLS(s.addr, "/todo", "/todo", s)
	}

	return http.ListenAndServe(s.addr, s)
}

func (s *Server) Close() {
	s.closeMu.Lock()

	s.closed = true

	s.closeMu.Unlock()
}

func (s *Server) buildRouter() {
	router := chi.NewRouter()

	router.Use(mw.Recoverer)
	router.Use(mw.RequestID)
	router.Use(s.handleMw)

	router.Use(render.SetContentType(render.ContentTypeJSON))

	router.Get("/")

	s.Mux = router
}

func (s *Server) handle(
	h func(w http.ResponseWriter, req *http.Request) ([]byte, error),
	method_name string,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		out, err := h(w, r)
		if err != nil {
			s.log.Error(
				"http error",
				slog.String("method_name", method_name),
				logger.Err(err),
			)

			s.responseError(w, err)

			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if _, err = w.Write(out); err != nil {
			s.log.Error(
				"error write http response",
				slog.String("method_name", method_name),
				logger.Err(err),
			)
		}
	}
}

func (s *Server) responseError(w http.ResponseWriter, e error) {
	s.log.Error("error handle request", logger.Err(e))

	apiErr := mapError(e)

	out, err := json.Marshal(apiErr)
	if err != nil {
		s.log.Error("error marshal api error", logger.Err(err))

		return
	}

	w.WriteHeader(apiErr.Code)
	w.Write(out)
}

func (s *Server) handleMw(next http.Handler) http.Handler {
	// todo add rate limiter && cirquit braker

	fn := func(w http.ResponseWriter, r *http.Request) {
		s.closeMu.RLock()
		defer s.closeMu.RUnlock()

		if s.closed { // keep mutex closed
			return
		}

		w.Header().Add("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
