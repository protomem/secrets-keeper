package api

import (
	"github.com/gorilla/mux"
	"github.com/protomem/secrets-keeper/pkg/requestid"
	"github.com/rs/cors"
)

func (s *Server) requestID() mux.MiddlewareFunc {
	return requestid.Middleware()
}

func (s *Server) CORS() mux.MiddlewareFunc {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	}).Handler
}
