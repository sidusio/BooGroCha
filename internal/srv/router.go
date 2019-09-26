package srv

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (s *server) newRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(30 * time.Second))

	r.Route("/api", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/", s.saveCredentials)
			r.Delete("/", s.clearCredentials)
		})
		r.Route("/", func(r chi.Router) {
			r.Use(s.middlewareRefreshCredentials)
			r.Use(s.middlewareExtractCredentials)
			r.Get("/ping", func(writer http.ResponseWriter, request *http.Request) {
				fmt.Fprintf(writer, "pong")
			})
		})
	})

	return r
}
