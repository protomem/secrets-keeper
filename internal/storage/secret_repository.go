package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/protomem/secrets-keeper/internal/model"
	"github.com/protomem/secrets-keeper/pkg/logging"
)

type (
	SecretTable struct {
		ID         int
		CreatedAt  string
		AccessKey  string
		SigningKey string
		Message    string
	}

	SecretRepository struct {
		logger logging.Logger
		db     *sql.DB
	}
)

func (s *Storage) SecretRepo() *SecretRepository {
	return &SecretRepository{
		logger: s.logger.With("repository", "secret"),
		db:     s.db,
	}
}

func (r *SecretRepository) GetSecret(ctx context.Context, accessKey string) (model.Secret, error) {
	const op = "storage.GetSecret"
	var err error

	query := `
        SELECT * FROM secrets WHERE access_key = $1 LIMIT 1
    `

	var secretTable SecretTable
	err = r.db.
		QueryRowContext(ctx, query, accessKey).
		Scan(
			&secretTable.ID,
			&secretTable.CreatedAt,
			&secretTable.AccessKey,
			&secretTable.SigningKey,
			&secretTable.Message,
		)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Secret{}, fmt.Errorf("%s: %w", op, model.ErrSecretNotFound)
		}

		return model.Secret{}, fmt.Errorf("%s: %w", op, err)
	}

	secret, err := mapSeacretTableToSecretModel(secretTable)
	if err != nil {
		return model.Secret{}, fmt.Errorf("%s: %w", op, err)
	}

	return secret, nil
}

func (r *SecretRepository) SaveSecret(ctx context.Context, secret model.Secret) (int, error) {
	const op = "storage.SaveSecret"
	var err error

	query := `
        INSERT INTO secrets (created_at, access_key, signing_key, message) VALUES ($1, $2, $3, $4) RETURNING id
    `

	err = r.db.
		QueryRowContext(
			ctx, query,
			secret.CreatedAt.Format(time.RFC3339),
			secret.AccessKey,
			secret.SigningKey,
			secret.Message,
		).
		Scan(&secret.ID)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return secret.ID, nil
}

func (r *SecretRepository) RemoveSecret(ctx context.Context, accessKey string) error {
	const op = "storage.RemoveSecret"
	var err error

	query := `
        DELETE FROM secrets WHERE access_key = $1
    `

	_, err = r.db.
		ExecContext(ctx, query, accessKey)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func mapSeacretTableToSecretModel(secret SecretTable) (model.Secret, error) {
	createdAt, err := time.Parse(time.RFC3339, secret.CreatedAt)
	if err != nil {
		return model.Secret{}, fmt.Errorf("parse created at: %w", err)
	}

	return model.Secret{
		ID:         secret.ID,
		CreatedAt:  createdAt,
		AccessKey:  secret.AccessKey,
		SigningKey: secret.SigningKey,
		Message:    secret.Message,
	}, nil
}
