package passhash

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var ErrWrongPassword = errors.New("wrong password")

func Generate(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("passhash.Generate: %w", err)
	}

	return string(hash), nil
}

func Compare(password string, hash string) error {
	const op = "passhash.Compare"

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return fmt.Errorf("%s: %w", op, ErrWrongPassword)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
