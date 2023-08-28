package storage

import (
	"context"
	"fmt"

	"github.com/protomem/secrets-keeper/assets"
)

func (s *Storage) Migrate(ctx context.Context) error {
	const op = "storage.Migrate"
	var err error

	migrationsFile, err := assets.Assets.ReadFile("migrations/migrations.sql")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.db.ExecContext(ctx, string(migrationsFile))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
