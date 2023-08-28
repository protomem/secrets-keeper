package api

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/protomem/secrets-keeper/pkg/requestid"
	"github.com/rs/cors"
)

func (s *Server) requestID() mux.MiddlewareFunc {
	return requestid.Middleware()
}

func (s *Server) requestLogger() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := s.logger.With(
				"middleware", "requestLogger",
				requestid.LogKey, requestid.Extract(r.Context()),
			)
			handlers.CombinedLoggingHandler(logger, next).ServeHTTP(w, r)
		})
	}
}

func (s *Server) recovery() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := s.logger.With(
				"middleware", "recovery",
				requestid.LogKey, requestid.Extract(r.Context()),
			)
			handlers.RecoveryHandler(
				handlers.RecoveryLogger(logger),
			)(next).ServeHTTP(w, r)
		})
	}
}

func (s *Server) CORS() mux.MiddlewareFunc {
	return cors.New(cors.Options{
		AllowCredentials: true,
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodOptions},
	}).Handler
}
