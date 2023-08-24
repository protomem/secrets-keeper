package api

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/protomem/secrets-keeper/internal/cryptor"
	"github.com/protomem/secrets-keeper/internal/model"
	"github.com/protomem/secrets-keeper/pkg/randstr"
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

		secretKeyRaw, err := hex.DecodeString(secretKey)
		if err != nil {
			logger.Error("failed to decode secret key", "error", err)

			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "invalid secret key",
			})

			return
		}

		secretKeyParts := bytes.Split(secretKeyRaw, []byte("$"))
		if len(secretKeyParts) != 2 {
			logger.Error("failed to split secret key", "error", err)

			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "invalid secret key",
			})

			return
		}

		accessKey := string(secretKeyParts[0])
		signingKey := string(secretKeyParts[1])

		logger.Debug("getting kyes", "accessKey", accessKey, "signingKey", signingKey, "secretKey", secretKey)

		secret, err := s.store.GetSecret(ctx, accessKey)
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

		signingKey = signingKey + secret.SigningKey
		decodedMessage, err := cryptor.Decode(secret.Message)
		if err != nil {
			logger.Error("failed to decode message", "error", err)

			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "failed to decode message",
			})

			return
		}

		decryptedMessage, err := cryptor.Decrypt(decodedMessage, []byte(signingKey))
		if err != nil {
			logger.Error("failed to decrypt message", "error", err)

			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "failed to decrypt message",
			})

			return
		}

		secret.Message = string(decryptedMessage)

		err = s.store.RemoveSecret(ctx, accessKey)
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

		accessKey := randstr.Gen(8)
		signingKey := randstr.Gen(8)
		secretKey := hex.EncodeToString(bytes.Join(
			[][]byte{[]byte(accessKey), []byte(signingKey[:4])},
			[]byte("$"),
		))

		encryptedMessage, err := cryptor.Encrypt([]byte(req.Message), []byte(signingKey))
		if err != nil {
			logger.Error("failed to encrypt message", "error", err)

			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "failed to encrypt message",
			})

			return
		}

		encodedMessage, err := cryptor.Encode(encryptedMessage)
		if err != nil {
			logger.Error("failed to encode message", "error", err)

			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "failed to encode message",
			})

			return
		}

		logger.Debug("proccessing message", "message", req.Message, "encryptedMessage", encryptedMessage)

		_, err = s.store.SaveSecret(ctx, model.Secret{
			CreatedAt:  time.Now(),
			AccessKey:  accessKey,
			SigningKey: signingKey[4:],
			Message:    encodedMessage,
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
			"secretKey": secretKey,
		})
	})
}
