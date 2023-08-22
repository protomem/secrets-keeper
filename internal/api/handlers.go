package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/protomem/secrets-keeper/internal/model"
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
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "server.GetSecret"
		var err error
		ctx := r.Context()
		logger := s.logger.With("operation", op)

		w.Header().Set("Content-Type", "application/json")

		secretKey, ok := mux.Vars(r)["key"]
		if !ok {
			logger.Error("failed to get secret key")

			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "missing secret key",
			})

			return
		}

		secretID, err := strconv.Atoi(secretKey)
		if err != nil {
			logger.Error("failed to convert secret key to id")

			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "invalid secret key",
			})

			return
		}

		secret, err := s.store.GetSecret(ctx, secretID)
		if err != nil {
			logger.Error("failed to get secret", "error", err)

			code := http.StatusInternalServerError
			res := map[string]string{
				"error": "failed to get secret",
			}

			if errors.Is(err, model.ErrSecretNotFound) {
				code = http.StatusNotFound
				res = map[string]string{
					"error": "secret not found",
				}
			}

			w.WriteHeader(code)
			_ = json.NewEncoder(w).Encode(res)

			return
		}

		err = s.store.RemoveSecret(ctx, secretID)
		if err != nil {
			logger.Error("failed to remove secret", "error", err)

			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "failed to get secret",
			})

			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]model.Secret{
			"secret": secret,
		})
	})
}

func (s *Server) handleCreateSecret() http.Handler {
	type request struct {
		Message string `json:"message"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "server.CreateSecret"
		var err error
		ctx := r.Context()
		logger := s.logger.With("operation", op)

		w.Header().Set("Content-Type", "application/json")

		var req request
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			logger.Error("failed to decode request", "error", err)

			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "invalid request",
			})

			return
		}

		// TODO: Add validation
		// TODO: Add expiration for Secret
		// TODO: Crypt Secret

		secretID, err := s.store.SaveSecret(ctx, model.Secret{
			CreatedAt: time.Now(),
			Message:   req.Message,
		})
		if err != nil {
			logger.Error("failed to save secret", "error", err)

			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "failed to save secret",
			})

			return
		}

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"secretKey": strconv.Itoa(secretID),
		})
	})
}
