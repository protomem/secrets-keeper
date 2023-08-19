package storage

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/protomem/secrets-keeper/pkg/logging"
)

type Storage struct {
	logger logging.Logger
	db     *sql.DB
}

func New(ctx context.Context, logger logging.Logger, database string) (*Storage, error) {
	const op = "storage.New"
	var err error

	db, err := sql.Open("sqlite3", database)
	if err != nil {
		return nil, fmt.Errorf("%w: open: %s", err, op)
	}

	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: ping: %s", err, op)
	}

	return &Storage{
		logger: logger.With("module", "storage"),
		db:     db,
	}, nil
}

func (s *Storage) Close(_ context.Context) error {
	err := s.db.Close()
	if err != nil {
		return fmt.Errorf("storage.Close: %w", err)
	}

	return nil
}
