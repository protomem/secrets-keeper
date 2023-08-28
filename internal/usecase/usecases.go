package usecase

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/protomem/secrets-keeper/internal/cryptor"
	"github.com/protomem/secrets-keeper/internal/model"
	"github.com/protomem/secrets-keeper/internal/storage"
)

type UseCaseFunc[P any, R any] func(context.Context, P) (R, error)

type GetSecretDTO struct {
	SecretKey string
}

func GetSecret(secretRepo *storage.SecretRepository) UseCaseFunc[GetSecretDTO, model.Secret] {
	return func(ctx context.Context, dto GetSecretDTO) (model.Secret, error) {
		const op = "usecase.GetSecret"
		var err error

		decodedSecretKey, err := hex.DecodeString(dto.SecretKey)
		if err != nil {
			return model.Secret{}, fmt.Errorf("%s: %w", op, err)
		}

		secretKeyParts := bytes.Split(decodedSecretKey, []byte("$"))
		if len(secretKeyParts) != 2 {
			return model.Secret{}, fmt.Errorf("%s: %w", op, errors.New("invalid secret key"))
		}

		accessKey := string(secretKeyParts[0])
		signingKey := string(secretKeyParts[1])

		secret, err := secretRepo.GetSecret(ctx, accessKey)
		if err != nil {
			return model.Secret{}, fmt.Errorf("%s: %w", op, err)
		}

		decodedMessage, err := cryptor.Decode(secret.Message)
		if err != nil {
			return model.Secret{}, fmt.Errorf("%s: %w", op, err)
		}

		signingKey = signingKey + secret.SigningKey
		decryptedMessage, err := cryptor.Decrypt(decodedMessage, []byte(signingKey))
		if err != nil {
			return model.Secret{}, fmt.Errorf("%s: %w", op, err)
		}

		secret.Message = string(decryptedMessage)

		err = secretRepo.RemoveSecret(ctx, accessKey)
		if err != nil {
			return model.Secret{}, fmt.Errorf("%s: %w", op, err)
		}

		return secret, nil
	}
}
