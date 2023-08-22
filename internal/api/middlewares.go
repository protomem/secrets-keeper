package api

import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func (s *Server) CORS() mux.MiddlewareFunc {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	}).Handler
}
