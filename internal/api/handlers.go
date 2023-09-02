package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/protomem/secrets-keeper/internal/model"
	"github.com/protomem/secrets-keeper/internal/usecase"
	"github.com/protomem/secrets-keeper/pkg/requestid"
)

func (*Server) handleHealthCheck() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	})
}

func (s *Server) handleGetSecret() http.Handler {
	type Request struct {
		SecretPhrase string `json:"secretPhrase"`
	}

	type Response struct {
		Secret model.Secret `json:"secret"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "server.GetSecret"
		var err error

		ctx := r.Context()
		logger := s.logger.With(
			"operation", op,
			requestid.LogKey, requestid.Extract(ctx),
		)

		defer func() {
			if err != nil {
				logger.Error("failed to handle request", "error", err)
			}
		}()

		w.Header().Set("Content-Type", "application/json")

		secretKey, ok := mux.Vars(r)["key"]
		if !ok {
			logger.Error("failed to get secret key")

			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(map[string]string{
				"error": "missing secret key",
			})

			return
		}

		var req Request
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			logger.Error("failed to decode request", "error", err)

			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(map[string]string{
				"error": "invalid request",
			})

			return
		}

		secret, err := usecase.GetSecret(
			s.store.SecretRepo(),
			s.hasher,
			s.encoder,
			s.encryptor,
		)(ctx, usecase.GetSecretDTO{
			SecretKey:    secretKey,
			SecretPhrase: req.SecretPhrase,
		})
		if err != nil {
			logger.Error("failed to get secret", "error", err)

			code := http.StatusInternalServerError
			res := map[string]string{
				"error": "failed to get secret",
			}

			if errors.Is(err, model.ErrSecretNotFound) {
				code = http.StatusNotFound
				res = map[string]string{
					"error": model.ErrSecretNotFound.Error(),
				}
			}

			w.WriteHeader(code)
			err = json.NewEncoder(w).Encode(res)

			return
		}

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(Response{
			Secret: secret,
		})
	})
}

func (s *Server) handleCreateSecret() http.Handler {
	type Request struct {
		Message      string `json:"message"`
		TTL          int64  `json:"ttl"`
		SecretPhrase string `json:"secretPhrase"`
	}

	type Response struct {
		SecretKey        string `json:"secretKey"`
		WithSecretPhrase bool   `json:"withSecretPhrase"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "server.CreateSecret"
		var err error

		ctx := r.Context()
		logger := s.logger.With(
			"operation", op,
			requestid.LogKey, requestid.Extract(ctx),
		)

		defer func() {
			if err != nil {
				logger.Error("failed to handle request", "error", err)
			}
		}()

		w.Header().Set("Content-Type", "application/json")

		var req Request
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			logger.Error("failed to decode request", "error", err)

			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(map[string]string{
				"error": "invalid request",
			})

			return
		}

		secretKey, err := usecase.CreateSecret(
			s.store.SecretRepo(),
			s.hasher,
			s.encoder,
			s.encryptor,
		)(ctx, usecase.CreateSecretDTO{
			Message:      req.Message,
			TTL:          req.TTL,
			SecretPhrase: req.SecretPhrase,
		})
		if err != nil {
			logger.Error("failed to create secret", "error", err)

			code := http.StatusInternalServerError
			res := map[string]string{
				"error": "failed to create secret",
			}

			w.WriteHeader(code)
			err = json.NewEncoder(w).Encode(res)

			return
		}

		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(Response{
			SecretKey:        secretKey,
			WithSecretPhrase: req.SecretPhrase != "",
		})
	})
}
