package usecase

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/protomem/secrets-keeper/internal/cryptor"
	"github.com/protomem/secrets-keeper/internal/model"
	"github.com/protomem/secrets-keeper/internal/storage"
	"github.com/protomem/secrets-keeper/pkg/randstr"
)

type UseCaseFunc[I any, O any] func(context.Context, I) (O, error)

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

type CreateSecretDTO struct {
	Message string
}

func CreateSecret(secretRepo *storage.SecretRepository) UseCaseFunc[CreateSecretDTO, string] {
	return func(ctx context.Context, dto CreateSecretDTO) (string, error) {
		const op = "usecase.CreateSecret"
		var err error

		accessKey := randstr.Gen(8)
		signingKey := randstr.Gen(8)
		secretKey := hex.EncodeToString(bytes.Join(
			[][]byte{[]byte(accessKey), []byte(signingKey[:4])},
			[]byte("$"),
		))

		encryptedMessage, err := cryptor.Encrypt([]byte(dto.Message), []byte(signingKey))
		if err != nil {
			return "", fmt.Errorf("%s: %w", op, err)
		}

		encodedMessage, err := cryptor.Encode(encryptedMessage)
		if err != nil {
			return "", fmt.Errorf("%s: %w", op, err)
		}

		_, err = secretRepo.SaveSecret(ctx, model.Secret{
			CreatedAt:  time.Now(),
			AccessKey:  accessKey,
			SigningKey: signingKey[4:],
			Message:    encodedMessage,
		})
		if err != nil {
			return "", fmt.Errorf("%s: %w", op, err)
		}

		return secretKey, nil
	}
}
