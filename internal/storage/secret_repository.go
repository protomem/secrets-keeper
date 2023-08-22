package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/protomem/secrets-keeper/internal/model"
)

type SecretTable struct {
	ID        int
	CreatedAt string
	Message   string
}

func (s *Storage) GetSecret(ctx context.Context, id int) (model.Secret, error) {
	const op = "storage.GetSecret"
	var err error

	query := `
        SELECT * FROM secrets WHERE id = $1 LIMIT 1
    `

	var secretTable SecretTable
	err = s.db.
		QueryRowContext(ctx, query, id).
		Scan(&secretTable.ID, &secretTable.CreatedAt, &secretTable.Message)
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

func (s *Storage) SaveSecret(ctx context.Context, secret model.Secret) (int, error) {
	const op = "storage.SaveSecret"
	var err error

	query := `
        INSERT INTO secrets (created_at, message) VALUES ($1, $2) RETURNING id
    `

	err = s.db.
		QueryRowContext(ctx, query, secret.CreatedAt.Format(time.RFC3339), secret.Message).
		Scan(&secret.ID)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return secret.ID, nil
}

func (s *Storage) RemoveSecret(ctx context.Context, id int) error {
	const op = "storage.RemoveSecret"
	var err error

	query := `
        DELETE FROM secrets WHERE id = $1
    `

	_, err = s.db.
		ExecContext(ctx, query, id)
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
		ID:        secret.ID,
		CreatedAt: createdAt,
		Message:   secret.Message,
	}, nil
}