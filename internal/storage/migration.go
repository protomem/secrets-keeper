package storage

import (
	"context"
	"fmt"

	"github.com/protomem/secrets-keeper/assets"
)

func (s *Storage) Migrate(ctx context.Context) error {
	const op = "storage.Migrate"
	var err error

	migrationsFile, err := assets.Assets.Open("migrations/migrations.sql")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer func() { _ = migrationsFile.Close() }()

	migrationsFileStat, err := migrationsFile.Stat()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	migrations := make([]byte, migrationsFileStat.Size())

	_, err = migrationsFile.Read(migrations)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.db.ExecContext(ctx, string(migrations))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
