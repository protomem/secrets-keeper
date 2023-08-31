package usecase

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/protomem/secrets-keeper/internal/cryptor"
	"github.com/protomem/secrets-keeper/internal/model"
	"github.com/protomem/secrets-keeper/internal/passhash"
	"github.com/protomem/secrets-keeper/internal/storage"
	"github.com/protomem/secrets-keeper/pkg/randstr"
)

type UseCaseFunc[I any, O any] func(context.Context, I) (O, error)

type GetSecretDTO struct {
	SecretKey    string
	SecretPhrase string
}

func GetSecret(
	secretRepo *storage.SecretRepository,
	encoder cryptor.Encoder,
	encryptor cryptor.Encryptor,
) UseCaseFunc[GetSecretDTO, model.Secret] {
	return func(ctx context.Context, dto GetSecretDTO) (model.Secret, error) {
		const op = "usecase.GetSecret"
		var err error
		now := time.Now()

		decodedSecretKey, err := encoder.Decode([]byte(dto.SecretKey))
		if err != nil {
			return model.Secret{}, fmt.Errorf("%s: %w", op, err)
		}

		secretKeyParts := bytes.Split(decodedSecretKey, []byte("$"))
		if len(secretKeyParts) != 2 {
			return model.Secret{}, fmt.Errorf("%s: %w", op, errors.New("invalid secret key"))
		}

		accessKey := secretKeyParts[0]
		signingKey := secretKeyParts[1]

		secret, err := secretRepo.GetSecret(ctx, string(accessKey))
		if err != nil {
			return model.Secret{}, fmt.Errorf("%s: %w", op, err)
		}

		if secret.ExpiredAt.Unix() < now.Unix() && secret.ExpiredAt.Unix() > secret.CreatedAt.Unix() {
			return model.Secret{}, fmt.Errorf("%s: %w", op, model.ErrSecretNotFound)
		}

		if secret.SecretPhrase != "" {
			if dto.SecretPhrase == "" {
				return model.Secret{}, fmt.Errorf("%s: %w", op, model.ErrSecretNotFound)
			}

			err = passhash.Compare(dto.SecretPhrase, secret.SecretPhrase)
			if err != nil {
				if errors.Is(err, passhash.ErrWrongPassword) {
					return model.Secret{}, fmt.Errorf("%s: %w", op, model.ErrSecretNotFound)
				}

				return model.Secret{}, fmt.Errorf("%s: %w", op, err)
			}
		}

		signingKey = append(signingKey, []byte(secret.SigningKey)...)
		decryptedMessage, err := encryptor.Decrypt([]byte(secret.Message), []byte(signingKey))
		if err != nil {
			return model.Secret{}, fmt.Errorf("%s: %w", op, err)
		}

		secret.Message = string(decryptedMessage)

		err = secretRepo.RemoveSecret(ctx, string(accessKey))
		if err != nil {
			return model.Secret{}, fmt.Errorf("%s: %w", op, err)
		}

		return secret, nil
	}
}

type CreateSecretDTO struct {
	Message      string
	TTL          int64 // in hours
	SecretPhrase string
}

func CreateSecret(
	secretRepo *storage.SecretRepository,
	encoder cryptor.Encoder,
	encryptor cryptor.Encryptor,
) UseCaseFunc[CreateSecretDTO, string] {
	return func(ctx context.Context, dto CreateSecretDTO) (string, error) {
		const op = "usecase.CreateSecret"
		var err error
		now := time.Now()

		accessKey := []byte(randstr.Gen(8))
		signingKey := []byte(randstr.Gen(16))

		secretKey, err := encoder.Encode(bytes.Join(
			[][]byte{[]byte(accessKey), signingKey[:6]},
			[]byte("$"),
		))
		if err != nil {
			return "", fmt.Errorf("%s: %w", op, err)
		}

		encryptedMessage, err := encryptor.Encrypt([]byte(dto.Message), signingKey)
		if err != nil {
			return "", fmt.Errorf("%s: %w", op, err)
		}

		if dto.SecretPhrase != "" {
			dto.SecretPhrase, err = passhash.Generate(dto.SecretPhrase)
			if err != nil {
				return "", fmt.Errorf("%s: %w", op, err)
			}
		}

		_, err = secretRepo.SaveSecret(ctx, model.Secret{
			CreatedAt:    now,
			ExpiredAt:    now.Add(time.Duration(dto.TTL) * time.Hour),
			AccessKey:    string(accessKey),
			SigningKey:   string(signingKey[6:]),
			SecretPhrase: dto.SecretPhrase,
			Message:      string(encryptedMessage),
		})
		if err != nil {
			return "", fmt.Errorf("%s: %w", op, err)
		}

		return string(secretKey), nil
	}
}
