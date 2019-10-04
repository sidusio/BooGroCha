package srv

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func (s *server) newRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(30 * time.Second))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		AllowCredentials:   true,
		MaxAge:             300, // maximum allowed by all mayor browsers
	})
	r.Use(c.Handler)

	r.Route("/api", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/", s.saveCredentials)
			r.Delete("/", s.clearCredentials)
			r.Get("/test", s.testCredentials)
		})
		r.Route("/", func(r chi.Router) {
			r.Use(s.middlewareRefreshCredentials)
			r.Use(s.middlewareExtractCredentials)

			r.Route("/booking", func(r chi.Router) {
				r.Get("/available", s.available)
				r.Get("/", s.bookings)
				r.Delete("/{bookingID}", s.delete)
			})
		})
		r.Get("/ping", func(writer http.ResponseWriter, request *http.Request) {
			fmt.Fprintf(writer, "pong")
		})
	})

	return r
}
